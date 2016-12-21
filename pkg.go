package request

import "time"

func New(method, url string) *RequestAgent {
	return newRequestAgent(method, url)
}

func GET(url string) *RequestAgent {
	return newRequestAgent("GET", url)
}

func POST(url string) *RequestAgent {
	return newRequestAgent("POST", url)
}

func PUT(url string) *RequestAgent {
	return newRequestAgent("PUT", url)
}

func DELETE(url string) *RequestAgent {
	return newRequestAgent("DELETE", url)
}

func SetTimeout(duration int64) {
	_httpClient.Timeout = time.Duration(duration) * time.Second
}
