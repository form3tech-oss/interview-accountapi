# Aymeric Hermant - Beginner with Go

Thanks for taking the time to review my work!🤘🏻🎉  
This is a client for the form3 account API developed in Go.

It includes the `Create`, `Fetch`, and `Delete` operations on the `accounts` resource of the [Form3 API](http://api-docs.form3.tech/api.html#organisation-accounts) and some BDD style `Ginkgo` integration and unit tests.

## Usage

```go

// Account creation
id := "ad27e265-9604-4b4b-a0e5-3000ea9cc8d0"
version := "0"
accountData:=	`{
		"data":{
			"type": "accounts",
			"id": "`+id+`",
			"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			"version" :"0",
			"attributes": {
				"name": ["ahtest"],
				"country": "GB"
			}
		}
	}`



resp, err := form3client.CreateAccount(accountData)
if err != nil {
  log.Fatal(err)
}

// Account fetching
resp, err := form3client.FetchAccount(id)
if err != nil {
  log.Fatal(err)
}

// Account deletion
resp, err := form3client.DeleteAccount(id, version)
if err != nil {
  log.Fatal(err)
}

```

## Environment variables

| Environment variable   | Description                               |
| :--------------------- | :---------------------------------------- |
| `ACCOUNT_API_BASE_URL` | API host URL, e.g. https://api.form3.tech |

## Testing

Tests are written with `Gingko` in BDD style.  
They are 2 parts:

- Some integration tests which require the API
- Some unit tests

The unit tests are using [Gomega ghttp](https://pkg.go.dev/github.com/onsi/gomega/ghttp) to mock the BE server.
They can be launched with `docker-compose`:

```
docker-compose up
```

Do not forget to rebuild if you are making changes and are testing multiple times:

```
docker-compose up --build
```

It can also be tested with `ginkgo` if you already have going installed and an up and running docker image for your account API. For example, with the coverage and detailed view:

```
cd go
ginkgo -v --cover
```

Or with `go test`:

```
cd go
go test -v --cover
```

## Todos for production

- Run tests in CI
- Provide a way to use the `AccountData` struct as an input/output
- Use SSL
- Mock `HttpRequest` and `io.Reader` for a 100% test coverage vs 92.2% now
- Remove all the environment variables values from `docker-compose.yml` and use `Gitlab secrets`, `Github secrets` or an external `Vault` instead to retrieve them
- Separate integration tests and unit tests within 2 different test suites
