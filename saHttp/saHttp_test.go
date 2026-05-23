package saHttp

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpload(t *testing.T) {
	//res, err := Upload("https://v3.biyingniao.com/web/bang/upload", map[string]string{"name": "file", "path": "./saHttp_test.go"}, nil, nil)
	//fmt.Println(err)
	//fmt.Println(res)
}

func TestErrCallback(t *testing.T) {
	SetErrCallback(func(request string) {
		fmt.Println(request)
	})

	_ = Do(Params{
		Url:             "vv",
		Query:           nil,
		Header:          nil,
		Body:            nil,
		Timeout:         0,
		CallbackWhenErr: true,
	}, nil)
}

func TestDoDecodesXMLResponse(t *testing.T) {
	type response struct {
		XMLName xml.Name `xml:"response"`
		Message string   `xml:"message"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(`<response><message>ok</message></response>`))
	}))
	defer server.Close()

	var out response
	err := Do(Params{Url: server.URL}, &out)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if out.Message != "ok" {
		t.Fatalf("Message = %q, want ok", out.Message)
	}
}

func TestDoFormRejectsEmptyURLWithoutPanic(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatalf("DoForm panicked for empty URL: %v", e)
		}
	}()

	err := DoForm(FormParams{}, nil)
	if err == nil {
		t.Fatal("DoForm returned nil error for empty URL")
	}
}

func TestDoFormReturnsErrorForMalformedURLWithoutPanic(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatalf("DoForm panicked for malformed URL: %v", e)
		}
	}()

	err := DoForm(FormParams{Url: "\n"}, nil)
	if err == nil {
		t.Fatal("DoForm returned nil error for malformed URL")
	}
}
