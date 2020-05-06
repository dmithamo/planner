package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

func (a *application) handleError(w http.ResponseWriter, r *http.Request, err error) {
	errTrace := fmt.Sprintf("%v\n%v", err.Error(), string(debug.Stack()))

	ts, err := template.ParseFiles([]string{
		"./ui/html/error.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}

	if err := ts.Execute(w, errTrace); err != nil {
		a.errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}
}

// insertErrMessage is a helper for inserting errors
func insertErrMessage(err error) string {
	return fmt.Sprintf(`{"msg":"something went wrong: %v"}`, err)
}
