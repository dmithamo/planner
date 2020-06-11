package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dmithamo/planner/pkg/form"
)

// listProjects handles requests to /projects
func (a *application) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := a.projects.SelectAll()
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.renderTemplate("projects.page.tmpl", w, r, templateData{Projects: projects})
}

// showCreateProjectForm handles requests at /projects/create
func (a *application) showCreateProjectForm(w http.ResponseWriter, r *http.Request) {
	// clear errs and form data if any
	a.infoLogger.Printf("app run::response::%v", http.StatusOK)
	a.renderTemplate("create.page.tmpl", w, r, templateData{Form: initialForm})
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

	formValues := url.Values{}
	formValues.Set("title", title)
	formValues.Set("description", description)
	requiredFields := form.Required{"title": {2, 100}, "description": {2, 300}}

	form := a.templateData.Form.New(formValues, requiredFields)
	if !form.IsValid() {
		a.errLogger.Printf("app run::response::%v", http.StatusBadRequest)
		a.renderTemplate("create.page.tmpl", w, r, templateData{Form: form})
		return
	}

	projectSlug, err := a.projects.Insert(strings.TrimSpace(title), strings.TrimSpace(description))
	if err != nil {
		a.errLogger.Printf("app run::create project::%s", err)

		errString := strings.ToLower(err.Error())
		switch {
		case strings.Contains(errString, "duplicate"):
			form.ValidationErrs.Add("title", "This title is already in use")

		case strings.Contains(errString, "projects.projectSlug"):
			form.ValidationErrs.Add("title", err.Error())

		case strings.Contains(errString, "projects.description"):
			form.ValidationErrs.Add("description", err.Error())

		default:
			form.ValidationErrs.Add("generic", err.Error())
		}

		a.errLogger.Printf("app run::response::%v", http.StatusBadRequest)
		a.renderTemplate("create.page.tmpl", w, r, templateData{Form: form})
		return
	}

	a.infoLogger.Printf("app run::response::%v", http.StatusPermanentRedirect)
	http.Redirect(w, r, fmt.Sprintf("/projects/slug/%v", projectSlug), http.StatusSeeOther)
}
