package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	// router.HandlerFunc(http.MethodGet, "/ws", wsConnectionHandler)
	router.GET("/ws", app.wsConnectionHandler)

	router.HandlerFunc(http.MethodGet, "/contacts", app.ListContactHandler)
	router.HandlerFunc(http.MethodPost, "/contacts", app.createContactHandler)
	router.HandlerFunc(http.MethodGet, "/contacts/:id", app.showContactHandler)
	router.HandlerFunc(http.MethodPatch, "/contacts/:id", app.updateContactHandler)
	router.HandlerFunc(http.MethodDelete, "/contacts/:id", app.deleteContactHandler)
	router.HandlerFunc(http.MethodGet, "/contacts/:id/history", app.GetContactHistoryHandler)

	return cors.AllowAll().Handler(router)
	// return app.recoverPanic(app.enableCORS(router))
}
