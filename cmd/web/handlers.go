package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// landingPage handles requests to /
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	if r.URL.Path != "/" {
		a.notFoundErr(w, r, fmt.Errorf("resource with URL `%v` not found", r.URL))
		return
	}

	ts, err := template.ParseFiles([]string{
		"./views/html/auth.page.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.errLogger.Println(err)
		a.serverError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.serverError(w, r, err)
		return
	}
	a.infoLogger.Println(http.StatusOK)
}

// ad handles requests to /list/listOfProjects
func (a *application) listOfProjects(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	projects, err := a.projects.SelectAll()
	if err != nil {
		a.serverError(w, r, err)
	}

	ts, err := template.ParseFiles([]string{
		"./views/html/projects.list.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.serverError(w, r, err)
		return
	}

	if err := ts.Execute(w, projects); err != nil {
		a.serverError(w, r, err)
		return
	}

	a.infoLogger.Println(http.StatusOK)
}

// singleProject handles requests to /projects/?id
func (a *application) singleProject(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	id := r.URL.Query().Get("id")
	id = strings.TrimSpace(id)

	if id == "" {
		err := errors.New("Invalid id")
		a.serverError(w, r, err)
		return
	}

	ts, err := template.ParseFiles([]string{
		"./views/html/projects.single.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.serverError(w, r, err)
		return
	}

	project, err := a.projects.SelectOne(id)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFoundErr(w, r, fmt.Errorf("project with id `%s` not found", id))
			return
		}

		a.serverError(w, r, err)
		return
	}

	if err := ts.Execute(w, project); err != nil {
		a.serverError(w, r, err)
		return
	}

	a.projects.Insert("I am the stone that the builder refused", "No description")
	a.infoLogger.Println(http.StatusOK)
}

// settings handles requests to /settings
func (a *application) settings(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	ts, err := template.ParseFiles([]string{
		"./views/html/settings.page.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.serverError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.serverError(w, r, err)
		return
	}
	a.infoLogger.Println(http.StatusOK)
}
