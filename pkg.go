package request

import "net/http"

// New create a new RequestAgent instance
func New() Agent {
	return newAgentWithClient(_httpClient)
}

// NewWithClient create a new RequestAgent instance with custom http client
func NewWithClient(client *http.Client) Agent {
	return newAgentWithClient(client)
}

// GET return RequestAgent that uses HTTP GET method with target URL
func GET(url string) Agent {
	return newAgentWithClient(_httpClient).GET(url)
}

// POST return RequestAgent that uses HTTP POST method with target URL
func POST(url string) Agent {
	return newAgentWithClient(_httpClient).POST(url)
}

// PUT return RequestAgent that uses HTTP PUT method with target URL
func PUT(url string) Agent {
	return newAgentWithClient(_httpClient).PUT(url)
}

// DELETE return RequestAgent that uses HTTP DELETE method with target URL
func DELETE(url string) Agent {
	return newAgentWithClient(_httpClient).DELETE(url)
}

// SetMethod return RequestAgent that uses HTTP method with target URL
func SetMethod(method, targetURL string) Agent {
	return newAgentWithClient(_httpClient).SetMethod(method, targetURL)
}
