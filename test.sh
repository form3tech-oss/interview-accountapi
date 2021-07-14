#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

docker-compose up --detach
export GOPATH=$GOPATH:$PWD
GOTEST=$(go test -ldflags="-X 'api.URL=http://localhost:8080'" ./... -coverprofile cover.out)
echo "${GOTEST}"
RESULT=$(echo -n ${GOTEST} | awk '{print $NF}')
if [ $RESULT = "FAIL" ]
then 
    docker-compose stop
    echo -e "${RED}Some tests failed${NC}"
    exit -1
fi
go tool cover --html=cover.out -o cover.html
COVERAGE="$(go tool cover --func=cover.out | tail -1 | awk '{print $NF}')"
EXPECTED="100.0%"
if [ $COVERAGE = $EXPECTED  ]
then
    docker-compose stop
    echo -e "${GREEN}All tests passed with expected code coverage of $EXPECTED${NC}"
    exit 0
else
    docker-compose stop
    echo -e "${RED}All tests passed but coverage is only at $COVERAGE. Please bring it up to $EXPECTED${NC}"
    exit -1
fi