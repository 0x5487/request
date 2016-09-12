package request

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	*http.Response
	OK   bool
	Body []byte
}

func (source *Response) setResp(aa *http.Response) {
	source.Response = aa
}

func (source *Response) String() string {
	return string(source.Body)
}

func (source *Response) JSON(val interface{}) error {
	err := json.Unmarshal(source.Body, val)
	if err != nil {
		return err
	}
	return nil
}
