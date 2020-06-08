# Projects

## A Go Monolith

### Description

Projects is an effort at learning Web Dev with Go.

It is also a tool I mean to use to curate projects I am working on, and projects I project I will one day work on, for fun and learning.

### Run, app, run

To test, execute the following commands in a Terminal (bash or similar):

```bash
#!/bin/bash
git clone <this-repo>
cd planner
git checkout develop
go run cmd/web/* -p=:3001 -db='dbuser:dbpassword@[host]/dbname  -rdb=true -ltf=true -dev=true
```

If all went well, the app should load the app at the provided port on localhost

### Flags definition

|Flag|Description|Default|
|:------|:-----|----:|
`p`| port at which the app will run | `:3001`
`db`| DSN of the database | `REQUIRED`
`rdb`| setting this to true will create the database tables, dropping them if they already exist | `false`
`ltf`| setting this to true means all logging is done to a file in the root dir of the project | `false`
`dev`|sets the env to `development`. Not too useful, beyond setting the depth of info  logged by error loggers | `false`

### The Stack

This project uses the following open source technologies

- [The Go programming Language](https://golang.org/)
- [MySQL Database](https://www.mysql.com/)
- [Alice Middleware Chaining Tool](https://github.com/justinas/alice)

### TODO

Populate this README more verbosely.
