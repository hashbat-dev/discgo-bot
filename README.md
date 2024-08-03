# discgo-bot

## To-Do
* Implement the "admin only" add gif property

## Push Request Checklist
1. Update the Version in version.go
2. List your changes in the string array variable in version.go
3. Push


## Automation Testing

This is handled by GitHub actions making use of the Dockerfile we already wrote to spin up a container
and run a suite of tests. The intention for this is for it to run the same kind of checks we run locally ahead
of looking to get changes solidifed. This includes:
- Linting
- Go test calls
- DB query schema validations? (TBD)

From the developers point of view, this change also includes decoupling local testing from the production DB, to prevent 
this causing any damage to already stored data (we should also look to isolate local instances from connecting to it entirely).

To get started, run the following:
`./scripts/initdb.sh &`
Which will spin up a docker container with an instance of our database, which we then seed with basic values as part of the startup
of the test suite.



