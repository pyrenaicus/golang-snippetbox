package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// define an application struct to hold app-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// define an "addr" command-line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	// read in the flag value & assign it to addr
	flag.Parse()
	// create a logger, additional infor flags are joined with the bitwise OR operator
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// init a new instance of application struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

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
	// initialize a new http.Server struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	// write messages using the new loggers
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
