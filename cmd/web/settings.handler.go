package main

import "net/http"

// settings handles requests to /settings
func (a *application) settings(w http.ResponseWriter, r *http.Request) {
	a.renderTemplate("settings.page.tmpl", nil, w, r)
}
