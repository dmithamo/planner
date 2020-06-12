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
func (a *application) renderTemplate(templateName string, w http.ResponseWriter, r *http.Request, data templateData) {
	// retrieve flash msgs, if any
	msg := app.session.PopString(r, "flashMsg")
	data.FlashMsg = msg

	ts, ok := a.templates[templateName]
	if !ok {
		a.serverError(w, r, fmt.Errorf("app run::template not found::`%s`", templateName))
		return
	}

	if err := ts.Execute(w, data); err != nil {
		a.serverError(w, r, fmt.Errorf("app run::template err::`%s`::%s", templateName, err))
		return
	}
}
