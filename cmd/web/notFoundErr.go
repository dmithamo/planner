package main

import (
	"net/http"
	"text/template"
)

func (a *application) notFoundErr(w http.ResponseWriter, r *http.Request, notFoundError error) {
	ts, err := template.ParseFiles([]string{
		"./views/html/notFound.page.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.serverError(w, r, err)
		return
	}

	err = ts.Execute(w, notFoundError)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.errLogger.Println(http.StatusNotFound)
}
