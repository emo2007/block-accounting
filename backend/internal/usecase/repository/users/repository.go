package users

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	sqltools "github.com/emochka2007/block-accounting/internal/pkg/sqlutils"
	"github.com/google/uuid"
)

type GetParams struct {
	Ids            uuid.UUIDs
	OrganizationId uuid.UUIDs
	Seed           []byte
}

// todo implement
type Repository interface {
	Get(ctx context.Context, params GetParams) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Activate(ctx context.Context, id string) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type repositorySQL struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositorySQL{
		db: db,
	}
}

func (s *repositorySQL) Conn(ctx context.Context) sqltools.DBTX {
	if tx, ok := ctx.Value(sqltools.TxCtxKey).(*sql.Tx); ok {
		return tx
	}

	return s.db
}

func (r *repositorySQL) Get(ctx context.Context, params GetParams) (*models.User, error) {
	var user *models.User

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	return user, nil
}

func (r *repositorySQL) Create(ctx context.Context, user *models.User) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		columns := []string{"id", "seed", "created_at"}

		values := []any{
			user.ID,
			user.Bip32Seed,
			user.CreatedAt,
		}

		if user.Activated {
			columns = append(columns, "activated_at")
			values = append(values, user.CreatedAt)
		}

		query := sq.Insert("users").Columns(
			columns...,
		).Values(
			values...,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error insert new user. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Activate(ctx context.Context, id string) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Update(ctx context.Context, user *models.User) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Delete(ctx context.Context, id string) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}
