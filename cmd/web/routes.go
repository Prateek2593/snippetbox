package main

import "net/http"

// the routes() method returns a servemux containing our application routes
// update the signature for the routes() method so that it returns a http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux, then register the home handler function at the root ("/") path.
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory. note that the path given to the http.Dir function is relative to the project root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use the mux.Handle() function to register the file server as the handler for all url paths that start with "/static/". for matching paths, we strip the "/static" prefix before the request reaches the file server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// pass the servermux as the 'next' parameter to the secureHeaders middleware. because secureHeaders is just a function and the function returns a http.Handler we dont need to do anything else
	// wrap the existing chain with logRequest middleware
	// wrap the existing chain with recoverPanic middleware
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
