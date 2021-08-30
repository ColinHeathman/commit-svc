package api

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestInternalServerError(t *testing.T) {

	responseRecorder := httptest.NewRecorder()
	InternalServerError(responseRecorder, fmt.Errorf("Test error"))
	if responseRecorder.Code != 500 {
		t.Log("Wrong http status code")
		t.Fail()
	}
	if responseRecorder.Body.String() != "{\"error\": \"Test error\"}" {
		t.Log("Wrong http body")
		t.Fail()
	}
}

func TestNotFound(t *testing.T) {

	responseRecorder := httptest.NewRecorder()
	NotFound(responseRecorder)
	if responseRecorder.Code != 404 {
		t.Log("Wrong http status code")
		t.Fail()
	}
	if responseRecorder.Body.String() != "{}" {
		t.Log("Wrong http body")
		t.Fail()
	}
}
