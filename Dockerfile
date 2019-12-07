FROM golang:1.13.5-stretch

ENV API_HOST ''
ENV API_PORT ''

WORKDIR /go/src/github.com/lioda/interview-accountapi

RUN apt-get update && apt-get install -y wait-for-it
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . /go/src/github.com/lioda/interview-accountapi
CMD wait-for-it "${API_HOST}:${API_PORT}" && go test ./... -v