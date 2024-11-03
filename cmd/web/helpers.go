package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.Any("method", method), slog.Any("uri", uri), slog.Any("trace", trace))
	app.clientError(w, http.StatusInternalServerError)
}
