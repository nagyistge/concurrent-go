FROM golang:1.4.2
MAINTAINER peter.edge@gmail.com

RUN mkdir -p /go/src/github.com/peter-edge/go-concurrent
ADD . /go/src/github.com/peter-edge/go-concurrent/
WORKDIR /go/src/github.com/peter-edge/go-concurrent
