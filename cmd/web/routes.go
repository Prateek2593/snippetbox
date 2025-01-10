package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// the routes() method returns a servemux containing our application routes
// update the signature for the routes() method so that it returns a http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux, then register the home handler function at the root ("/") path.
	// mux := http.NewServeMux()

	// initialize the router
	router := httprouter.New()

	// create a handler function which wraps our notFound() helper, and then assign it as the custom handler for 404 not found responses. you can also set a custom handler for 405 Method Not Allowed responses by setting router.MethodNotAllowed in same way
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// create a file server which serves files out of the "./ui/static" directory. note that the path given to the http.Dir function is relative to the project root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use the mux.Handle() function to register the file server as the handler for all url paths that start with "/static/". for matching paths, we strip the "/static" prefix before the request reaches the file server
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// update the pattern for the route for static files
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)

	// create a new middleware chaing containing the middleware specific to our dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// and then create routes using the appropriate methods, patterns and handlers
	// update these routes to use the dynamic middleware chain followed by the appropriate handler function. note that because the alice ThenFunc() method returns a http.Handler(rather than a http.HandlerFunc) we also need to switch to registering the route using router.Handler method
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// create a middleware chain containing our standard middlewares which will be used for every request our application receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// pass the servermux as the 'next' parameter to the secureHeaders middleware. because secureHeaders is just a function and the function returns a http.Handler we dont need to do anything else
	// wrap the existing chain with logRequest middleware
	// wrap the existing chain with recoverPanic middleware
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// return the standard middleware chain followed by the servemux
	return standard.Then(router)
}
