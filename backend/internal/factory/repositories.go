package factory

import (
	"database/sql"

	"github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/users"
)

func provideUsersRepository(db *sql.DB) users.Repository {
	return users.NewRepository(db)
}

func provideOrganizationsRepository(db *sql.DB) organizations.Repository {
	return organizations.NewRepository(db)
}
