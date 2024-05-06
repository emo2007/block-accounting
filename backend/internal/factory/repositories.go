package factory

import (
	"fmt"

	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/usecase/repository"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/users"
)

func provideUsersRepository(c config.Config) (users.Repository, func(), error) {
	db, close, err := repository.ProvideDatabaseConnection(c)
	if err != nil {
		return nil, func() {}, fmt.Errorf("error connect to database. %w", err)
	}

	return users.NewRepository(db), close, nil
}
