package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError renders server errors in a special page.
func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	errTrace := ""
	// limit error info output for production env
	if a.isDevEnv {
		errTrace = fmt.Sprintf("%v\n%v", err.Error(), string(debug.Stack()))
	} else {
		errTrace = fmt.Sprintf("%v", err.Error())
	}

	ts, ok := a.templates["serverError.page.tmpl"]
	if !ok {
		a.errLogger.Fatal("app run::server error page::fail ", errors.New("template not found"))
		return
	}

	if err := ts.Execute(w, errTrace); err != nil {
		a.errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeError := w.Write([]byte(insertErrMessage(err)))
		if writeError != nil {
			a.errLogger.Fatal("error rendering error page: ", writeError) // this is funny.
		}
		return
	}

	a.errLogger.Println(errTrace)
}

// insertErrMessage is a helper for inserting errors
func insertErrMessage(err error) string {
	return fmt.Sprintf(`{"msg":"something went wrong: %v"}`, err)
}
