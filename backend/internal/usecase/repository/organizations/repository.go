package organizations

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
	Ids uuid.UUIDs
}

type Repository interface {
	Create(ctx context.Context, org models.Organization) error
	Get(ctx context.Context, params GetParams) ([]*models.Organization, error)
	Update(ctx context.Context, org models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
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

func (r *repositorySQL) Create(ctx context.Context, org models.Organization) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Insert("organizations").Columns(
			"id, name",
		).Values(
			org.ID,
			org.Name,
		)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error insert new organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Get(ctx context.Context, params GetParams) ([]*models.Organization, error) {
	panic("implement me!")

	return nil, nil
}

func (r *repositorySQL) Update(ctx context.Context, org models.Organization) error {
	panic("implement me!")

	return nil
}

func (r *repositorySQL) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me!")

	return nil
}
