FROM golang:alpine

WORKDIR /go/src
COPY . github.com/xdimgg/cheater
WORKDIR /go/src/github.com/xdimgg/cheater

RUN go mod vendor
RUN go build -o /bin/cheater
RUN apk del golang*

ENTRYPOINT /bin/cheater