package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/colinheathman/commit-svc/app/api"
	"github.com/colinheathman/commit-svc/app/endpoints/frequency"
	"github.com/colinheathman/commit-svc/app/endpoints/users"
	"github.com/colinheathman/commit-svc/pkg/github"
	"github.com/colinheathman/commit-svc/pkg/http_util"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"go.uber.org/fx"
)

func main() {
	flag.Parse()

	// Optional variable
	proxyPath, ok := os.LookupEnv("PROXY_PATH")
	if !ok {
		proxyPath = "/"
	}

	// Required variables
	initFail := false
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		glog.Error("REDIS_URL is not set")
		initFail = true
	}

	githubRepo, ok := os.LookupEnv("GITHUB_REPO")
	if !ok {
		glog.Error("GITHUB_REPO is not set")
		initFail = true
	}

	listenAddr, ok := os.LookupEnv("LISTEN_ADDR")
	if !ok {
		glog.Error("LISTEN_ADDR is not set")
		initFail = true
	}
	if initFail {
		glog.Exit("Initialization failure")
	}

	app := fx.New(
		// Provide all the constructors we need, which teaches Fx how we'd like to initialize our application
		// Constructors are called lazily, so this block doesn't do much on its own.
		fx.Provide(

			// Redis client
			func() *redis.Client {
				return redis.NewClient(&redis.Options{
					Addr: redisURL,
					DB:   0, // use default DB
				})
			},

			// Caching HTTP client
			func(rdb *redis.Client) http_util.HTTPClient {

				return http_util.NewCachingClient(
					rdb,
					http.DefaultClient,
					context.Background(),
				)
			},

			// Github commits reader
			func(client http_util.HTTPClient) (github.CommitReader, error) {
				return github.NewCommitReader(client, githubRepo)
			},

			// API configuration
			func() api.APIConfiguration {
				return api.APIConfiguration{
					ProxyPath: proxyPath,
					Addr:      listenAddr,
				}
			},
			// API
			api.NewAPI,

			// Endpoints for API
			users.NewEndpoint,
			frequency.NewEndpoint,
		),

		// Since constructors are called lazily, we need some invocations to
		// kick-start our application. Calling it requires Fx to build those
		// types using the constructors above.
		fx.Invoke(
			// Register Endpoints with the API
			users.RegisterRoutes,
			frequency.RegisterRoutes,
		),
	)

	// Start the application
	app.Run()

	// Block here with <-app.Done()
	<-app.Done()

}
