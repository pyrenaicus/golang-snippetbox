package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// define an "addr" command-line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	// read in the flag value & assign it to addr
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register the file server as the handler for all URL paths
	// starting with "/static/". We strip the "/static" before
	// the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// register other app routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// the value returned from flag.String() is a pointer to the flag value,
	// so we need to dereference the pointer before using it.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
