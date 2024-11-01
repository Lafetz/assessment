package web

import (
	_ "github.com/lafetz/assessment/docs"
	"github.com/lafetz/assessment/internal/web/handlers"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *App) initAppRoutes() {
	a.Router.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	a.Router.HandleFunc("GET /api/v1/persons", a.recoverPanic(a.enableCORS(handlers.GetPersons(a.PersonSvc, a.logger))))
	a.Router.HandleFunc("GET /api/v1/persons/{personId}", a.recoverPanic(a.enableCORS(handlers.GetPersonByID(a.PersonSvc, a.logger))))
	a.Router.HandleFunc("POST /api/v1/persons", a.recoverPanic(a.enableCORS(handlers.AddPerson(a.PersonSvc, a.logger, a.validate))))
	a.Router.HandleFunc("PUT /api/v1/persons/{personId}", a.recoverPanic(a.enableCORS(handlers.UpdatePerson(a.PersonSvc, a.logger, a.validate))))
	a.Router.HandleFunc("DELETE /api/v1/persons/{personId}", a.recoverPanic(a.enableCORS(handlers.DeletePerson(a.PersonSvc, a.logger))))
	a.Router.HandleFunc("/", a.recoverPanic(a.enableCORS(handlers.NotFound())))
}
