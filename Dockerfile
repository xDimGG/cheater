FROM golang:alpine

WORKDIR /go/src
COPY . github.com/xdimgg/cheater
WORKDIR /go/src/github.com/xdimgg/cheater

RUN apk update
RUN apk upgrade
RUN apk add git curl --no-cache
RUN GO111MODULE=on go mod vendor
RUN go build -o /bin/cheater
RUN apk del golang*

ENTRYPOINT /bin/cheater