package main

import "net/http"

// settings handles requests to /settings
func (a *application) settings(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.renderTemplate("settings.page.tmpl", w, templateData{})
}
