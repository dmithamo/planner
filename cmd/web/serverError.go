package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// errTrace := fmt.Sprintf("%v\n%v", err.Error(), string(debug.Stack()))
	errTrace := fmt.Sprintf("%v", err.Error())

	ts, err := template.ParseFiles([]string{
		"./views/html/error.page.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeError := w.Write([]byte(insertErrMessage(err)))

		if writeError != nil {
			a.errLogger.Fatal("error rendering error page: ", writeError)
		}
		return
	}

	if err := ts.Execute(w, errTrace); err != nil {
		a.errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeError := w.Write([]byte(insertErrMessage(err)))
		if writeError != nil {
			a.errLogger.Fatal("error rendering error page: ", writeError)
		}
		return
	}

	a.errLogger.Println(errTrace)
}

// insertErrMessage is a helper for inserting errors
func insertErrMessage(err error) string {
	return fmt.Sprintf(`{"msg":"something went wrong: %v"}`, err)
}
