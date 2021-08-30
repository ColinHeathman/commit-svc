FROM golang:1.17-buster

RUN mkdir -p /go/src/github.com/colinheathman/commit-svc

WORKDIR /go/src/github.com/colinheathman/commit-svc

COPY app/ app/
COPY cmd/ cmd/
COPY pkg/ pkg/

COPY test/ test/
ENV TEST_JSON_DIR="/go/src/github.com/colinheathman/commit-svc/test"

COPY [ "go.mod", "go.sum", "./"]

RUN go mod download

ENV GOOS=linux

RUN go build -a -o /usr/bin/commit-svc github.com/colinheathman/commit-svc/cmd

