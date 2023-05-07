FROM golang:1.20

WORKDIR /go/src/app

COPY . .

RUN go build -v ./...

ENTRYPOINT ["sh", "/go/src/app/entrypoint.sh"]