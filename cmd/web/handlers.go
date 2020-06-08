package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// landingPage handles requests to / <that is to say, root>
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFoundErr(w, r, fmt.Errorf("resource with URL `%v` not found", r.URL))
		return
	}

	a.renderTemplate("auth.page.tmpl", nil, w, r)
}

// ad handles requests to /list/listOfProjects
func (a *application) listOfProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := a.templateData.projects.SelectAll()
	if err != nil {
		panic(err)
	}

	a.renderTemplate("projects.page.tmpl", projects, w, r)
}

// singleProject handles requests to /projects/details/?id
func (a *application) singleProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["projectID"]
	projectID = strings.TrimSpace(projectID)

	if projectID == "" {
		err := errors.New("Invalid projectID")
		panic(err)
	}

	project, err := a.templateData.projects.SelectOne(projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFoundErr(w, r, fmt.Errorf("project with projectID `%s` not found", projectID))
			return
		}
		panic(err)
	}

	a.renderTemplate("project.page.tmpl", project, w, r)
}

// settings handles requests to /settings
func (a *application) settings(w http.ResponseWriter, r *http.Request) {
	a.renderTemplate("settings.page.tmpl", nil, w, r)
}
