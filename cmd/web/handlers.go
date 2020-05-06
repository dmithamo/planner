package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// insertErrMessage is a helper for inserting errors
func insertErrMessage(err error) string {
	return fmt.Sprintf(`{"msg":"something went wrong: %v"}`, err)
}

func home(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println(r.Method, r.URL)

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"msg":"not found"}`))
		return
	}

	ts, err := template.ParseFiles([]string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println(r.Method, r.URL)

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"msg": "list of todos"}`))
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"msg": "invalid id: %v"}`, id)))
		return
	}

	ts, err := template.ParseFiles([]string{
		"./ui/html/todo.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}

	if err := ts.Execute(w, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println(r.Method, r.URL)

	ts, err := template.ParseFiles([]string{
		"./ui/html/add.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}
}

func settings(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println(r.Method, r.URL)

	ts, err := template.ParseFiles([]string{
		"./ui/html/settings.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}...)

	if err != nil {
		errLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}

	if err := ts.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(insertErrMessage(err)))
		return
	}
}
