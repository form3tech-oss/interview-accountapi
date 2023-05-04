FROM golang:1.20

WORKDIR /go/src/app

COPY . .

RUN go build -v ./...

CMD ["go", "test", "-v", "./..."]

