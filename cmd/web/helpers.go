package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// serverError writes an error message & stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
}

// clientError sends a specific status code & description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// convenience wrapper around clientError which sends a 404 to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// helper which returns a pointer to a templateData struct initialized without
// the current year
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// retrieve template set from cache based on the page name, if no entry exists
	// create a new error & call serverError() helper method
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write template to buffer, if there's an error call serverError()
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// if template is written in buffer without errors, write HTTP status code
	w.WriteHeader(status)

	// write contents of buffer to http.ResponseWriter
	buf.WriteTo(w)
}

// decodePostForm() helper method, the second parameter, dst, is the target
// destination where we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// if we try to use an invalid target destination, Decode() will return
		// an error with type *form.InvalidDecoderError. Check with errors.As()
		// and raise a panic if found, instead of returnin the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// for all other errors, return as normal
		return err
	}

	return nil
}
