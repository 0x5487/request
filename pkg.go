package request

func Get(url string) *RequestAgent {
	return newRequestAgent(url)
}

func Post(url string) *RequestAgent {
	return newRequestAgent(url)
}

func Put(url string) *RequestAgent {
	return newRequestAgent(url)
}

func Delete(url string) *RequestAgent {
	return newRequestAgent(url)
}

func (source *RequestAgent) Set(key, val string) *RequestAgent {
	source.Header["key"] = val
	return source
}

func (source *RequestAgent) Send(data string) *RequestAgent {
	source.SendStr = data
	return source
}

func (source *RequestAgent) Query(data string) *RequestAgent {
	source.QueryStr = data
	return source
}

func (source *RequestAgent) End() (*Response, error) {

}
