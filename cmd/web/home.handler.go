package main

import (
	"net/http"
)

// landingPage handles requests to /auth
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.renderTemplate("auth.page.tmpl", w, r)
}
