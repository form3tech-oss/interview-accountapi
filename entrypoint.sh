#!/bin/sh
export TEST_CONTEXT="container"

echo $TEST_CONTEXT

go test -v ./... -cover