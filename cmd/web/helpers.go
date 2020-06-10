package main

import (
	"fmt"
)

// checkErrsHelper checks for fatal errors, logs them, and shuts down the app
func checkFatalErrorsHelper(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s::%s", msg, err))
	}
}
