package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// define a home handler function which writes a byte slide containing
// "Hello from Snippetbox!" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	// check if the current request URL path exactly matches"/". If it doesn't, use
	// the http.NotFount() function to send a 404 response to client
	// Importantly, we then return from the handler. If we don't return the handler would
	// keep executing and also write the "Hello from SnippetBox" message
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox!"))
}

// add a snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	// extract the value of id parameter from the query string and try to
	// convert it to an interger using the strconv.Atoi() function. If it can't be converted to an integer, or the value is less than 1, we return a 404 page not fount response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// use the fmt.Fprintf function to interpolate the id value with our response and write it to the http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {

	// use r.Method to check whether the request is using POST or not
	if r.Method != "POST" {
		// if not, use the w.WriteHeader() function to send a 405 status code
		// and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed

		// use Header().Set() method to add an "Allow:POST" header to the response header map.
		// the first parameter is the name of the header and the second parameter is header value
		//w.Header().Set("Allow", "POST")
		w.Header().Set("Allow", http.MethodPost)

		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))

		// use http.Error() function to send a custom HTTP error
		//response. The first parameter is the response writer, the
		//second parameter is the error message, and the third parameter
		//is the HTTP status code. In this case, we use 405 status code
		//which means "Method Not Allowed"
		//http.Error(w, "Method Not Allowed", 405)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
