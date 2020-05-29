package main

import (
	"net/http"
)

// notFoundErr renders the 404 page
func (a *application) notFoundErr(w http.ResponseWriter, r *http.Request, notFoundError error) {
	a.renderTemplate("notFound.page.tmpl", nil, w, r)
	a.errLogger.Println(http.StatusNotFound)
}
