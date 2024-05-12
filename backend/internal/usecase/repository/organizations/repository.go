package organizations

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
	Ids    uuid.UUIDs
	UserId uuid.UUID

	OffsetDate time.Time
	CursorId   uuid.UUID
	Limit      int64
}

type AddParticipantParams struct {
	OrganizationId uuid.UUID
	UserId         uuid.UUID
	EmployeeId     uuid.UUID
	IsAdmin        bool
}

type Repository interface {
	Create(ctx context.Context, org models.Organization) error
	Get(ctx context.Context, params GetParams) ([]*models.Organization, error)
	Update(ctx context.Context, org models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddParticipant(ctx context.Context, params AddParticipantParams) error
	CreateAndAdd(ctx context.Context, org models.Organization, user *models.User) error
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
			"id",
			"name",
			"address",
			"wallet_seed",
			"created_at",
			"updated_at",
		).Values(
			org.ID,
			org.Name,
			org.Address,
			org.WalletSeed,
			org.CreatedAt,
			org.UpdatedAt,
		).PlaceholderFormat(sq.Dollar)

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
	organizations := make([]*models.Organization, 0, params.Limit)

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Select(
			"o.id",
			"o.name",
			"o.address",
			"o.wallet_seed",
			"o.created_at",
			"o.updated_at",
		).From("organizations as o").
			Limit(uint64(params.Limit)).
			PlaceholderFormat(sq.Dollar)

		if params.UserId != uuid.Nil {
			query = query.InnerJoin("organizations_users as ou on o.id = ou.organization_id").
				Where(sq.Eq{
					"ou.user_id": params.UserId,
				})
		}

		if params.CursorId != uuid.Nil {
			query = query.Where(sq.Gt{
				"o.id": params.CursorId,
			})
		}

		if params.Ids != nil {
			query = query.Where(sq.Eq{
				"o.id": params.Ids,
			})
		}

		if !params.OffsetDate.IsZero() {
			query = query.Where(sq.GtOrEq{
				"o.updated_at": params.OffsetDate,
			})
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch organizations from database. %w", err)
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				err = errors.Join(fmt.Errorf("error close rows. %w", closeErr), err)
			}
		}()

		for rows.Next() {
			var (
				id         uuid.UUID
				name       string
				address    string
				walletSeed []byte
				createdAt  time.Time
				updatedAt  time.Time
			)

			if err = rows.Scan(
				&id,
				&name,
				&address,
				&walletSeed,
				&createdAt,
				&updatedAt,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			organizations = append(organizations, &models.Organization{
				ID:         id,
				Name:       name,
				Address:    address,
				WalletSeed: walletSeed,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	return organizations, nil
}

func (r *repositorySQL) Update(ctx context.Context, org models.Organization) error {
	panic("implement me!")

	return nil
}

func (r *repositorySQL) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me!")

	return nil
}

func (r *repositorySQL) AddParticipant(ctx context.Context, params AddParticipantParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Insert("organizations_users").Columns(
			"organization_id",
			"user_id",
			"employee_id",
			"added_at",
			"updated_at",
			"is_admin",
		).Values(
			params.OrganizationId,
			params.UserId,
			params.EmployeeId,
			time.Now(),
			time.Now(),
			params.IsAdmin,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add new participant to organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) CreateAndAdd(ctx context.Context, org models.Organization, user *models.User) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		if err := r.Create(ctx, org); err != nil {
			return fmt.Errorf("error create organization. %w", err)
		}

		if err := r.AddParticipant(ctx, AddParticipantParams{
			OrganizationId: org.ID,
			UserId:         user.Id(),
			IsAdmin:        true,
		}); err != nil {
			return fmt.Errorf("error add user to newly created organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}
