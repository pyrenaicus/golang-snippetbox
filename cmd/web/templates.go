package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.cnoua.org/internal/models"
)

// define a templateData type to act as the holding structure for
// any dynamic data passed to our html templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// fn returns a formatted string of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// initialize a template.FuncMap object & store it in a global variable. it acts as a
// lookup table for our custom template functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize a new map to act as the cache
	cache := map[string]*template.Template{}
	// get a slice of all filepaths for 'page' templates
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract the file name from the full filepath
		name := filepath.Base(page)

		// use template.New() to create an empty template set, use the Funcs() method
		// to register the template.FuncMap, and then
		// parse the base template into a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// call ParseGlob() on template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// call ParseFiles() to add page template to template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add template set to the map, using page name as the key
		cache[name] = ts
	}

	return cache, nil
}
