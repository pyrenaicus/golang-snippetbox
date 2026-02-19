package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.cnoua.org/internal/models"
)

// define a templateData type to act as the holding structure for
// any dynamic data passed to our html templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

		// create a slice containing the filepaths for our base template, partials and page
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// parse files into a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add template set to the map, using page name as the key
		cache[name] = ts
	}

	return cache, nil
}
