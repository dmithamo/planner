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
	projectSlug := mux.Vars(r)["projectSlug"]
	projectSlug = strings.TrimSpace(projectSlug)

	if projectSlug == "" {
		err := errors.New("Invalid projectSlug")
		panic(err)
	}

	project, err := a.projects.SelectOne(projectSlug)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFoundErr(w, r)
			return
		}
		panic(err)
	}

	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.renderTemplate("project.page.tmpl", w, templateData{Project: project})
}
