package main

import (
	"fmt"
	"html/template"
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

	// use the template.ParseFile() function to read the template file into a template set. If there's an error, we log the detailed error message and use the http.Error() function to send a generic 500 internal server error response to the user
	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// we then use the Execute() method on the template set to write the template content as the response body. the last parameter to Execute() represents the dynamic data that we want to pass in, which for now we'll leave as nil
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
