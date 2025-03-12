package main

import (
	"autobiography/assets"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /technology/{name}", app.technologyHandler)

	return app.recoverPanic(app.securityHeaders(app.setCandidateId(1, mux)))
}
