package request

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Response respresent result of http request
type Response struct {
	*http.Response
	OK   bool
	Body []byte
}

func (resp *Response) setResp(aa *http.Response) {
	resp.Response = aa
}

func (resp *Response) String() string {
	return string(resp.Body)
}

// JSON allow result of http request bind to json struct
func (resp *Response) JSON(val interface{}) error {
	err := json.Unmarshal(resp.Body, val)
	if err != nil {
		return err
	}
	return nil
}

// XML allow result of http request bind to xml struct
func (resp *Response) XML(val interface{}) error {
	err := xml.Unmarshal(resp.Body, val)
	if err != nil {
		return err
	}
	return nil
}
