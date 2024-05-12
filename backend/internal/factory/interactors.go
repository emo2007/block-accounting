package factory

import (
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/cache"
	orepo "github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	urepo "github.com/emochka2007/block-accounting/internal/usecase/repository/users"
)

func provideUsersInteractor(
	log *slog.Logger,
	usersRepo urepo.Repository,
) users.UsersInteractor {
	return users.NewUsersInteractor(log.WithGroup("users-interactor"), usersRepo)
}

func provideJWTInteractor(c config.Config, usersInteractor users.UsersInteractor) jwt.JWTInteractor {
	return jwt.NewWardenJWT(c.Common.JWTSecret, usersInteractor)
}

func provideOrganizationsInteractor(
	log *slog.Logger,
	orgRepo orepo.Repository,
	cache cache.Cache,
) organizations.OrganizationsInteractor {
	return organizations.NewOrganizationsInteractor(log, orgRepo, cache)
}
