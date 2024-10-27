package web

import (
	_ "github.com/lafetz/assessment/docs"
	"github.com/lafetz/assessment/internal/web/handlers"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *App) initAppRoutes() {
	a.router.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	a.router.HandleFunc("GET /api/v1/persons", a.recoverPanic(handlers.GetPersons(a.PersonSvc, a.logger)))

	a.router.HandleFunc("GET /api/v1/persons/{personId}", a.recoverPanic(handlers.GetPersonByID(a.PersonSvc, a.logger)))

	a.router.HandleFunc("POST /api/v1/persons", a.recoverPanic(handlers.AddPerson(a.PersonSvc, a.logger, a.validate)))

	a.router.HandleFunc("PUT /api/v1/persons/{personId}", a.recoverPanic(handlers.UpdatePerson(a.PersonSvc, a.logger, a.validate)))

	a.router.HandleFunc("DELETE /api/v1/persons/{personId}", a.recoverPanic(handlers.DeletePerson(a.PersonSvc, a.logger)))
}
