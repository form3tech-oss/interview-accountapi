Mateusz Gepert new to golang.

 - started task with narrowing down how to make http call with root lib from go (found that `net/http` and `encoding/json` will do the job). Since I was completetly new I started doing all in one big main function (not so bit around 150 lines TBH :))
 - decided to do simple split of application to separate model and actual implemetation (keeping in mind that actual API is an implementation detail also I introduced the level of abstraction for form3 accounts resources) - added also unit tests. I used here native testing lib and also one for assert and another one to mock the rest http calls to mimic API as close as possible. Unit tests are simple and covers only the happy path scenario (as an improvement I would recommend to make that tests more deep handle different error codes, malformed response bodies etc.)
 - at this time I found minimum scenario covered so I've enriched the docker-compose file and provided basic Dockerfile. I've prepared basic end to end test to check if all good with connection etc. (that might be tested with `docker-compose up --build`)
 - since I've seen in API swagger documentation that there are some validation done decided to prepare simple validation framework and implemented two validations one for BankID second one for BIC with relation to country field in account (that also might be improved but the validation logic and validator is present only new validators needs to be added and injected to AccountValidator which is a kind of a proxy for all Validators). To be able to use that validator I needed to create a simple decorator for `AccountRespository` that will trigger validation before delegating request to actual API (I found better to check that before making http request)
  - at the end I've decided to restructure the whole app and provide simple `handler` file with above decorator - it might be treated as and entrypoint and factory method for the whole lib

 Cut corners:
  - I decided to add limited amount of tests (just to get the idea how it works and how to extend that in future)
  - `Account` model is shared as contract to library and also as entry point to API (ideally it might be decoupled but for such small project I don't find it necessary). Also in real production ready env I would provide the simple fluent builder for account to be able to create them easily (possible driven by the actual account type and limitations based on country)
  - package structure I'm not fully familiar with go package management (I'm java origin and it was quite hard for me to understand how to best split the funcionality accros the lib) Ideally I would like to have abstracted entry point with clear definition of contract and provide simple factories to create dedicated client of library. Partialy it can be achived with `handler.NewForm3AccountHandler` which wraps and hides all implementation details only what need to be passed is url (which probably might be better abstracted)

 



# Form3 Take Home Exercise

Engineers at Form3 build highly available distributed systems in a microservices environment. Our take home test is designed to evaluate real world activities that are involved with this role. We recognise that this may not be as mentally challenging and may take longer to implement than some algorithmic tests that are often seen in interview exercises. Our approach however helps ensure that you will be working with a team of engineers with the necessary practical skills for the role (as well as a diverse range of technical wizardry). 

## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](http://api-docs.form3.tech/api.html#organisation-accounts) for information on how to interact with the API. Please note that the fake account API does not require any authorisation or authentication.

If you encounter any problems running the fake account API we would encourage you to do some debugging first,
before reaching out for help.

## Submission Guidance

### Shoulds

The finished solution **should:**
- Be written in Go.
- Be a client library suitable for use in another software project.
- Implement the `Create`, `Fetch`, and `Delete` operations on the `accounts` resource.
- Be well tested to the level you would expect in a commercial environment.
- Contain documentation of your technical decisions.
- Be simple and concise.
- Have tests that run from `docker-compose up` - our reviewers will run `docker-compose up` to assess if your tests pass.

### Should Nots

The finished solution **should not:**
- Use a code generator to write the client library.
- Use (copy or otherwise) code from any third party without attribution to complete the exercise, as this will result in the test being rejected.
- Use a library for your client (e.g: go-resty). Libraries to support testing or types like UUID are fine.
- Implement client-side validation.
- Implement an authentication scheme.
- Implement support for the fields `data.attributes.private_identification`, `data.attributes.organisation_identification`
  and `data.relationships`, as they are omitted in the provided fake account API implementation.
- Have advanced features, however discussion of anything extra you'd expect a production client to contain would be useful in the documentation.
- Be a command line client or other type of program - the requirement is to write a client library.
- Implement the `List` operation.
> We give no credit for including any of the above in a submitted test, so please only focus on the "Shoulds" above.

## How to submit your exercise

- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, copy the `docker-compose` from this repository
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) @form3tech-interviewer-1 to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team

## License

Copyright 2019-2021 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
