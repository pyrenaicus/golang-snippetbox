package main

import "net/http"

// it returns a http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// pass the servemux as the 'next' parameter to the secureHandlers middleware.
	// It returns a http.Handler
	return secureHeaders(mux)
}
