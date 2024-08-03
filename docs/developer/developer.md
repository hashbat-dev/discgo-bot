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

  ## Coding Standards

  ## Pull requests 