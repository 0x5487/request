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
	_httpClient *http.Client
	_timeout    time.Duration
	// ErrTimeout means http request have been timeout
	ErrTimeout = errors.New("request: request timeout")
)

func init() {
	_httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	_timeout = 30 * time.Second
}

// Agent the main struct to handle all http requests
type Agent struct {
	client *http.Client
	err    error

	URL      string
	Method   string
	Header   map[string]string
	QueryStr string
	Body     []byte
	Timeout  time.Duration
}

func newAgent() *Agent {
	agent := &Agent{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
		Header: map[string]string{},
	}
	agent.Header["Accept"] = "application/json"
	agent.Timeout = _timeout
	return agent
}

func newAgentWithClient(client *http.Client) *Agent {
	agent := &Agent{
		client: client,
		Header: map[string]string{},
	}
	agent.Header["Accept"] = "application/json"
	agent.Timeout = _timeout
	return agent
}

func (agent *Agent) getTransport() *http.Transport {
	trans, _ := agent.client.Transport.(*http.Transport)
	return trans
}

// GET return Agent that uses HTTP GET method with target URL
func (agent *Agent) GET(targetURL string) *Agent {
	agent.Method = "GET"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.URL = targetURL
	return agent
}

// POST return Agent that uses HTTP POST method with target URL
func (agent *Agent) POST(targetURL string) *Agent {
	agent.Method = "POST"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.URL = targetURL
	return agent
}

// PUT return Agent that uses HTTP PUT method with target URL
func (agent *Agent) PUT(targetURL string) *Agent {
	agent.Method = "PUT"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.URL = targetURL
	return agent
}

// DELETE return Agent that uses HTTP PUT method with target URL
func (agent *Agent) DELETE(targetURL string) *Agent {
	agent.Method = "DELETE"
	_, err := url.Parse(targetURL)
	if err != nil {
		agent.err = err
	}
	agent.URL = targetURL
	return agent
}

// Set that set HTTP header to agent
func (agent *Agent) Set(key, val string) *Agent {
	agent.Header[key] = val
	return agent
}

// SetTimeout set timeout for agent.  The default value is 30 seconds.
func (agent *Agent) SetTimeout(timeout time.Duration) *Agent {
	if timeout > 0 {
		agent.Timeout = timeout
	}
	return agent
}

// SetProxyURL set the simple proxy with fixed proxy url
func (agent *Agent) SetProxyURL(proxyURL string) *Agent {
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
func (agent *Agent) SendBytes(bytes []byte) *Agent {
	agent.Body = bytes
	return agent
}

// SendJSON send json to target URL
func (agent *Agent) SendJSON(v interface{}) *Agent {
	agent.Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		agent.err = err
	}
	agent.Body = b
	return agent
}

// SendXML send json to target URL
func (agent *Agent) SendXML(v interface{}) *Agent {
	agent.Set("Content-Type", "application/xml")
	b, err := xml.Marshal(v)
	if err != nil {
		agent.err = err
	}
	agent.Body = b
	return agent
}

// Send send string to target URL
func (agent *Agent) Send(body string) *Agent {
	agent.Set("Content-Type", "application/x-www-form-urlencoded")
	agent.Body = []byte(body)
	return agent
}

// Query set  querystring to target URL
func (agent *Agent) Query(querystring string) *Agent {
	agent.QueryStr = querystring
	return agent
}

// End start execute agent
func (agent *Agent) End() (*Response, error) {
	if agent.err != nil {
		return nil, agent.err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if agent.Timeout > 0 {
		ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, agent.Timeout)
		ctx = ctxWithTimeout
		cancel = cancelWithTimeout
	}

	// create new request
	url := agent.URL + agent.QueryStr
	outReq, err := http.NewRequest(agent.Method, url, bytes.NewReader(agent.Body))
	if err != nil {
		return nil, err
	}

	// copy Header
	for k, val := range agent.Header {
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
