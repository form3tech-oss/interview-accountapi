# Form3 Take Home Exercise

Engineers at Form3 build highly available distributed systems in a microservices environment. Our take home test is designed to evaluate real world activities that are involved with this role. We recognise that this may not be as mentally challenging and may take longer to implement than some algorithmic tests that are often seen in interview exercises. Our approach however helps ensure that you will be working with a team of engineers with the necessary practical skills for the role (as well as a diverse range of technical wizardry). 

## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](https://www.api-docs.form3.tech/api/tutorials/getting-started/create-an-account) for information on how to interact with the API. Please note that the fake account API does not require any authorisation or authentication.

A mapping of account attributes can be found in [models.go](./models.go). Can be used as a starting point, usage of the file is not required.

If you encounter any problems running the fake account API we would encourage you to do some debugging first,
before reaching out for help.

## Submission Guidance

### Shoulds

The finished solution **should:**
- Be written in Go.
- Use the `docker-compose.yaml` of this repository.
- Be a client library suitable for use in another software project.
- Implement the `Create`, `Fetch`, and `Delete` operations on the `accounts` resource.
- Be well tested to the level you would expect in a commercial environment. Note that tests are expected to run against the provided fake account API.
- Run the tests when `docker-compose up` is run - our reviewers will run `docker-compose up` and expect to see the test results in the output.
- Be simple and concise.

### Should Nots

The finished solution **should not:**
- Use a code generator to write the client library.
- Use (copy or otherwise) code from any third party without attribution to complete the exercise, as this will result in the test being rejected.
    - **We will fail tests that plagiarise others' work. This includes (but is not limited to) other past submissions or open-source libraries.**
- Use a library for your client (e.g: go-resty). Anything from the standard library (such as `net/http`) is allowed. Libraries to support testing or types like UUID are also fine.
- Implement client-side validation.
- Implement an authentication scheme.
- Implement support for the fields `data.attributes.private_identification`, `data.attributes.organisation_identification`
  and `data.relationships` or any other fields that are not included in the provided `models.go`, as they are omitted from the provided fake account API implementation.
- Have advanced features, however discussion of anything extra you'd expect a production client to contain would be useful in the documentation.
- Be a command line client or other type of program - the requirement is to write a client library.
- Implement the `List` operation.
> We give no credit for including any of the above in a submitted test, so please only focus on the "Shoulds" above.

## How to submit your exercise

- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, by copying all files you deem necessary for your submission
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) [@form3tech-interviewer-1](https://github.com/form3tech-interviewer-1) to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team

## License

Copyright 2019-2023 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

# Notes

## usage
Following the guidelines and the documentation I have implemented this simple client Library for the form3 api:

The account api client can be created as a struct or using the following constructor:
```
cfg := account.Config{
		BaseUrl: "localhost:8080",
		Version: "v1",
    MaxRetries: 3,
    Wait: 500,
	}

ac := account.NewAccountClient(&cfg)
```
The MaxRetries and Wait are the configuration for the Exponential back off strategy. The implementation is just a literal translation of the suggested pseudocode at form3 API documentation.

The client library provides 3 functions:
```
ctx := context.context.Background()

res, err := client.FetchAccount(ctx, id)

err = client.DeleteAccount(ctx, ir, version)

res, err = client.FetchAccount(ctx, id)

res, err = client.CreateAccount(ctx, &account.AccountData{})

```
Api responses with errors can be handled specifically as account.ErrorResponse.
```
	res, err = client.CreateAccount(ctx, createRequest())
	
  var errorResponse *account.ErrorResponse
  if err == nil {
    ...
	} else if errors.As(err, &errorResponse) {
		fmt.Println(errorResponse)
	} else {
		fmt.Println(err)
	}

```

Each function receives a context that can be used for cancellation, timeout, and error handling.
```
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.FetchAccount(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")

```
The context will be used for the request and any needed retry. This means that the original request and timeouts, retries or back offs should happen within the original context timeout.

The client library can be fine tuned via the http.Client to set a time out for individual requests and many other details.

## Running the code
To run the tests using docker compose:

```
# docker compose up
```

This will run the original containers **accountapi**, **postgresql**, **vault** and **integration_tests**. The last container is created with the Dockerfile and copies source files, build the go project and run the tests. Integration tests will use the internal urls to access other containers. 

To make changes in the code and re run the tests with docker compose:
```
docker compose up --build integration_tests
```

To run the tests locally and outside the container:

```
# docker compose up accountapi postgresql vault &

# go test -v ./...
```
The containers will map the service's ports to localhost and the test will use localhost url to access to the containers.

## Tests cases
There are three test case files:
- `client_test.go`: unit tests mocking a http.Server validating that requests are well built and responses well handled.
- `integration_tests.go`: create/fetch/delete account methods using the docker containers with the **accountapi** and the **postgresql**.
- `back_off_test.go`: unit tests to validate retry, back off, cancelling, timeout on different scenarios.

For the integration tests a direct connection with the database was chosen for **set up** and **tear down** scenarios. This has the disadvantage of coupling the test cases to the internal implementation of the the **accountapi**: any database changes can break the test cases. Another option for this would be using the same API services to create and re create the scenarios.

I started with the db approach to have more end to end control. In the end just implemented a simple `DELETE FROM "Account"` for clean up and an `INSERT INTO "Account" (...) VALUES (...)` to create some accounts. It was simple enough leave it this way and change in the future if needed.

As stated above, a simple configuration was added to the integration test to use a local configurations for development or other options for CI/CD etc.

Last, I've tried to add some simple stress tests here to get a 429 and assert some info about the retry and back off strategy, but could not stress the server without getting first a 500 error caused by postgres amount of connections before reaching any throttling error.

## external dependencies
The client was implemented using only the golang standard library. The dependencies at the `go.mod` are exclusive for test cases: testify and pq.
```
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.8.2
```