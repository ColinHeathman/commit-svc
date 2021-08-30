package users

import (
	"net/http"
	"time"

	"github.com/colinheathman/commit-svc/app/api"
	"github.com/colinheathman/commit-svc/pkg/dateutil"
	"github.com/colinheathman/commit-svc/pkg/github"
	"github.com/gorilla/mux"
)

type Endpoint struct {
	jsonResponder api.JSONResponder
	reader        github.CommitReader
}

func RegisterRoutes(api api.API, enp *Endpoint) {
	api.RegisterRoutes(func(router *mux.Router) {
		// Map the JSON responder to /users GET requests
		router.HandleFunc("/users", enp.jsonResponder.ServeHTTP).Methods("GET")
	})
}

func NewEndpoint(reader github.CommitReader) *Endpoint {
	return &Endpoint{
		api.JSONResponder{
			// From an http request, return a JSON marshallable object
			ModelFromRequest: func(r *http.Request) (interface{}, error) {

				// start and end dates
				var start *time.Time
				var end *time.Time

				q := r.URL.Query()

				// Replace start_date with one from query parameter if it exists
				sd, ok := q["start"]
				if ok {
					var err error
					parsedTime, err := time.Parse("2006-01-02", sd[0])
					if err != nil {
						return nil, err
					}
					start = &parsedTime
				}
				// Replace end_date with one from query parameter if it exists
				ed, ok := q["end"]
				if ok {
					var err error
					parsedTime, err := time.Parse("2006-01-02", ed[0])
					if err != nil {
						return nil, err
					}
					end = &parsedTime
				}

				dateRange := dateutil.NewDateRange(start, end)

				return CalculateUsers(reader, dateRange)
			},
		},
		reader,
	}
}
