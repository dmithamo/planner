package main

import "net/http"

// ad handles requests to /list/listProjects
func (a *application) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := a.templateData.projects.SelectAll()
	if err != nil {
		panic(err)
	}

	a.renderTemplate("projects.page.tmpl", projects, w, r)
}
