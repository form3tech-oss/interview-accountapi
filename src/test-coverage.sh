go test ${1:-./...} -coverprofile coverage.out 
go tool cover -html=coverage.out