package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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
