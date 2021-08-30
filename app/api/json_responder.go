package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

// JSONResponder is a generic http handler for responding to a request with a JSON object
type JSONResponder struct {
	ModelFromRequest func(*http.Request) (interface{}, error)
}

// ServeHTTP so that JSONResponder implements http.Handler
func (rspdr *JSONResponder) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {

	// Use the configured function to get an interface{} to marshal
	model, err := rspdr.ModelFromRequest(request)
	if err != nil {
		InternalServerError(response, err)
		return
	}

	// Marshal the JSON using the data
	dat, err := json.Marshal(model)
	if err != nil {
		err := fmt.Errorf("failed to marshal json: %v", err)
		InternalServerError(response, err)
		return
	}

	// Special case for Nil interface
	if string(dat) == "null" {
		NotFound(response)
		return
	}

	// Set response headers and status code
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.WriteHeader(http.StatusOK)
	// Write body
	_, err = response.Write(dat)
	if err != nil {
		glog.Errorf("failed to write response: %v", err)
	}
}
