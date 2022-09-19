# Form3 Take Home Exercise - Submission by Gianni Massi

## Solution

The solution is provided as a go module (github.com/giannimassi/interview-accountapi/accountapi) exposing functions to instantiate a client to the Account API with options including host and custom http client (e.g. useful for customizing timeouts). The client provides methods txo create, fetch and delete an Account. The methods support providing a context for controlling cancellation and expect the type `AccountData`. This might not be the most ergonomic in many cases and I would expect some utility methods to be added as needed, not expecting a full `AccountData` struct but only the data required to make the call (making more clear what is expected).

### Testing

All tests are run via the docker compose command. Tests are organized as follows:

- unit tests: these can be found in the accountapi/accountapi_test.go and they test the behaviour of the functions is as expected, by providing a mocked backend;
- integration tests: these can be found accountapi/tests/account_api_test.go and they test that the outcome of calling these functions against the actual backend works.

End-to-end and acceptance tests might also be worth implementing, depending on the context where this API is deployed and who its users are (e.g. if it is to be exposed as a product of its own, acceptance tests would be desirable to make sure the product satisfies the business requirements, if it is to be used by a web page, a end-to-end test would also be good).

#### Running tests locally with docker-compose

The following command can be used to run the tests via docker-compose:
```bash
docker-compose up
````

This will start the backend docker image with its dependencies and then run both unit and integration tests. This is well suited for running on push in CI to make sure there are no regressions from changes.

Note: a different approach would be to not modify the docker-compose.yml but to add an override file that adds the test run on top of the existing docker compose file, separating the concerns (testing vs running locally). Assuming the original _docker-compose.yml_ did not get modified, and a second file called _docker-compose.tests.yml_ was created in the root dir with the following contents:

```yaml
services:
  accountapi:
    hostname: accountapi
  tests:
    image: golang:1.19.1
    depends_on:
      - accountapi
    environment:
      - ACCOUNTAPI_HOST=http://accountapi:8080
    volumes:
      - $PWD:/go/src/github.com/giannimassi/accountapi
    working_dir: /go/src/github.com/giannimassi/accountapi
    command: go test -count=1 -race ./...
```

Test would be run with the following command:

```bash
docker-compose -f docker-compose.yml -f docker-compose.tests.yml run tests
```
This might make sense if the first docker-compose.yml is used for other scopes as well (e.g. manual testing, qa, etc).

[What follows is the provided excercise brief for reference]

## Introduction
Engineers at Form3 build highly available distributed systems in a microservices environment. Our take home test is designed to evaluate real world activities that are involved with this role. We recognise that this may not be as mentally challenging and may take longer to implement than some algorithmic tests that are often seen in interview exercises. Our approach however helps ensure that you will be working with a team of engineers with the necessary practical skills for the role (as well as a diverse range of technical wizardry). 

## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](http://api-docs.form3.tech/api.html#organisation-accounts) for information on how to interact with the API. Please note that the fake account API does not require any authorisation or authentication.

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
- Be simple and concise.
- Have tests that run from `docker-compose up` - our reviewers will run `docker-compose up` to assess if your tests pass.

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

Copyright 2019-2022 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
