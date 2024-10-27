package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/lafetz/assessment/internal/config"
	person "github.com/lafetz/assessment/internal/core/service"
	customlogger "github.com/lafetz/assessment/internal/logger"
	"github.com/lafetz/assessment/internal/repository"

	"github.com/lafetz/assessment/internal/web"
	customvalidator "github.com/lafetz/assessment/internal/web/validation"
)

func main() {
	config := config.NewConfig()
	logger := customlogger.NewLogger(config.LogLevel, config.Env)
	repo := repository.NewRepository()
	repo.SeedData()
	personSvc := person.NewPersonSvc(repo)
	val := validator.New()
	custonmVal := customvalidator.NewCustomValidator(val)
	web := web.NewApp(config.Port, logger, personSvc, custonmVal)
	logger.Info("running web server")
	err := web.Run()
	if err != nil {
		logger.Error("web server error", "error", err)
	}
}
