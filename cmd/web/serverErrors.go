package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError renders server errors in a pop up.
func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var errTrace error
	// limit error info output for production env
	if a.isDevEnv {
		errTrace = fmt.Errorf("%v\n%v", err.Error(), string(debug.Stack()))
	} else {
		errTrace = fmt.Errorf("%v", err.Error())
	}

	a.renderTemplate("serverError.page.tmpl", w, r, templateData{Error: errTrace})
	a.errLogger.Printf("app run::err %v::%s", http.StatusInternalServerError, errTrace)
}
