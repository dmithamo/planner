package main

import (
	"fmt"
	"net/http"
)

// notFoundErr renders the 404 page
func (a *application) notFoundErr(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("resource with url %s not found", r.URL.Path)
	a.errLogger.Printf("app run::err %v::%s", http.StatusNotFound, err)
	a.renderTemplate("notFound.page.tmpl", w, templateData{Error: err})
}
