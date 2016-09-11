package request

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestAgent struct {
	client   *http.Client
	URL      string
	Method   string
	Header   map[string]string
	QueryStr string
	Body     []byte
	Errors   []error
}

func newRequestAgent(method, url string) *RequestAgent {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(30) * time.Second,
	}
	return &RequestAgent{
		client: client,
		Method: method,
		URL:    url,
	}
}

func (source *RequestAgent) Set(key, val string) *RequestAgent {
	source.Header["key"] = val
	return source
}

func (source *RequestAgent) Query(data string) *RequestAgent {
	source.QueryStr = data
	return source
}

func (source *RequestAgent) End() (*Response, []error) {
	if len(source.Errors) > 0 {
		return nil, source.Errors
	}

	// create new request
	outReq, err := http.NewRequest(source.Method, source.URL, bytes.NewReader(source.Body))
	if err != nil {
		source.Errors = append(source.Errors, err)
		return nil, source.Errors
	}

	// copy Header
	for k, val := range source.Header {
		outReq.Header.Add(k, val)
	}

	// send to target
	resp, err := source.client.Do(outReq)
	if err != nil {
		source.Errors = append(source.Errors, err)
		return nil, source.Errors
	}
	defer respClose(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)

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
