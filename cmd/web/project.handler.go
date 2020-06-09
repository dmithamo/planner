package main

import (
	"database/sql"
	"errors"
	"fmt"
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

	project, err := a.templateData.projects.SelectOne(projectSlug)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFoundErr(w, r)
			return
		}
		panic(err)
	}

	a.renderTemplate("project.page.tmpl", project, w, r)
}

type formData struct {
	title       string
	description string
}
type formErrs map[string]string

// createProject handles requests at /projects/create
func (a *application) createProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		a.renderTemplate("create.page.tmpl", nil, w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		a.clientError(w, r, err)
		return
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")

	formErrs := validateFormData(formData{title, description})
	if len(formErrs) > 0 {
		a.renderTemplate("create.page.tmpl", formErrs, w, r)
		return
	}

	projectSlug, err := a.templateData.projects.Insert(title, description)
	if err != nil {
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%v", projectSlug), http.StatusSeeOther)
}

// validateFormData validates form data
func validateFormData(formData formData) formErrs {
	return formErrs{}
}
