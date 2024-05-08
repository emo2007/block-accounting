package factory

import (
	"log/slog"
	"os"

	"github.com/google/wire"

	"github.com/emochka2007/block-accounting/internal/interface/rest"
	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
)

var interfaceSet wire.ProviderSet = wire.NewSet(
	provideAuthController,
	provideControllers,

	provideAuthPresenter,
)

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

func provideAuthPresenter(
	jwtInteractor jwt.JWTInteractor,
) presenters.AuthPresenter {
	return presenters.NewAuthPresenter(jwtInteractor)
}

func provideAuthController(
	log *slog.Logger,
	usersInteractor users.UsersInteractor,
	authPresenter presenters.AuthPresenter,
) controllers.AuthController {
	return controllers.NewAuthController(
		log.WithGroup("auth-controller"),
		authPresenter,
		usersInteractor,
	)
}

func provideControllers(
	log *slog.Logger,
	authController controllers.AuthController,
) *controllers.RootController {
	return &controllers.RootController{
		Ping: controllers.NewPingController(log.WithGroup("ping-controller")),
		Auth: authController,
	}
}

func provideRestServer(
	log *slog.Logger,
	controllers *controllers.RootController,
	c config.Config,
) *rest.Server {
	return rest.NewServer(
		log.WithGroup("rest"),
		c.Rest,
		controllers,
	)
}
