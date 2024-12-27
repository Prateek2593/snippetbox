package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Prateek2593/snippetbox/internal/models"
)

// define a templateData type to act as the holding structure for any dynamic data that we want to pass to our HTML template. at the moment it only contains one field, but will add more
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet // include a snippets field in templateData struct
	Form        any
}

// create a humanDate function which returns a nicely formatted string representation of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// initialize a template.FuncMap object and store it in global variable. this is essentially a string-keyed map which acts as a lookup between the names of our custom template functions and functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize a new map to act as a cache
	cache := map[string]*template.Template{}

	// use the filepath.GLob() function to get a slice of all filepaths that match the pattern "./ui/html/pages/*.tmpl". this will essentially gives us a slice of all filepaths for our application 'page' templates, like:[ui/html/pages/home.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// loop through the page filepaths one by one
	for _, page := range pages {
		// extract the name of the file without the extension (e.g., "home") using filepath.Base() function and store it in the name variable. this will be used as the key for our cache map.  e.g., (home.tmpl)
		name := filepath.Base(page)

		// create a slice containing the file path for our base template, any partial and the page
		/*files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// parse the files into a new template set (ts) using the template.ParseFiles() function. this function takes a slice of file paths and returns a new TemplateSet object that can be used to execute the parsed templates. if there's an error during parsing, we return the error immediately, so we don't continue parsing the remaining files.  e.g., ["./ui/html/base.tmpl", "./ui/html/partials/nav.tmpl"]
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// add the parsed template set to our cache map with the key being the name of the template file without the extension (e.g., home.tmpl) and the value being the parsed template set  e.g., home: *template.Template{Name: "home", Templates: []*template.Template{...}}  where... are the parsed templates for the home.tmpl file
		cache[name] = ts
		*/

		// parse the base template file into template set
		// the template.FuncMap must be registered with the template set before you call the ParseFiles(). this means we have to use template.New() to create an empty template set, use the Funcs() method to register the template.FuncMap() and then parse the file as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// call ParseGlob() *on this template set* to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// call the ParseFiles() *on this template set* to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add the template set to mao as normal
		cache[name] = ts
	}

	// return the map
	return cache, nil
}
