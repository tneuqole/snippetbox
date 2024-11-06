package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/tneuqole/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", ping)

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.getSnippet))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.getUserSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.postUser))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.getUserLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.loginUser))

	protected := dynamic.Append(app.requireAuth)
	mux.Handle("POST /user/logout", protected.ThenFunc(app.logoutUser))
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.getSnippetForm))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.postSnippet))

	standard := alice.New(app.recoverPanic, app.logRequest, setCommonResponseHeaders)

	return standard.Then(mux)
}
