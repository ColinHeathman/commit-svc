package api

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponder(t *testing.T) {

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	type TestMsg struct {
		Msg string `json:"msg"`
	}

	testStruct := &TestMsg{
		"test",
	}

	getter := &JSONResponder{
		ModelFromRequest: func(h *http.Request) (interface{}, error) {
			return &testStruct, nil
		},
	}

	getter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Error("Wrong http status code")
	}
	if resp.Body.String() != "{\"msg\":\"test\"}" {
		t.Error("Wrong http body")
	}
}

func TestJSONResponderError(t *testing.T) {

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	getter := &JSONResponder{
		ModelFromRequest: func(h *http.Request) (interface{}, error) {
			return nil, fmt.Errorf("Test error")
		},
	}

	getter.ServeHTTP(resp, req)

	if resp.Code != 500 {
		t.Error("Wrong http status code")
	}

	if resp.Body.String() != "{\"error\": \"Test error\"}" {
		t.Error("Wrong http body")
	}
}

func TestJSONResponderNull(t *testing.T) {

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	getter := &JSONResponder{
		ModelFromRequest: func(h *http.Request) (interface{}, error) {
			return nil, nil
		},
	}

	getter.ServeHTTP(resp, req)

	if resp.Code != 404 {
		t.Error("Wrong http status code")
	}

	if resp.Body.String() != "{}" {
		t.Error("Wrong http body")
	}
}

func TestJSONResponderMarshalFailure(t *testing.T) {

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	type TestMsg struct {
		Msg float64 `json:"msg"`
	}

	testStruct := &TestMsg{
		math.Inf(1),
	}

	getter := &JSONResponder{
		ModelFromRequest: func(h *http.Request) (interface{}, error) {
			return &testStruct, nil
		},
	}

	getter.ServeHTTP(resp, req)

	if resp.Code != 500 {
		t.Error("Wrong http status code")
	}

	if resp.Body.String() != "{\"error\": \"failed to marshal json: json: unsupported value: +Inf\"}" {
		t.Error("Wrong http body")
	}
}
