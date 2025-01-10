package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Prateek2593/snippetbox/internal/models"
	"github.com/Prateek2593/snippetbox/internal/validator"
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

	// use the PopString() method to retrieve the value for the "flash" key. PopString also deletes the key and the value from session data, so it acts like a one time fetch. if there is no matching key in session data this will return the empty string
	// flash := app.sessionManager.PopString(r.Context(), "flash")

	// call the newTemplateData() helper to get a templateData struct containing the 'default' data(which for now is just the current year) and add the snippet slice to it
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// pass the flash data to the template
	// data.Flash = flash

	// pass the data to render() as normal
	app.render(w, http.StatusOK, "view.tmpl", data)
}

// add a new snippetCreate handler
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// initialize a new createSnippetForm instance and pass it to the template. notice how this is also a great opportunity to set any default or initial values for the form ---  here we set the initial value for the snippet expiry to 365 days
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// remove the explicit FieldErrors struct field and instead embed the Validator type, embedding this means that out snippetCreateForm "inherits" all the fields and methods of our Validator type
// update our snippetCreateForm struct to include struct tags which tell the decoder how to map HTML form values into different struct fields.
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// FieldErrors map[string]string
	validator.Validator `form:"-"` // "-" tells decoder to completely ignore a field during decoding
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
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// use the r.PostForm.Get() method to retrieve the title and content from the r.PostForm map
	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")

	// the r.PostForm.Get() method always returns the form data as a string. however we are expecting our expires value to be a number, and want to represent it in out go code as an integer. so we need to manually convert the form data to an integer using strcov.Atoi, and we send a 400 bad request if conversion fails
	/*
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := snippetCreateForm{
			Title:   r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			Expires: expires,
			// FieldErrors: map[string]string{},
		}
	*/
	// initialize a map to hold any validation errors for the form fields
	// fieldErrors := make(map[string]string)

	/*
		// chech that the title value is not blank and is not more than 100 characters long, if it fails either of those checks add a message to errors map using the field name as the key
		if strings.TrimSpace(form.Title) == "" {
			form.FieldErrors["title"] = "This field cannot be empty"
		} else if utf8.RuneCountInString(form.Title) > 100 {
			form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
		}

		// check that the content value isnt blank
		if strings.TrimSpace(form.Content) == "" {
			form.FieldErrors["content"] = "This field cannot be empty"
		}

		// check the expires value matches one of the permitted values(1,7,365)
		if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
			form.FieldErrors["expires"] = "This field must be equal 1,7, or 365"
		}

		// if there are any errors, dump them in a plain text HTTP response and return from handler
		if len(form.FieldErrors) > 0 {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}
	*/

	// declare a new empty instance of snippetCreateForm struct
	var form snippetCreateForm

	// call the Decode() method on the form decoder, passing in the current request and a pointer to our snippetCreateForm struct. this will essentially fill our struct with the relevant data from the HTML form. if there's a problem, we return a 400 Bad Request error
	// err = app.formDecoder.Decode(&form, r.PostForm)
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// because the Validator type is embedded by the snippetCreateForm struct, we can call checkField() directly on it to execute our validation checks. checkField() will add the provided key and error message to the FieldErrors map if the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal to 1, 7, or 365")

	// use the valid method to see if any of the checks failed. if they did, then re render the template passing in the form in same way as before
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use the Put() method to add a string value("Snippet created successfully") and the corresponding key ("flash") to session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for signup")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for login...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
