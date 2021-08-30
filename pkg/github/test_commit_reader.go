package github

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	commitsPG1 []byte
	commitsPG2 []byte
)

// NewTestCommitReader initializes a CommitReader that returns test data
func NewTestCommitReader() CommitReader {

	// Load test data

	jsonDir, ok := os.LookupEnv("TEST_JSON_DIR")
	if !ok {
		jsonDir = "test/"
	}

	// Load testing data commitsPG1.json
	file, err := os.Open(path.Join(jsonDir, "commits_pg1.json"))
	if err != nil {
		log.Fatal(err)
	}

	commitsPG1, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Load testing data commitsPG2.json
	file, err = os.Open(path.Join(jsonDir, "commits_pg2.json"))
	if err != nil {
		log.Fatal(err)
	}

	commitsPG2, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	rdr, _ := NewCommitReader(&DummyHTTPClient{}, "owner/repo")
	return rdr
}

// Dummy client to return the test data to HTTP requests
type DummyHTTPClient struct{}

func (*DummyHTTPClient) Do(req *http.Request) (*http.Response, error) {

	// default to page 1 if there isn't one found in the GET query
	q := req.URL.Query()
	pageParameters := q["page"]
	page := "1"
	if len(pageParameters) > 0 {
		page = pageParameters[0]
	}

	if page == "1" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(commitsPG1)),
		}, nil

	} else if page == "2" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(commitsPG2)),
		}, nil

	}

	return &http.Response{
		StatusCode: 404,
		Body:       ioutil.NopCloser(strings.NewReader("{}")),
	}, nil
}
