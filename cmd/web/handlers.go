package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// landingPage handles requests to / <that is to say, root>
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	if r.URL.Path != "/" {
		a.notFoundErr(w, r, fmt.Errorf("resource with URL `%v` not found", r.URL))
		return
	}

	a.renderTemplate("auth.page.tmpl", nil, w, r)
}

// ad handles requests to /list/listOfProjects
func (a *application) listOfProjects(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	projects, err := a.templateData.projects.SelectAll()
	if err != nil {
		a.serverError(w, r, err)
	}

	a.renderTemplate("projects.page.tmpl", projects, w, r)
}

// singleProject handles requests to /projects/details/?id
func (a *application) singleProject(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	id := r.URL.Query().Get("id")
	id = strings.TrimSpace(id)

	if id == "" {
		err := errors.New("Invalid id")
		a.serverError(w, r, err)
		return
	}

	project, err := a.templateData.projects.SelectOne(id)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFoundErr(w, r, fmt.Errorf("project with id `%s` not found", id))
			return
		}
		a.serverError(w, r, err)
		return
	}

	a.renderTemplate("project.page.tmpl", project, w, r)
}

// settings handles requests to /settings
func (a *application) settings(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)
	a.renderTemplate("settings.page.tmpl", nil, w, r)
}
