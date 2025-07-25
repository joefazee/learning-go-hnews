package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	defaultMiddleware := alice.New(app.recover, app.loggger)
	secureMiddleware := alice.New(app.session.Enable, app.authenticate)

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(app.publicPath))))
	mux.Handle("/", secureMiddleware.ThenFunc(app.home))
	mux.Handle("/login", secureMiddleware.ThenFunc(app.login))
	mux.Handle("/logout", secureMiddleware.ThenFunc(app.logout))
	mux.Handle("/submit", secureMiddleware.Append(app.requireAuth).ThenFunc(app.submit))
	mux.Handle("/vote", secureMiddleware.Append(app.requireAuth).ThenFunc(app.vote))
	mux.Handle("/comments", secureMiddleware.Append(app.requireAuth).ThenFunc(app.comments))
	mux.Handle("/register", secureMiddleware.ThenFunc(app.register))
	mux.Handle("/contact", secureMiddleware.ThenFunc(app.contact))
	mux.Handle("/about", secureMiddleware.ThenFunc(app.about))

	return defaultMiddleware.Then(mux)
}
