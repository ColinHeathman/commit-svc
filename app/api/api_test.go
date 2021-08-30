package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx/fxtest"
)

func TestAPI(t *testing.T) {

	lc := fxtest.NewLifecycle(t)

	// use a high numbered random port for test
	cfg := APIConfiguration{
		ProxyPath: "/",
		Addr:      ":31234",
	}

	api := NewAPI(lc, cfg)

	// set up test handler to check if it's called
	handled := false

	var testHandler http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
		handled = true
	}

	// register the test handler
	api.RegisterRoutes(func(router *mux.Router) {
		router.Handle("/test", testHandler)
	})

	// start the test application
	lc.Start(context.Background())
	defer lc.Stop(context.Background())

	// wait for the app to start - this is necessary
	<-time.After(10 * time.Millisecond)

	// test HTTP GET request to the test API
	resp, err := http.Get("http://localhost:31234/test")
	if err != nil {
		t.Error(err)
		return
	}

	// Ensure the handler was called
	if !handled {
		t.Fail()
		return
	}

	// Got the right status code
	if resp.StatusCode != http.StatusOK {
		t.Fail()
		return
	}
}
