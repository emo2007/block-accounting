package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	sqltools "github.com/emochka2007/block-accounting/internal/pkg/sqlutils"
	"github.com/google/uuid"
)

type GetParams struct {
	Ids            uuid.UUIDs
	OrganizationId uuid.UUID
	Seed           []byte
}

// todo implement
type Repository interface {
	Get(ctx context.Context, params GetParams) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Activate(ctx context.Context, id uuid.UUID) error
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

func (r *repositorySQL) Get(ctx context.Context, params GetParams) ([]*models.User, error) {
	var users []*models.User = make([]*models.User, 0, len(params.Ids))

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Select("id, seed, created_at, activated_at").
			From("users").
			PlaceholderFormat(sq.Dollar)

		if len(params.Ids) > 0 {
			query = query.Where("id", params.Ids)
		}

		// if params.OrganizationId != uuid.Nil {
		// 	// todo join org users
		// }

		if params.Seed != nil {
			query = query.Where("seed = ?", params.Seed)
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch data from database. %w", err)
		}

		defer func() {
			if cErr := rows.Close(); cErr != nil {
				err = errors.Join(fmt.Errorf("error close database rows. %w", cErr), err)
			}
		}()

		for rows.Next() {
			var (
				id   uuid.UUID
				seed []byte
				//isAdmin     bool
				createdAt   time.Time
				activatedAt sql.NullTime
			)

			if err = rows.Scan(&id, &seed, &createdAt, &activatedAt); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			users = append(users, &models.User{
				ID:        id,
				Bip32Seed: seed,
				//Admin:     isAdmin,
				CreatedAt: createdAt,
				Activated: activatedAt.Valid,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	return users, nil
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

func (r *repositorySQL) Activate(ctx context.Context, id uuid.UUID) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Update("users").
			SetMap(sq.Eq{
				"activated_at": time.Now(),
			}).
			Where(sq.Eq{
				"id": id,
			})

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error mark user as activated in database. %w", err)
		}

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
