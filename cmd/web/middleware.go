package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a deferred function (which will always be run in the event of a panic as GO unwinds the stack)
		defer func() {
			// use the builtin recover function to check if there has been a panic or not, if there has...
			if err := recover(); err != nil {
				// set a "Connection:close" header on the response
				w.Header().Set("Connection", "close")
				// call the app.serverError helper function to return a 500 internal server error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the user is not authenticated, redirect them to login page and return from the middleware chain so that no subsequent handlers in the chain are executed
		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// otherwise set the "Cache-Control:no-store" header so that pages require authentication are not stored in the users browser cache
		w.Header().Add("Cache-Control", "no-store")

		// and call the next handler in chain
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf Middleware function which uses a customized CSRF cookie with secure, path and HttpOnly attribures set
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve the authenticatedUserID value from the session using the GetInt() method. this will return the zero value for an int(0) if no "authenticatedUserID" value is in the session -- in which case we call the next handler in the chain as normal and return
		id := app.sessionManager.GetInt(r.Context(), "authenticationUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// otherwise, we check to see if the user with that ID exists in our database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// if a matching user is found, we know that the request is coming from an authenticated user who exists in our database, we create a new copy of the request(with an isAuthenticatedContextKey value of true in the request context) and assign it to r
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
