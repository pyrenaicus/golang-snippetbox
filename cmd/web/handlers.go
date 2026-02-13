package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// change the signature of home handler so it is defined as a method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w, r) // use notFound() helper
		return
	}

	// init a slice containing path to files
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// read template file into a template set
	// passing the slice as a variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // user serverError() helper
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// write the template content as the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err) // use serverError() helper
		http.Error(w, "Internal server error", 500)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notfound(w) // use notFoud() helper
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed) // use clientError() helper
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
