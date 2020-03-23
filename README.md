# Form3 Take Home Exercise

## Instructions
This exercise has been designed to be completed in 4-8 hours. The goal of this exercise is to write a client library 
in Go to access our fake [account API](http://api-docs.form3.tech/api.html#organisation-accounts) service. 

### Should
- Client library should be written in Go
- Document your technical decisions
- Implement the `Create`, `Fetch`, `List` and `Delete` operations on the `accounts` resource. Note that filtering of the List operation is not required, but you should support paging
- Ensure your solution is well tested to the level you would expect in a commercial environment. Make sure your tests are easy to read.
- If you encounter any problems running the fake accountapi we would encourage you to do some debugging first, 
before reaching out for help

#### Docker-compose
 - Add your solution to the provided docker-compose file
 - We should be able to run `docker-compose up` and see your tests run against the provided account API service 

### Please don't
- Use a code generator to write the client library
- Use a library for your client (e.g: go-resty). Only test libraries are allowed.
- Implement an authentication scheme

## How to submit your exercise
- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise 
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, copy the `docker-compose` from this repository
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) @form3tech-interviewer-1 to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team
