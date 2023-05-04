FROM golang:1.20

WORKDIR /go/src/app

COPY . .

RUN go build -v ./...

#RUN go test -v ./...

CMD ["go", "test", "-v", "./..."]

#CMD ["go", "run", "main.go"]
