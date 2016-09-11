package request

import "net/http"

type Response struct {
	*http.Response
	OK   bool
	Body []byte
}

func (source *Response) setResp(aa *http.Response) {
	source.Response = aa
}

func (source *Response) Text() string {
	return string(source.Body)
}
