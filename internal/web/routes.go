package web

import "github.com/lafetz/assessment/internal/web/handlers"

func (a *App) initAppRoutes() {
	a.router.HandleFunc("/person", a.recoverPanic(handlers.GetPersons(a.PersonSvc, a.logger)))
	a.router.HandleFunc("/person/{personId}", a.recoverPanic(handlers.GetPersonByID(a.PersonSvc, a.logger)))
	a.router.HandleFunc("/person", a.recoverPanic(handlers.AddPerson(a.PersonSvc, a.logger, a.validate)))
	a.router.HandleFunc("/person/{personId}", a.recoverPanic(handlers.UpdatePerson(a.PersonSvc, a.logger, a.validate)))
	a.router.HandleFunc("/person/{personId}", a.recoverPanic(handlers.DeletePerson(a.PersonSvc, a.logger)))
}
