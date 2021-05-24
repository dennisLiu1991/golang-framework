package main

import (
	"bytes"
	"igin/igin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sendHttpForUnitTest(uri string, method string, reqBody []byte, e *igin.Engine) (int, []byte, error) {
	req := httptest.NewRequest(method, uri, bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)

	result := w.Result()
	if result.Body != nil {
		defer result.Body.Close()
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return 0, nil, err
	}
	return result.StatusCode, body, nil
}

func Test_main(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		body         []byte
		wantCode     int
		wantResponse string
	}{
		{
			name:         "404 request",
			method:       "GET",
			path:         "/aabababab",
			wantCode:     http.StatusNotFound,
			wantResponse: "404 not found\n",
		},
		{
			name:         "success request to /a",
			method:       "GET",
			path:         "/a",
			body:         nil,
			wantCode:     http.StatusOK,
			wantResponse: "success",
		},
		{
			name:         "success request to /b,json response",
			method:       "GET",
			path:         "/b",
			body:         nil,
			wantCode:     http.StatusOK,
			wantResponse: `{"hello":"b"}`,
		},
	}
	e := initEngine()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := sendHttpForUnitTest(tt.path, tt.method, tt.body, e)
			if err != nil {
				t.Errorf("err : %v", err)
				return
			}
			if code != tt.wantCode {
				t.Errorf("code = %d,wanted code = %d", code, tt.wantCode)
				return
			}
			if string(body) != tt.wantResponse {
				t.Errorf("body = %s,wanted body = %s", string(body), tt.wantResponse)
				return
			}
		})
	}
}
