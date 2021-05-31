FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY cmd/app /app/

CMD ["go", "test", "-v", "./..."]
