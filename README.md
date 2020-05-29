# Projects

## A Go Monolith

### Description

Projects is an effort at learning Web Dev with Go.

It is also a tool I mean to use to curate projects I am working on, and projects I project I will one day work on, for fun and learning.

### Run, app, run

To test, in a Terminal (bash or similar), do:

```bash
#!/bin/bash
git clone <this-repo>
cd planner
git checkout develop
go run cmd/web/* -p=:1024 -db='dbuser:dbpassword@[host]/dbname  -rdb=true -ltf=true -dev=true

# flags
# -p - this is the port at which the app will run. Default is :3001
# -db - this is the DSN of the database. The project uses a MySQL db
# -rdb - setting this to true will create the database tables, dropping them if they already exist
# -ltf - setting this to true means all logging is done to a file. Requires that a directory /logs exist in the root dir of the project
# -dev - sets the env to dev. Not too useful, beyond setting the depth of info  logged by error loggers
```

If all went well, the app should load the app at the provided port on localhost

### TODO

Populate this README more verbosely.
