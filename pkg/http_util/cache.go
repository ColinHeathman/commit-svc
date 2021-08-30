package http_util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// ErrorHTTPClient client to return errors
type CachingClient struct {
	rdb          *redis.Client
	parentClient HTTPClient
	ctx          context.Context
}

// NewCachingClient initializes an HTTPClient
func NewCachingClient(
	rdb *redis.Client,
	parentClient HTTPClient,
	ctx context.Context,
) HTTPClient {
	return &CachingClient{
		rdb,
		parentClient,
		ctx,
	}
}

// Do executes an HTTP request possibly using the cache
func (client *CachingClient) Do(req *http.Request) (*http.Response, error) {

	val, err := client.rdb.Get(client.ctx, req.URL.String()).Result()
	if err == redis.Nil {
		// Update cache
		return client.Update(req)
	}
	if err != nil {
		return nil, err
	}
	// Return cached response
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(val)),
	}, nil
}

// Update executes an HTTP request and updates the cache with the result
func (client *CachingClient) Update(req *http.Request) (*http.Response, error) {

	// Use parent client to do request
	resp, err := client.parentClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Only works for status code OK
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code from API response: %v", resp.StatusCode)
	}

	// Read data from body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Put data in Redis for the request URL with a 2 minute timeout
	err = client.rdb.Set(client.ctx, req.URL.String(), data, 2*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	// Return cached response
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(data)),
	}, nil
}
