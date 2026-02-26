package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// it returns a http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
	// initialize the router
	router := httprouter.New()

	// custom handler for 404 responses
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// route for the static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// middleware chain containing the middleware specific to dynamic
	// application routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// update routes to use the new dynamic middleware chain followed by the
	// appropriate handler fn. Note that because the alice ThenFunc() returns
	// a http.Handler (rather than a http.HandlerFunc) we also need to switch
	// to registering the route using the router.Handler() method.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// create a middleware chain used for every request
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
