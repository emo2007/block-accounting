package factory

import (
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/auth"
	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/cache"
	orepo "github.com/emochka2007/block-accounting/internal/infrastructure/repository/organizations"
	txRepo "github.com/emochka2007/block-accounting/internal/infrastructure/repository/transactions"
	urepo "github.com/emochka2007/block-accounting/internal/infrastructure/repository/users"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/chain"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/transactions"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
)

func provideUsersInteractor(
	log *slog.Logger,
	usersRepo urepo.Repository,
	chainInteractor chain.ChainInteractor,
) users.UsersInteractor {
	return users.NewUsersInteractor(log.WithGroup("users-interactor"), usersRepo, chainInteractor)
}

func provideJWTInteractor(
	c config.Config,
	usersInteractor users.UsersInteractor,
	authRepository auth.Repository,
) jwt.JWTInteractor {
	return jwt.NewJWT(c.Common.JWTSecret, usersInteractor, authRepository)
}

func provideOrganizationsInteractor(
	log *slog.Logger,
	orgRepo orepo.Repository,
	cache cache.Cache,
) organizations.OrganizationsInteractor {
	return organizations.NewOrganizationsInteractor(log, orgRepo, cache)
}

func provideTxInteractor(
	log *slog.Logger,
	txRepo txRepo.Repository,
	orgInteractor organizations.OrganizationsInteractor,
	chainInteractor chain.ChainInteractor,
) transactions.TransactionsInteractor {
	return transactions.NewTransactionsInteractor(
		log.WithGroup("transaction-interactor"),
		txRepo,
		orgInteractor,
		chainInteractor,
	)
}

func provideChainInteractor(
	log *slog.Logger,
	config config.Config,
	txRepository txRepo.Repository,
	orgInteractor organizations.OrganizationsInteractor,
) chain.ChainInteractor {
	return chain.NewChainInteractor(log, config, txRepository, orgInteractor)
}
