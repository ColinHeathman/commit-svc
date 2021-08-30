package frequency

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/colinheathman/commit-svc/pkg/github"
)

func TestEndpoint(t *testing.T) {
	enp := NewEndpoint(github.NewTestCommitReader())

	req := httptest.NewRequest("GET", "/users?start=2019-09-07&end=2020-09-02", nil)
	resp := httptest.NewRecorder()

	enp.jsonResponder.ServeHTTP(resp, req)

	// Should be able to unmarshal the JSON back into a struct
	result := []CommitCount{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}

}
