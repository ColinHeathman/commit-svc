package api

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

// InternalServerError returns error code 500 on the given HTTP response
func InternalServerError(response http.ResponseWriter, err error) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(response, "{\"error\": \"%v\"}", err)
	glog.Errorf("Internal Server Error: %v", err)
}

// NotFound returns error code 404 on the given HTTP response
func NotFound(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(response, "{}")
}
