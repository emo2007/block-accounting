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
		query := sq.Select("u.id, u.name, u.email, u.phone, u.tg, u.seed, u.created_at, u.activated_at, u.public_key, u.mnemonic").
			From("users as u").
			PlaceholderFormat(sq.Dollar)

		if len(params.Ids) > 0 {
			query = query.Where(sq.Eq{
				"u.id": params.Ids,
			})
		}

		if params.OrganizationId != uuid.Nil {
			query = query.InnerJoin(
				"organizations_users as ou on ou.user_id = u.id",
			).Where(sq.Eq{
				"ou.organization_id": params.OrganizationId,
			})
		}

		if params.Seed != nil {
			query = query.Where("u.seed = ?", params.Seed)
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
				id uuid.UUID

				name  string
				email string
				phone string
				tg    string

				seed []byte
				pk   []byte
				//isAdmin     bool
				createdAt   time.Time
				activatedAt sql.NullTime
				mnemonic    string
			)

			if err = rows.Scan(
				&id,
				&name,
				&email,
				&phone,
				&tg,
				&seed,
				&createdAt,
				&activatedAt,
				&pk,
				&mnemonic,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			users = append(users, &models.User{
				ID:   id,
				Name: name,
				Credentails: &models.UserCredentials{
					Email:    email,
					Phone:    phone,
					Telegram: tg,
				},
				Bip39Seed: seed,
				PK:        pk,
				Mnemonic:  mnemonic,
				//Admin:     isAdmin,
				CreatedAt: createdAt,
				Activated: activatedAt.Valid,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repositorySQL) Create(ctx context.Context, user *models.User) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		columns := []string{"id", "name", "email", "phone", "tg", "seed", "public_key", "mnemonic", "created_at"}

		values := []any{
			user.ID,
			user.Name,
			user.Credentails.Email,
			user.Credentails.Phone,
			user.Credentails.Telegram,
			user.Bip39Seed,
			user.PK,
			user.Mnemonic,
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
		return err
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
		return err
	}

	return nil
}

func (r *repositorySQL) Update(ctx context.Context, user *models.User) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repositorySQL) Delete(ctx context.Context, id string) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {

		return nil
	}); err != nil {
		return err
	}

	return nil
}
