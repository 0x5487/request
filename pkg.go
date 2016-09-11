package request

func Get(url string) *RequestAgent {
	return newRequestAgent("GET", url)
}

func Post(url string) *RequestAgent {
	return newRequestAgent("POST", url)
}

func Put(url string) *RequestAgent {
	return newRequestAgent("PUT", url)
}

func Delete(url string) *RequestAgent {
	return newRequestAgent("DELETE", url)
}
