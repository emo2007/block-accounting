// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package factory

import (
	"github.com/emochka2007/block-accounting/internal/interface/rest"
	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/emochka2007/block-accounting/internal/service"
	"log/slog"
	"os"
)

// Injectors from service.go:

func ProvideService(c config.Config) (service.Service, func(), error) {
	logger := provideLogger(c)
	rootController := provideControllers(logger)
	server := provideRestServer(logger, rootController, c)
	serviceService := service.NewService(logger, server)
	return serviceService, func() {
	}, nil
}

// service.go:

func provideLogger(c config.Config) *slog.Logger {
	lb := new(logger.LoggerBuilder).WithLevel(logger.MapLevel(c.Common.LogLevel)).WithWriter(os.Stdout)

	if c.Common.LogLocal {
		lb.Local()
	}

	if c.Common.LogFile != "" {
		logFile, err := os.Open(c.Common.LogFile)
		if err != nil {
			panic(err)
		}

		lb.WithWriter(logFile)
	}

	if c.Common.LogAddSource {
		lb.WithSource()
	}

	return lb.Build()
}

func provideControllers(
	log *slog.Logger,
) *controllers.RootController {
	return &controllers.RootController{
		Ping: controllers.NewPingController(log.WithGroup("ping-controller")),
		Auth: controllers.NewAuthController(
			log.WithGroup("auth-controller"), presenters.NewAuthPresenter(),
		),
	}
}

func provideRestServer(
	log *slog.Logger, controllers2 *controllers.RootController,
	c config.Config,
) *rest.Server {
	return rest.NewServer(
		log.WithGroup("rest"),
		c.Rest, controllers2,
	)
}
