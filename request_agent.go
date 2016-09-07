package request

import "net/http"

type RequestAgent struct {
	client   *http.Client
	URL      string
	Method   string
	Header   map[string]string
	QueryStr string
	SendStr  string
	Errors   []error
}

func newRequestAgent(url string) *RequestAgent {
	return &RequestAgent{
		URL: url,
	}
}