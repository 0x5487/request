package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	_httpClient *http.Client
	ErrTimeout  = errors.New("request: request timeout")
)

func init() {
	_httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(30) * time.Second,
	}
}

type RequestAgent struct {
	client *http.Client
	err    error

	URL      string
	Method   string
	Header   map[string]string
	QueryStr string
	Body     []byte
}

func newRequestAgent(method, targetURL string) *RequestAgent {
	agent := &RequestAgent{
		client: _httpClient,
		Method: method,
		URL:    targetURL,
		Header: map[string]string{},
	}
	agent.Header["Accept"] = "application/json"

	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	return agent
}

func (source *RequestAgent) Set(key, val string) *RequestAgent {
	source.Header[key] = val
	return source
}

func (source *RequestAgent) SendBytes(bytes []byte) *RequestAgent {
	source.Body = bytes
	return source
}

func (source *RequestAgent) SendJSON(v interface{}) *RequestAgent {
	source.Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		source.err = err
	}
	source.Body = b
	return source
}

func (source *RequestAgent) Send(body string) *RequestAgent {
	source.Body = []byte(body)
	return source
}

func (source *RequestAgent) Query(querystring string) *RequestAgent {
	source.QueryStr = querystring
	return source
}

func (source *RequestAgent) End() (*Response, error) {
	if source.err != nil {
		return nil, source.err
	}

	// create new request
	url := source.URL + source.QueryStr
	outReq, err := http.NewRequest(source.Method, url, bytes.NewReader(source.Body))
	if err != nil {
		return nil, err
	}

	// copy Header
	for k, val := range source.Header {
		outReq.Header.Add(k, val)
	}

	// send to target
	resp, err := source.client.Do(outReq)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer respClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	result.setResp(resp)
	result.Body = body

	if resp.StatusCode >= 200 && resp.StatusCode < 300 || resp.StatusCode == 304 {
		result.OK = true
	}

	return result, nil
}

func respClose(body io.ReadCloser) error {
	if body == nil {
		return nil
	}
	if _, err := io.Copy(ioutil.Discard, body); err != nil {
		return err
	}
	return body.Close()
}
