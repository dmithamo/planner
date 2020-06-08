package main

import (
	"fmt"
	"net/http"
)

// landingPage handles requests to / <that is to say, root>
func (a *application) landingPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFoundErr(w, r, fmt.Errorf("resource with URL `%v` not found", r.URL))
		return
	}

	a.renderTemplate("auth.page.tmpl", nil, w, r)
}
