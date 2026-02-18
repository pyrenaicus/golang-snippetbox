package main

import "snippetbox.cnoua.org/internal/models"

// define a templateData type to act as the holding structure for
// any dynamic data passed to our html templates.
type templateData struct {
	Snippet *models.Snippet
}
