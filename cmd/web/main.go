package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// define an "addr" command-line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	// read in the flag value & assign it to addr
	flag.Parse()
	// create a logger, additional infor flags are joined with the bitwise OR operator
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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

	// write messages using the new loggers
	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
