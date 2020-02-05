package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// httpClient should be kept for reuse purpose
	_httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	_timeout = 30 * time.Second

	// ErrTimeout means http request have been timeout
	ErrTimeout = errors.New("request: request timeout")
)

// Agent the main struct to handle all http requests
type Agent struct {
	client *http.Client
	err    error

	URL     string
	Method  string
	Header  map[string]string
	Body    []byte
	Timeout time.Duration
}

func newAgentWithClient(client *http.Client) Agent {
	agent := Agent{
		client: client,
		Header: map[string]string{},
	}
	agent.Header["Accept"] = "application/json"
	agent.Timeout = _timeout
	return agent
}

func (a Agent) getTransport() *http.Transport {
	trans, _ := a.client.Transport.(*http.Transport)
	return trans
}

// SetMethod return Agent that uses HTTP method with target URL
func (a Agent) SetMethod(method, targetURL string) Agent {
	a.Method = strings.ToUpper(method)

	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// SetClient return Agent with target URL
func (a Agent) SetClient(client *http.Client) Agent {
	a.client = client
	return a
}

// GET return Agent that uses HTTP GET method with target URL
func (a Agent) GET(targetURL string) Agent {
	a.Method = "GET"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// POST return Agent that uses HTTP POST method with target URL
func (a Agent) POST(targetURL string) Agent {
	a.Method = "POST"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// PUT return Agent that uses HTTP PUT method with target URL
func (a Agent) PUT(targetURL string) Agent {
	a.Method = "PUT"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// DELETE return Agent that uses HTTP PUT method with target URL
func (a Agent) DELETE(targetURL string) Agent {
	a.Method = "DELETE"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// Set that set HTTP header to agent
func (a Agent) Set(key, val string) Agent {
	newHeader := map[string]string{}

	if a.Header != nil {
		for k, val := range a.Header {
			newHeader[k] = val
		}
	}

	newHeader[key] = val
	a.Header = newHeader

	return a
}

// SetTimeout set timeout for agent.  The default value is 30 seconds.
func (a Agent) SetTimeout(timeout time.Duration) Agent {
	if timeout > 0 {
		a.Timeout = timeout
	}
	return a
}

// SetProxyURL set the simple proxy with fixed proxy url
func (a Agent) SetProxyURL(proxyURL string) Agent {
	trans := a.getTransport()
	if trans == nil {
		a.err = errors.New("request: no transport")
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		a.err = err
	}
	trans.Proxy = http.ProxyURL(u)
	return a
}

// SendBytes send bytes to target URL
func (a Agent) SendBytes(bytes []byte) Agent {
	a.Body = bytes
	return a
}

// SendJSON send json to target URL
func (a Agent) SendJSON(v interface{}) Agent {
	newAgent := a.Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		newAgent.err = err
	}
	return newAgent.SendBytes(b)
}

// SendXML send json to target URL
func (a Agent) SendXML(v interface{}) Agent {
	newAgent := a.Set("Content-Type", "application/xml")
	b, err := xml.Marshal(v)
	if err != nil {
		newAgent.err = err
	}
	return newAgent.SendBytes(b)
}

// Send send string to target URL
func (a Agent) Send(body string) Agent {
	newAgent := a.Set("Content-Type", "application/x-www-form-urlencoded")
	return newAgent.SendBytes([]byte(body))
}

// End start execute agent
func (a Agent) End() (*Response, error) {
	if a.err != nil {
		return nil, a.err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if a.Timeout > 0 {
		ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, a.Timeout)
		ctx = ctxWithTimeout
		cancel = cancelWithTimeout
	}

	// create new request
	url := a.URL
	outReq, err := http.NewRequest(a.Method, url, bytes.NewReader(a.Body))
	if err != nil {
		return nil, err
	}

	// copy Header
	for k, val := range a.Header {
		outReq.Header.Add(k, val)
	}

	// send to target
	resp, err := a.client.Do(outReq.WithContext(ctx))
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
