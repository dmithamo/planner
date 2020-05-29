package main

import (
	"html/template"
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
