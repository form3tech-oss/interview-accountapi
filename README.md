# Form3 Take Home Exercise

## Instructions

This exercise has been designed to be completed in 2-4 hours. The goal of this exercise is to write a client library in Go to access our fake [account API](http://api-docs.form3.tech/api.html#organisation-accounts) service. The library should be able to perform all of the operations outlined in the [documentation](http://api-docs.form3.tech/api.html#organisation-accounts)

### Should
- Client library should be written in Go
- Focus on writing full-stack tests that cover the full range of expected and unexpected use-cases
 - Tests can be written in Go idomatic style or in BDD style. Make sure tests are easy to read
 - Engineers at Form3 generally favor BDD style tests or idomatic Go tests for this kind of task.
 
- Docker-compose
 - Add your solution to the provided docker-compose file
 - We should be able to `docker-compose up` and see your tests run against the provided account API service 

 
### Please don't
- Please do not use a code generator to write the client library


## How to submit your exercise
- Create a private repository, copy the `docker-compose` from this repository
- Email our recruitment team to let them know you have finished the assignment, they will then ask you to add to reviewers
