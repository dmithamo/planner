package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// listProjects handles requests to /projects
func (a *application) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := a.projects.SelectAll()
	if err != nil {
		panic(err)
	}

	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.templateData.Projects = projects
	a.renderTemplate("projects.page.tmpl", w, r)
}

// isValid reports validation errs in formData
func isValid(f *url.Values) (bool, formErrs) {
	errs := formErrs{}

	isValidString := func(value string) (bool, string) {
		value = strings.Trim(value, " \t \n")
		// empty values not allowed!
		if len(value) == 0 || value == "" {
			return false, "This should not be empty"
		}

		if utf8.RuneCountInString(value) < 4 {
			return false, "Too short! Make it at leat 4 chars long"
		}

		if utf8.RuneCountInString(value) > 144 {
			return false, "Too long! Keep it at 144 chars max"
		}

		return true, ""
	}

	// validate title
	if valid, err := isValidString(f.Get("title")); !valid {
		errs["title"] = err
	}
	// further restrict title to 100 chars
	if utf8.RuneCountInString(f.Get("title")) > 100 {
		errs["title"] = "Too long! Keep it at 100 chars max"
	}

	// validate description
	if valid, err := isValidString(f.Get("description")); !valid {
		errs["description"] = err
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, formErrs{}
}

// showCreateProjectForm handles requests at /projects/create
func (a *application) showCreateProjectForm(w http.ResponseWriter, r *http.Request) {
	// clear errs and form data if any
	a.templateData.FormData = &url.Values{}
	a.templateData.FormErrs = &formErrs{}

	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.renderTemplate("create.page.tmpl", w, r)
}

// showCreateProjectForm handles requests at /projects/create
func (a *application) createproject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")

	formData := url.Values{}
	formData.Set("title", title)
	formData.Set("description", description)
	isValid, errs := isValid(&formData)

	if !isValid {
		a.templateData.FormErrs = &errs
		a.templateData.FormData = &formData

		a.errLogger.Printf("app run::response::%v", http.StatusBadRequest)
		a.renderTemplate("create.page.tmpl", w, r)
		return
	}

	projectSlug, err := a.projects.Insert(strings.TrimSpace(title), strings.TrimSpace(description))

	if err != nil {
		a.errLogger.Printf("app run::create project::%s", err)

		a.templateData.FormData = &formData
		errString := strings.ToLower(err.Error())

		switch {
		case strings.Contains(errString, "duplicate"):
			a.templateData.FormErrs = &formErrs{"title": "This title is already in use"}

		case strings.Contains(errString, "projects.projectSlug"):
			a.templateData.FormErrs = &formErrs{"title": err.Error()}

		case strings.Contains(errString, "projects.description"):
			a.templateData.FormErrs = &formErrs{"description": err.Error()}

		default:
			a.templateData.FormErrs = &formErrs{"generic": err.Error()}
		}

		a.errLogger.Printf("app run::response::%v", http.StatusBadRequest)
		a.renderTemplate("create.page.tmpl", w, r)
		return
	}

	a.infoLogger.Printf("app run::response::%v", http.StatusPermanentRedirect)
	http.Redirect(w, r, fmt.Sprintf("/projects/slug/%v", projectSlug), http.StatusSeeOther)
}
