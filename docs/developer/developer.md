# Welcome to the developer docs!

### This is where we set out coding standards, logic flow, etc.

## TOC
1. [Set up to run locally](#set-up-to-run-locally)
2. [Control flow](#control-flow)
    1. [General Control Flow](#general-control-flow)
3. [Contribution standards](#coding-standards)
    1. [Branches](#branches)
    1. [Coding Standards](#coding-standards)
    1. [Pull Requests](#pull-requests)
1. [Automation Testing](#automation-testing)

---
---

# Set up to run locally

# Control Flow
## General control flow
```mermaid 
flowchart TB
  node_1["Command Received"]
  node_2["Decode Command"]
  node_3["Determine permissions"]
  node_4["Call handler to get return message"]
  node_5["Return message"]
  node_6(["Dictionary for permissions, \nhave an int table of perms. \nIf user perm #gt;= required perm then allow"])
  node_7(["Determine where and how to return.\nShould it delete original msg?\nShould it spawn a thread?"])
  node_1 --> node_2
  node_2 --> node_3
  node_3 --> node_4
  node_4 --> node_5
  node_3 --> node_6
  node_5 --> node_7
  ```

# Contribution Standards
## Branches
Branches to add new features should be called ```feature/X```

Branches to fix identified issues should be called ```issue/X```

Branches specifically for documentation updates should be called ```docs/X```
## Coding Standards
Code should be clean, readible, and commented!

## Pull requests
  1. Update the Version in version.go
  1. List your changes in the string array variable in version.go
  1. Push
  1. Make a pull request detailing your changes. Be Desciptive!

**Note: Pull requests can't be merged if they don't pass all tests!**

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