package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home handles requests to /
func (a *application) home(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	if r.URL.Path != "/" {
		a.handleError(w, r, fmt.Errorf("resource with URL %v not found", r.URL))
		return
	}

	ts, err := template.ParseFiles([]string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.errLogger.Println(err)
		a.handleError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.handleError(w, r, err)
		return
	}
}

// list handles requests to /list
func (a *application) list(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"msg": "list of todos"}`))
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		a.handleError(w, r, err)
		return
	}

	ts, err := template.ParseFiles([]string{
		"./ui/html/todo.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.handleError(w, r, err)
		return
	}

	if err := ts.Execute(w, id); err != nil {
		a.handleError(w, r, err)
	}
}

// ad handles requests to /list/add
func (a *application) add(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	ts, err := template.ParseFiles([]string{
		"./ui/html/add.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.handleError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.handleError(w, r, err)
		return
	}
}

// settings handles requests to /settings
func (a *application) settings(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println(r.Method, r.URL)

	ts, err := template.ParseFiles([]string{
		"./ui/html/settings.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		a.handleError(w, r, err)
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		a.handleError(w, r, err)
		return
	}
}
