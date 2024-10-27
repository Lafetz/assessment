package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	person "github.com/lafetz/assessment/internal/core/service"
	customvalidator "github.com/lafetz/assessment/internal/web/validation"
)

type App struct {
	port      int
	router    *http.ServeMux
	logger    *slog.Logger
	PersonSvc person.PersonSvcApi
	validate  *customvalidator.CustomValidator
}

func NewApp(port int, logger *slog.Logger, personSvc person.PersonSvcApi, validate *customvalidator.CustomValidator) *App {
	a := &App{
		router:    http.NewServeMux(),
		logger:    logger,
		port:      port,
		PersonSvc: personSvc,
		validate:  validate,
	}
	a.initAppRoutes()
	return a
}
func (a *App) Run() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(a.port)),
		Handler:      a.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		a.logger.Info("shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	a.logger.Info("server stopped")
	return nil
}
