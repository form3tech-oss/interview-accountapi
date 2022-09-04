FROM golang:1.19-alpine

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Set working directory
WORKDIR /go/src/github.com/jsebasct/account-api-lib/

# Run tests
CMD CGO_ENABLED=0 go test ./...