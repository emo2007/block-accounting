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
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/chain"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/transactions"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/auth"
)

var interfaceSet wire.ProviderSet = wire.NewSet(
	provideAuthController,
	provideOrganizationsController,
	provideControllers,
	provideTxController,
	provideParticipantsController,

	provideAuthPresenter,
	provideOrganizationsPresenter,
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

func provideOrganizationsPresenter() presenters.OrganizationsPresenter {
	return presenters.NewOrganizationsPresenter()
}

func provideAuthController(
	log *slog.Logger,
	usersInteractor users.UsersInteractor,
	authPresenter presenters.AuthPresenter,
	jwtInteractor jwt.JWTInteractor,
	repo auth.Repository,
) controllers.AuthController {
	return controllers.NewAuthController(
		log.WithGroup("auth-controller"),
		authPresenter,
		usersInteractor,
		jwtInteractor,
		repo,
	)
}

func provideOrganizationsController(
	log *slog.Logger,
	organizationsInteractor organizations.OrganizationsInteractor,
	presenter presenters.OrganizationsPresenter,
) controllers.OrganizationsController {
	return controllers.NewOrganizationsController(
		log.WithGroup("organizations-controller"),
		organizationsInteractor,
		presenter,
	)
}

func provideTxController(
	log *slog.Logger,
	txInteractor transactions.TransactionsInteractor,
	chainInteractor chain.ChainInteractor,
	organizationsInteractor organizations.OrganizationsInteractor,
) controllers.TransactionsController {
	return controllers.NewTransactionsController(
		log.WithGroup("transactions-controller"),
		txInteractor,
		presenters.NewTransactionsPresenter(),
		chainInteractor,
		organizationsInteractor,
	)
}

func provideParticipantsController(
	log *slog.Logger,
	orgInteractor organizations.OrganizationsInteractor,
	usersInteractor users.UsersInteractor,
) controllers.ParticipantsController {
	return controllers.NewParticipantsController(
		log.WithGroup("participants-controller"),
		orgInteractor,
		usersInteractor,
		presenters.NewParticipantsPresenter(),
	)
}

func provideControllers(
	log *slog.Logger,
	authController controllers.AuthController,
	orgController controllers.OrganizationsController,
	txController controllers.TransactionsController,
	participantsController controllers.ParticipantsController,
) *controllers.RootController {
	return controllers.NewRootController(
		controllers.NewPingController(log.WithGroup("ping-controller")),
		authController,
		orgController,
		txController,
		participantsController,
	)
}

func provideRestServer(
	log *slog.Logger,
	controllers *controllers.RootController,
	c config.Config,
	jwt jwt.JWTInteractor,
) *rest.Server {
	return rest.NewServer(
		log.WithGroup("rest"),
		c.Rest,
		controllers,
		jwt,
	)
}
