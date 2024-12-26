package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Prateek2593/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

// define a home handler function which writes a byte slide containing
// "Hello from Snippetbox!" as the response body
// change the signature of the home handler so it is defined as a method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check if the current request URL path exactly matches"/". If it doesn't, use
	// the http.NotFount() function to send a 404 response to client
	// Importantly, we then return from the handler. If we don't return the handler would
	// keep executing and also write the "Hello from SnippetBox" message

	// because httprouter matches the "/" path exactly, we can now remove the manual check of r.URL.Path != "/" from this handler
	// if r.URL.Path != "/" {
	// 	app.notFound(w) // use the notFound() helper
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// initialize a slice containing the paths to the tow files. its important to note that the file containing our base template must be the *first* file in the slice
	/*files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl", // include navigation partial in template file
		"./ui/html/pages/home.tmpl",
	}

	// use the template.ParseFile() function to read the template file into a template set. If there's an error, we log the detailed error message and use the http.Error() function to send a generic 500 internal server error response to the user
	// notice that we can pass the slice of file paths as a variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// because the home handler function is now a method against application it can access its fields, including the error logger
		// app.errorLog.Println(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	// we then use the Execute() method on the template set to write the template content as the response body. the last parameter to Execute() represents the dynamic data that we want to pass in, which for now we'll leave as nil
	//err = ts.Execute(w, nil)

	//use the ExecuteTemplate method to write tehe content of base template as response body
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		// app.errorLog.Println(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		app.serverError(w, err)
	}
	*/

	// call the newTemplateData() helper to get a templateData struct containing the 'default' data(which for now is just the current year) and add the snippet slice to it
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// pass the data to render() as normal
	app.render(w, http.StatusOK, "home.tmpl", data)
	// w.Write([]byte("Hello from Snippetbox!"))
}

// add a snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// when httprouter is parsing a request, the values of any named parameters will be stored in the request contect, ParamsFromContext() function is used to retrieve a slice containing these parameter  names and values
	params := httprouter.ParamsFromContext(r.Context())

	// extract the value of id parameter from the query string and try to
	// convert it to an interger using the strconv.Atoi() function. If it can't be converted to an integer, or the value is less than 1, we return a 404 page not fount response.
	// we can then use the ByName() method to get the value of the "id" named parameter from the slice and validate it as normal
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// use the SnippetModel objects Get method to retrieve the data for a specific record based on its ID. if no matching recored is found, return a 404 not found response
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// initialize a slice containing the paths to the view.tmpl file, plus the base layout and navigation partial that we made earlier
	/*files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl", // include navigation partial in template file
		"./ui/html/pages/view.tmpl",
	}

	// parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// create a instance of templateData struct for holding snippet data
	data := &templateData{
		Snippet: snippet,
	}

	// and then execute them, notice how we are passing in the snippet data (a models.Snippet struct) as final parameter
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	// use the fmt.Fprintf function to interpolate the id value with our response and write it to the http.ResponseWriter
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	// write the snippet data as a plain text HTTP response body
	// pass in templateData struct when executing the template
	fmt.Fprintf(w, "%+v", data)
	*/

	// call the newTemplateData() helper to get a templateData struct containing the 'default' data(which for now is just the current year) and add the snippet slice to it
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// pass the data to render() as normal
	app.render(w, http.StatusOK, "view.tmpl", data)
}

// add a new snippetCreate handler
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Add a snippetCreate handler function
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// checking if the request method is POST is now superfluous and can be removed, because this is done automatically by httprouter
	/*
		// use r.Method to check whether the request is using POST or not
		if r.Method != "POST" {
			// if not, use the w.WriteHeader() function to send a 405 status code and the w.Write() method to write a "Method Not Allowed" response body. We then return from the function so that the subsequent code is not executed

			// use Header().Set() method to add an "Allow:POST" header to the response header map. the first parameter is the name of the header and the second parameter is header value
			// w.Header().Set("Allow", "POST")
			w.Header().Set("Allow", http.MethodPost)

			//w.WriteHeader(405)
			//w.Write([]byte("Method Not Allowed"))

			// use http.Error() function to send a custom HTTP error response. The first parameter is the response writer, the second parameter is the error message, and the third parameter is the HTTP status code. In this case, we use 405 status code which means "Method Not Allowed"

			// http.Error(w, "Method Not Allowed", 405)
			// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}
	*/

	// first we call the r.ParseForm() which adds any data in POST request bodies to the r.PostForm map. this also workds in same way for PUT and PATCH requests. if there are any errors, we use our clientError helper to send a 400 bad request response to user
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// use the r.PostForm.Get() method to retrieve the title and content from the r.PostForm map
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// the r.PostForm.Get() method always returns the form data as a string. however we are expecting our expires value to be a number, and want to represent it in out go code as an integer. so we need to manually convert the form data to an integer using strcov.Atoi, and we send a 400 bad request if conversion fails
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
