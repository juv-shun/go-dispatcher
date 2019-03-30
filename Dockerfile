FROM golang:1.12-alpine

LABEL maintainer="Shun Fukusumi(shun.fukusumi@gmail.com)"

RUN apk update \
    && apk add --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep

ADD ./src/ /go/src/github.com/juv-shun/go-worker_template/src/
WORKDIR /go/src/github.com/juv-shun/go-worker_template/src/
RUN dep ensure && go build -o main .
CMD ["./main"]
