package main

import (
	"fmt"
	"net/http"
)

// notFoundErr renders the 404 page
func (a *application) notFoundErr(w http.ResponseWriter, r *http.Request) {
	a.renderTemplate("notFound.page.tmpl", fmt.Errorf("resource with url `%s` not found", r.URL.Path), w, r)
}
