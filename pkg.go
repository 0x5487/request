package request

// New create a new RequestAgent instance
func New() *Agent {
	return newAgent()
}

// GET return RequestAgent that uses HTTP GET method with target URL
func GET(url string) *Agent {
	return newAgentWithClient(_httpClient).GET(url)
}

// POST return RequestAgent that uses HTTP POST method with target URL
func POST(url string) *Agent {
	return newAgentWithClient(_httpClient).POST(url)
}

// PUT return RequestAgent that uses HTTP PUT method with target URL
func PUT(url string) *Agent {
	return newAgentWithClient(_httpClient).PUT(url)
}

// DELETE return RequestAgent that uses HTTP DELETE method with target URL
func DELETE(url string) *Agent {
	return newAgentWithClient(_httpClient).DELETE(url)
}
