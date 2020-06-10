package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

var customFuncs = template.FuncMap{
	"localizedDate": func(t time.Time) string { return t.Format("02 Jan 2006 at 15:00") },
}

// buildTemplatesCache builds a cache of the templates
func buildTemplatesCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		pageName := filepath.Base(page)
		ts, err := template.New(pageName).Funcs(customFuncs).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add partials to ts
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// add layouts to ts
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[pageName] = ts
	}

	return cache, nil
}

// renderTemplate renders tempates used in handlers
func (a *application) renderTemplate(templateName string, w http.ResponseWriter, r *http.Request) {
	ts, ok := a.templates[templateName]
	if !ok {
		panic(fmt.Errorf("app run::%s::template not found", templateName))
	}

	var data interface{}
	switch {
	case templateName == "serverError.page.tmpl":
		data = a.templateData.ServerErr

	case templateName == "notFound.page.tmpl":
		data = fmt.Errorf("resource with url `%s` not found", r.URL.Path)

	case templateName == "auth.page.tmpl":
		data = nil

	case templateName == "settings.page.tmpl":
		data = nil

	case templateName == "projects.page.tmpl":
		data = a.templateData.Projects

	case templateName == "project.page.tmpl":
		data = a.templateData.Project

	case templateName == "create.page.tmpl":
		data = map[string]interface{}{
			"formData": a.templateData.FormData,
			"formErrs": a.templateData.FormErrs,
		}

	default:
		panic(fmt.Errorf("app run::unknown template::%s", templateName))
	}

	if err := ts.Execute(w, data); err != nil {
		panic(fmt.Errorf("app run::%s::template err::%s", templateName, err))
	}
}
