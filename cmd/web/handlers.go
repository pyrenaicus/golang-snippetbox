package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.cnoua.org/internal/models"
)

// change the signature of home handler so it is defined as a method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// get a templateData struct containing the default data (current year)
	// and add the snippets slice to it
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// use render helper
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// retrieve named parameters from request
	params := httprouter.ParamsFromContext(r.Context())

	// get value of "id" parameter
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w) // use notFoud() helper
		return
	}

	// use the snippetModel Get method to retrieve the data for a specific ID, if no
	// matching record found, return a 404 Not Found response
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// call newTempletaData() and use render helper
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

// for now return a placeholder response
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// call r.ParseForm() which adds any data in POST request body to
	// r.PostForm map.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// use r.PostForm.Get() to retrieve the title and content
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// r.PostForm.Get() always return form data as a string, however we expect
	// expires value to be an integer
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// initialize a map to hold any validation errors for the form fields
	fieldErrors := make(map[string]string)

	// check title is not blank and not more than 100 char long, if it fails
	// add a message to the errors map using field name as key
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// chack that content isn't blank
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	// check expires value matches one of the permitted values
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// if there are any errors, dump them in a plain text response and return
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// update redirect path to use new clean URL format
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
