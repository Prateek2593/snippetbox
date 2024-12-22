package main

import (
	"log"
	"net/http"
)

// define a home handler function which writes a byte slide containing
// "Hello from Snippetbox!" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox!"))
}

// add a snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a specific snippet..."))
}

// Add a snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}
func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home handler function at the root ("/") path.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: TCP network address to listen on(in this case":4000)
	// and the servemux we just created. If http.ListenAndServe() return an error,
	// we use the log.Fatal() function to log the error message and exit.
	// note that any error returned by http.ListenAndServe() is always non-nil
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
