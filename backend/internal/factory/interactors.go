package factory

import (
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	urepo "github.com/emochka2007/block-accounting/internal/usecase/repository/users"
)

func provideUsersInteractor(
	log *slog.Logger,
	usersRepo urepo.Repository,
) users.UsersInteractor {
	return users.NewUsersInteractor(log.WithGroup("users-interactor"), usersRepo)
}

func provideJWTInteractor(c config.Config) jwt.JWTInteractor {
	return jwt.NewWardenJWT(c.Common.JWTSecret)
}
