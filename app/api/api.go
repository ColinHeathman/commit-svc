package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// API used for serving api endpoints through HTTPS
type API interface {
	RegisterRoutes(cb func(router *mux.Router))
}

// DefaultAPI default implementation of API
type DefaultAPI struct {
	router *mux.Router
	srv    *http.Server
}

// APIConfiguration config parameter
type APIConfiguration struct {
	ProxyPath string
	Addr      string
}

// NewAPI initializes an API
func NewAPI(
	lc fx.Lifecycle,
	cfg APIConfiguration,
) API {

	// Application router
	router := mux.NewRouter().PathPrefix(cfg.ProxyPath).Subrouter()
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// Hook the API to start and stop with the Fx application
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			// HTTPS
			go func() {

				paths := ""
				router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
					tmp, _ := route.GetPathTemplate()
					paths = fmt.Sprintf("%s\n %s", paths, tmp)
					return nil
				})

				glog.Infof("Starting up HTTP server @ %s%s", cfg.Addr, paths)

				srv.ListenAndServe()
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					// unexpected error. port in use?
					glog.Fatalf("ListenAndServe(): %v", err)
				}
				glog.Infof("HTTP server shutting down")
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return &DefaultAPI{
		router,
		srv,
	}
}

// RegisterRoutes registers the API paths to the given router
func (api *DefaultAPI) RegisterRoutes(cb func(router *mux.Router)) {
	// API paths
	cb(api.router)
}
