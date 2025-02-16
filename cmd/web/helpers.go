package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// the serverError helper writes an error message and stack trace to the errorLog, then sends a generic 500 Internal Server Error response to user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding error to user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// for consistency, we'll also implement a notFound helper. this is simply a convenience wrapper around clientError which sends a 404 not found error
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// retrieve the appropriate template set from the cache based on the page name (like 'home.tmpl). tf no entry exists in the cache with the provided name, then create a new error and call the serverError() helper method
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write the template to the buffer, instead of straight to the http.ResponseWriter , if there's an error, call our serverError helper and return
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// if the template is written to buffer without any errors, we are safe to go ahead and write the HTTP status code to http.ResponseWriter
	w.WriteHeader(status)

	//execute the template set and write the response body. again if there's any error we call the serverError
	// err := ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// write the contents of the buffer to http.ResponseWriter, note:this is another time where we pass out http.ResponseWriter to a function that takes an io.Writer.
	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// add the flash message to template data if exists
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		// add the authentication status to template data if exists
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r), // add the CSRFToken to template data
	}
}

// create a new decodePostForm() helper method. the second parameter here, dst, is the target destination that we want to decode the form data into
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// call the ParseForm method on the request, in the same way that we did in our createSnippetPost handler
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// call the Decode method instance, passing the target destination as the first argument
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// if we try to use an invalid destination, the decode() method will return an error with the type *form.InvalidDecoderError. we use the errors.As() to check for this and raise a panic rather then returning an error
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// for other errors we return them as normal
		return err
	}
	return nil
}

// return true if the current request is from an authenticated user, otherwise return false
func (app *application) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
