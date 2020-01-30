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
	"time"
)

var (
	// httpClient should be kept for reuse purpose
	_httpClient *http.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	// ErrTimeout means http request have been timeout
	ErrTimeout = errors.New("request: request timeout")
)

// Agent the main struct to handle all http requests
type agent struct {
	client   *http.Client
	err      error
	url      string
	method   string
	header   map[string]string
	queryStr string
	body     []byte
	timeout  time.Duration
}

// Agenter return represents an interface which handles request actions and do chain job
type Agenter interface {
	GET(targetURL string) Agenter
	POST(targetURL string) Agenter
	PUT(targetURL string) Agenter
	DELETE(targetURL string) Agenter
	Set(key, val string) Agenter
	SetClient(client *http.Client) Agenter
	Timeout() time.Duration
	SetTimeout(timeout time.Duration) Agenter
	SetProxyURL(proxyURL string) Agenter
	SendBytes(bytes []byte) Agenter
	SendJSON(v interface{}) Agenter
	SendXML(v interface{}) Agenter
	Send(body string) Agenter
	Query(querystring string) Agenter
	End() (*Response, error)
}

func newAgentWithClient(client *http.Client) Agenter {
	agent := agent{
		client: client,
		header: map[string]string{},
	}
	agent.header["Accept"] = "application/json"
	agent.timeout = 30 * time.Second
	return agent
}

func (agent agent) getTransport() *http.Transport {
	trans, _ := agent.client.Transport.(*http.Transport)
	return trans
}

// GET return Agent that uses HTTP GET method with target URL
func (agent agent) GET(targetURL string) Agenter {
	agent.method = "GET"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.url = targetURL
	return agent
}

// POST return Agent that uses HTTP POST method with target URL
func (agent agent) POST(targetURL string) Agenter {
	agent.method = "POST"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.url = targetURL
	return agent
}

// PUT return Agent that uses HTTP PUT method with target URL
func (agent agent) PUT(targetURL string) Agenter {
	agent.method = "PUT"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.url = targetURL
	return agent
}

// DELETE return Agent that uses HTTP PUT method with target URL
func (agent agent) DELETE(targetURL string) Agenter {
	agent.method = "DELETE"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.url = targetURL
	return agent
}

// Set that set HTTP header to agent
func (agent agent) Set(key, val string) Agenter {
	agent.header[key] = val
	return agent
}

// SetClient allow to set a custom client to agent
func (agent agent) SetClient(client *http.Client) Agenter {
	agent.client = client
	return agent
}

// Timeout returns timeout value. The default value is 30 seconds.
func (agent agent) Timeout() time.Duration {
	return agent.timeout
}

// SetTimeout set timeout for agent.  The default value is 30 seconds.
func (agent agent) SetTimeout(timeout time.Duration) Agenter {
	if timeout > 0 {
		agent.timeout = timeout
	}
	return agent
}

// SetProxyURL set the simple proxy with fixed proxy url
func (agent agent) SetProxyURL(proxyURL string) Agenter {
	trans := agent.getTransport()
	if trans == nil {
		agent.err = errors.New("request: no transport")
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		agent.err = err
	}
	trans.Proxy = http.ProxyURL(u)
	return agent
}

// SendBytes send bytes to target URL
func (agent agent) SendBytes(bytes []byte) Agenter {
	agent.body = bytes
	return agent
}

// SendJSON send json to target URL
func (agent agent) SendJSON(v interface{}) Agenter {
	agent.Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		agent.err = err
	}
	agent.body = b
	return agent
}

// SendXML send json to target URL
func (agent agent) SendXML(v interface{}) Agenter {
	agent.Set("Content-Type", "application/xml")
	b, err := xml.Marshal(v)
	if err != nil {
		agent.err = err
	}
	agent.body = b
	return agent
}

// Send send string to target URL
func (agent agent) Send(body string) Agenter {
	agent.Set("Content-Type", "application/x-www-form-urlencoded")
	agent.body = []byte(body)
	return agent
}

// Query set  querystring to target URL
func (agent agent) Query(querystring string) Agenter {
	agent.queryStr = querystring
	return agent
}

// End start execute agent
func (agent agent) End() (*Response, error) {
	if agent.err != nil {
		return nil, agent.err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if agent.timeout > 0 {
		ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, agent.timeout)
		ctx = ctxWithTimeout
		cancel = cancelWithTimeout
	}

	// create new request
	url := agent.url + agent.queryStr
	outReq, err := http.NewRequest(agent.method, url, bytes.NewReader(agent.body))
	if err != nil {
		return nil, err
	}

	// copy Header
	for k, val := range agent.header {
		outReq.Header.Add(k, val)
	}

	// send to target
	resp, err := agent.client.Do(outReq.WithContext(ctx))
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
