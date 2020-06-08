package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// viewProject handles requests to /projects/details/?id
func (a *application) viewProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["projectID"]
	projectID = strings.TrimSpace(projectID)

	if projectID == "" {
		err := errors.New("Invalid projectID")
		panic(err)
	}

	project, err := a.templateData.projects.SelectOne(projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFoundErr(w, r)
			return
		}
		panic(err)
	}

	a.renderTemplate("project.page.tmpl", project, w, r)
}

func (a *application) createProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		a.renderTemplate("create.page.tmpl", nil, w, r)
	}
}
