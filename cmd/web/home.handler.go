package main

import (
	"net/http"
)

// landingPage handles requests to / <that is to say, root>
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	a.renderTemplate("auth.page.tmpl", nil, w, r)
}
