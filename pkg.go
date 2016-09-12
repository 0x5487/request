package request

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
