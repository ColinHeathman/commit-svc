# commit-svc

This is a simple web application that reads data from Github's REST API and determines the people who have committed to a repository in a date range, as well as the total number of commits made during that date range.

## building 

### Docker

This project can be built easily using Docker.

`docker build .`

or

`docker-compose build`


### Go

To build this project directly without Docker.

Download the dependencies

`go mod download`

Build the application using the go command line.

`go build github.com/colinheathman/commit-svc/cmd`


### Running tests

Unit tests can be run within docker-compose.

`docker-compose run commit-svc go test -coverprofile=c.out -v ./...`

### Running the application

This application uses a Redis DB to cache Github's API responses. This can be run along side the main application with docker-compose.

`docker-compose up`

Then the application API can accessed using curl. jq is recommended for formatting JSON on the command line.

`curl -s 'localhost:8080/users?start=2019-06-01&end=2020-05-31' | jq .`

`curl -s 'localhost:8080/most-frequent?start=2019-06-01&end=2020-05-31' | jq .`

#### Configuration

The docker-compose.yml has a few environment variables set for the application.

`REDIS_URL=redis:6379` - The redis service to connect to (auth not supported)

`GITHUB_REPO=teradici/deploy` - The Github repository to read from

`LISTEN_ADDR=:8080` - The listen address for the HTTP server (https not supported)

`PROXY_PATH` can also be used to change the URI prefix for the application (eg. "/api/v1")
