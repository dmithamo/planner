package main

import (
	"fmt"
	"log"
)

// checkErrsHelper checks for fatal errors, logs them, and shuts down the app
func checkFatalErrorsHelper(err error, logger *log.Logger, msg string) {
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s::%s", msg, err))
	}
}
