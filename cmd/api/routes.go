package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)

	// Public routes
	mux.Group(func(r chi.Router) {
		// Monitor
		r.Get("/.monitor/health", app.status)

		// Authentication
		r.Post("/register", app.createUser)
		r.Post("/login", app.createAuthenticationToken)
	})

	mux.Route("/api/v1", func(r chi.Router) {

		// Protected routes via JWT tokens
		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthenticatedUser)

			r.Get("/user", app.Me) // Current User
			r.Post("/user/{userID}", app.updateUser)
			r.Delete("/user/{userID}", app.deleteUser)

			r.Get("/protected", app.protected)
		})

		// Protected routes via Basic Auth Protection - BAP
		r.Group(func(r chi.Router) {
			r.Use(app.requireBasicAuthentication)

			r.Get("/bap", app.protected)
		})
	})

	return mux
}
