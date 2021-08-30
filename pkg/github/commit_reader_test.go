package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// TestReadAllCommits tests that the TestCommitReader can iterate through the paginated data
func TestReadAllCommits(t *testing.T) {
	reader := NewTestCommitReader()
	for page := range reader.StreamAllCommits() {
		if page.Err != nil {
			t.Error(page.Err)
		}
		if len(page.Commits) == 0 {
			t.Errorf("got empty page")
		}
	}
}

// ErrorHTTPClient client to return errors
type ErrorHTTPClient struct{}

func (*ErrorHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("test error")
}

// TestBadRepository tests that NewCommitReader will return an error if a bad repository path is used
func TestBadRepository(t *testing.T) {
	_, err := NewCommitReader(&ErrorHTTPClient{}, "foobar")

	// Expect an error here
	if err == nil {
		t.Errorf("didn't get error")
	}
}

// TestErrorReturned tests that the TestCommitReader can iterate through the paginated data
func TestErrorReturned(t *testing.T) {
	rdr, err := NewCommitReader(&ErrorHTTPClient{}, "owner/repo")
	if err != nil {
		t.Fatal(err)
	}

	for page := range rdr.StreamAllCommits() {

		// Expect an error here
		if page.Err == nil {
			t.Errorf("didn't get error")
		}
	}
}

type ErrorReader struct{}

func (*ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

// ErrorHTTPClient client to return errors
type ErrorReadingHTTPClient struct{}

func (*ErrorReadingHTTPClient) Do(req *http.Request) (*http.Response, error) {

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(&ErrorReader{}),
	}, nil
}

// TestErrorReturned tests that the TestCommitReader can iterate through the paginated data
func TestErrorReadingResponse(t *testing.T) {
	rdr, err := NewCommitReader(&ErrorReadingHTTPClient{}, "owner/repo")
	if err != nil {
		t.Fatal(err)
	}

	for page := range rdr.StreamAllCommits() {

		// Expect an error here
		if page.Err == nil {
			t.Errorf("didn't get error")
		}
	}
}

// BadDataClient client to return bad data
type BadDataClient struct{}

func (*BadDataClient) Do(req *http.Request) (*http.Response, error) {

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("something not json ")),
	}, nil
}

// TestErrorUnmarshallingResponse tests that the TestCommitReader can iterate through the paginated data
func TestErrorUnmarshallingResponse(t *testing.T) {
	rdr, err := NewCommitReader(&BadDataClient{}, "owner/repo")
	if err != nil {
		t.Fatal(err)
	}

	for page := range rdr.StreamAllCommits() {

		// Expect an error here
		if page.Err == nil {
			t.Errorf("didn't get error")
		}
	}
}

// WrongStatusCodeClient client to return the wrong status code
type WrongStatusCodeClient struct{}

func (*WrongStatusCodeClient) Do(req *http.Request) (*http.Response, error) {

	return &http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(strings.NewReader("{\"error\":\"something\"}")),
	}, nil
}

// TestErrorReturned tests that the TestCommitReader can iterate through the paginated data
func TestErrorWrongStatusCode(t *testing.T) {
	rdr, err := NewCommitReader(&WrongStatusCodeClient{}, "owner/repo")
	if err != nil {
		t.Fatal(err)
	}

	for page := range rdr.StreamAllCommits() {

		// Expect an error here
		if page.Err == nil {
			t.Errorf("didn't get error")
		}
	}
}
