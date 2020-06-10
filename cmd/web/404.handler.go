package main

import (
	"net/http"
)

// notFoundErr renders the 404 page
func (a *application) notFoundErr(w http.ResponseWriter, r *http.Request) {
	a.errLogger.Printf("app run::err %v::resource with url `%s` not found", http.StatusNotFound, r.URL.Path)
	a.renderTemplate("notFound.page.tmpl", w, r)
}
