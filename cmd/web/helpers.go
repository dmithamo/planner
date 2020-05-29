package main

import (
	"fmt"
	"log"
	"net/http"
)

// checkErrsHelper checks for fatal errors, logs them, and shuts down the app
func checkErrorsHelper(err error, logger *log.Logger, msg string) {
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s::%s", msg, err))
	}
}

// renderTemplate renders tempates used in handlers
func (a *application) renderTemplate(templateName string, data interface{}, w http.ResponseWriter, r *http.Request) {
	ts, ok := a.templates[templateName]
	if !ok {
		a.serverError(w, r, fmt.Errorf("app run::%s::template not found", templateName))
		return
	}

	if err := ts.Execute(w, data); err != nil {
		a.serverError(w, r, fmt.Errorf("app run::%s::template err::%s", templateName, err))
		return
	}

	a.infoLogger.Println(http.StatusOK)
}
