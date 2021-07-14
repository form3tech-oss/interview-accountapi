#!/bin/bash

export GOPATH=$GOPATH:$PWD
go build -ldflags="-X 'api.URL=http://localhost:8080'" main.go