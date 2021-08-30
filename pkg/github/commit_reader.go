package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/colinheathman/commit-svc/pkg/dateutil"
	"github.com/colinheathman/commit-svc/pkg/http_util"
)

// CommitReader for reading and unmarshalling github commit data
type CommitReader interface {
	StreamCommits(dateRange dateutil.DateRange) <-chan CommitPage
	StreamAllCommits() <-chan CommitPage
}

// Default implementation of CommitReader
type DefaultCommitReader struct {
	client     http_util.HTTPClient
	repository string // eg. teradici/deploy
}

// NewCommitReader a CommitReader
func NewCommitReader(client http_util.HTTPClient, repository string) (CommitReader, error) {

	pattern := regexp.MustCompile("[A-Za-z0-9]+/[A-Za-z0-9]+")

	if !pattern.Match([]byte(repository)) {
		return nil, fmt.Errorf("repository %s does not match pattern [A-Za-z0-9]+/[A-Za-z0-9]+ (ie. {owner}/{repo})", repository)
	}

	return &DefaultCommitReader{
		client,
		repository,
	}, nil
}

// StreamAllCommits reads all of the commits from https://api.github.com/repos/{owner}/{repo}/commits
func (rdr *DefaultCommitReader) StreamAllCommits() <-chan CommitPage {
	return rdr.StreamCommits(dateutil.DateRange{})
}

// StreamCommits reads the commits from https://api.github.com/repos/{owner}/{repo}/commits
func (rdr *DefaultCommitReader) StreamCommits(dateRange dateutil.DateRange) <-chan CommitPage {

	out := make(chan CommitPage, 256)

	go func() {
		defer close(out)

		// curl \
		//   -H "Accept: application/vnd.github.v3+json" \
		//   https://api.github.com/repos/{owner}/{repo}/commits?per_page={#}&page={#}&since={YYYY-MM-DDTHH:MM:SSZ}&until={YYYY-MM-DDTHH:MM:SSZ}

		// ?since=	string	query
		// Only show notifications updated after the given time. This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.
		since := ""
		if dateRange.Start != nil {
			since = fmt.Sprintf("&since=%s", dateRange.Start.Format("2006-01-02T15:04:05Z"))
		}

		// ?until=	string	query
		// Only commits before this date will be returned. This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.
		until := ""
		if dateRange.End != nil {
			until = fmt.Sprintf("&until=%s", dateRange.End.Format("2006-01-02T15:04:05Z"))
		}

		per_page := 100 // 30 is the default, 100 is max
		page := 1       // pagination starts at 1
		for {

			// Assemble the URL
			url, _ := url.Parse(fmt.Sprintf("https://api.github.com/repos/%s/commits?per_page=%d&page=%d%s%s",
				rdr.repository,
				per_page,
				page,
				since,
				until,
			))

			// Assemble the HTTP request
			request := http.Request{
				Method: "GET",
				URL:    url,
				Header: http.Header{
					"Accept": []string{"application/vnd.github.v3+json"},
				},
			}

			// Perform the request with the configured client
			response, err := rdr.client.Do(&request)
			if err != nil {
				out <- CommitPage{
					Commits: nil,
					Err:     err,
				}
				return
			}

			// Only http status 200 is valid
			if response.StatusCode != http.StatusOK {
				out <- CommitPage{
					Commits: nil,
					Err:     fmt.Errorf("wrong status code: %d", response.StatusCode),
				}
				return
			}

			// Read body []byte
			body, err := io.ReadAll(response.Body)
			if err != nil {
				out <- CommitPage{
					Commits: nil,
					Err:     err,
				}
				return
			}

			// Unmarshal JSON body into structs
			var list []CommitData

			err = json.Unmarshal(body, &list)
			if err != nil {
				out <- CommitPage{
					Commits: nil,
					Err:     err,
				}
				return
			}

			// Return result
			out <- CommitPage{
				Commits: list,
				Err:     nil,
			}

			// Stop pagination if the page size is less than expected
			if len(list) < per_page {
				return
			} else {
				page += 1
				continue
			}

		}
	}()

	return out
}
