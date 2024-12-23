package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	// define a new command line flag with name addr, a default value of ":4000" and some short help text explaining what the flag controls. the value of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// importantly, we use the flag.Parse() function to parse the command line flag. this reads in command line flag value and assigns it to the addr variable. you need to call this *before* you use the addr variable otherwise it will always contain the default value of ":4000". if any errors are encountered during parsing the application will be terminated
	flag.Parse()

	// Use the http.NewServeMux() function to initialize a new servemux, then register the home handler function at the root ("/") path.
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory. note that the path given to the http.Dir function is relative to the project root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use the mux.Handle() function to register the file server as the handler for all url paths that start with "/static/". for matching paths, we strip the "/static" prefix before the request reaches the file server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	//the value returned from the flag.String() function is a pointer to the flag value, not the value itself. so we need to dereference the pointer(i.e. prefix it with the *symbol) before using it. note that we're using the log.Printf() function to interpolate the address with log message.
	log.Printf("Starting server on %s", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in two parameters: TCP network address to listen on(in this case":4000) and the servemux we just created. If http.ListenAndServe() return an error, we use the log.Fatal() function to log the error message and exit. note that any error returned by http.ListenAndServe() is always non-nil
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
