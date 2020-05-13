package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// landingPage handles requests to /
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	if r.URL.Path != "/" {
		a.serveError(w, r, fmt.Errorf("resource with URL %v not found", r.URL))
		return
	}

	ts, err := template.ParseFiles([]string{
		"./views/html/auth.page.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.errLogger.Println(err)
		a.serveError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.serveError(w, r, err)
		return
	}
	a.infoLogger.Println(http.StatusOK)
}

// ad handles requests to /list/listOfTodos
func (a *application) listOfTodos(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	ts, err := template.ParseFiles([]string{
		"./views/html/home.page.tmpl",
		"./views/html/base.layout.tmpl",
		"./views/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.serveError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.serveError(w, r, err)
		return
	}

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
		a.serveError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.serveError(w, r, err)
		return
	}
	a.infoLogger.Println(http.StatusOK)
}
