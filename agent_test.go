package request

import (
	"encoding/xml"
	"encoding/json"
	"net/http/httptest"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBody(t *testing.T) {
	body := "hell world"
	handler := func(w http.ResponseWriter, r *http.Request) {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(bs) != body {
			t.Errorf("body = %s; want = %s", bs, body)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))

	// string
	_, err := POST(ts.URL).Send(body).End()
	if err != nil {
		t.Fatal(err)
	}

	// []byte
	_, err = POST(ts.URL).SendBytes([]byte(body)).End()
	if err != nil {
		t.Fatal(err)
	}
}


func TestBodyJSON(t *testing.T) {
	type content struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	c := content{
		Code: 1,
		Msg:  "ok",
	}
	checkData := func(data []byte) {
		var cc content
		err := json.Unmarshal(data, &cc)
		if err != nil {
			t.Fatal(err)
		}
		if cc != c {
			t.Errorf("request body = %+v; want = %+v", cc, c)
		}
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		checkData(data)
	})

	ts := httptest.NewServer(handler)
	_, err := POST(ts.URL).SendJSON(c).End()
	if err != nil {
		t.Fatal(err)
	}
}


func TestBodyXML(t *testing.T) {
	type content struct {
		Code int    `xml:"code"`
		Msg  string `xml:"msg"`
	}
	c := content{
		Code: 1,
		Msg:  "ok",
	}
	checkData := func(data []byte) {
		var cc content
		err := xml.Unmarshal(data, &cc)
		if err != nil {
			t.Fatal(err)
		}
		if cc != c {
			t.Errorf("request body = %+v; want = %+v", cc, c)
		}
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		checkData(data)
	})

	ts := httptest.NewServer(handler)
	_, err := PUT(ts.URL).SendXML(c).End()
	if err != nil {
		t.Fatal(err)
	}
}

type Header map[string]string
func TestHeader(t *testing.T) {
	header := Header{
		"User-Agent":    "V1.0.0",
		"Authorization": "abc",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		for key, value := range header {
			if v := r.Header.Get(key); value != v {
				t.Errorf("header %q = %s; want = %s", key, v, value)
			}
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	_, err := POST(ts.URL).
		Set("Authorization", "abc").
		Set("User-Agent", "V1.0.0").
		End()

	if err != nil {
		t.Fatal(err)
	}

}