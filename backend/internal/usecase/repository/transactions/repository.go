package transactions

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

type GetTransactionsParams struct {
	Ids            uuid.UUIDs
	OrganizationId uuid.UUID
	CreatedById    uuid.UUID
	To             []byte
	Limit          int64
	CursorId       uuid.UUID

	WithCancelled bool
	WithConfirmed bool
	WithCommited  bool
	WithExpired   bool
	WithPending   bool

	WithConfirmations bool
}

type ConfirmTransactionParams struct {
	TxId           uuid.UUID
	UserId         uuid.UUID
	OrganizationId uuid.UUID
}

type CancelTransactionParams struct {
	TxId           uuid.UUID
	UserId         uuid.UUID
	OrganizationId uuid.UUID
}

type Repository interface {
	GetTransactions(ctx context.Context, params GetTransactionsParams) ([]*models.Transaction, error)
	CreateTransaction(ctx context.Context, tx models.Transaction) error
	UpdateTransaction(ctx context.Context, tx models.Transaction) error
	DeleteTransaction(ctx context.Context, tx models.Transaction) error

	ConfirmTransaction(ctx context.Context, params ConfirmTransactionParams) error
	CancelTransaction(ctx context.Context, params CancelTransactionParams) error
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

func (r *repositorySQL) GetTransactions(
	ctx context.Context,
	params GetTransactionsParams,
) ([]*models.Transaction, error) {
	var txs []*models.Transaction = make([]*models.Transaction, 0, len(params.Ids))
	err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := buildGetTransactionsQuery(params)

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch transactions data from database. %w", err)
		}

		defer func() {
			if cErr := rows.Close(); cErr != nil {
				err = errors.Join(fmt.Errorf("error close database rows. %w", cErr), err)
			}
		}()

		fmt.Println(query.ToSql())

		for rows.Next() {
			var (
				id             uuid.UUID
				description    string
				organizationId uuid.UUID
				amount         float64
				toAddr         []byte
				maxFeeAllowed  float64
				deadline       sql.NullTime
				createdAt      time.Time
				updatedAt      time.Time
				confirmedAt    sql.NullTime
				cancelledAt    sql.NullTime
				commitedAt     sql.NullTime

				createdById          uuid.UUID
				createdBySeed        []byte
				createdByCreatedAt   time.Time
				createdByActivatedAt sql.NullTime
				createdByIsAdmin     bool
			)

			if err = rows.Scan(
				&id,
				&description,
				&organizationId,
				&createdById,
				&amount,
				&toAddr,
				&maxFeeAllowed,
				&deadline,
				&createdAt,
				&updatedAt,
				&confirmedAt,
				&cancelledAt,
				&commitedAt,

				&createdBySeed,
				&createdByCreatedAt,
				&createdByActivatedAt,
				&createdByIsAdmin,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			tx := &models.Transaction{
				Id:             id,
				Description:    description,
				OrganizationId: organizationId,
				Amount:         amount,
				ToAddr:         toAddr,
				MaxFeeAllowed:  maxFeeAllowed,
				CreatedBy: &models.OrganizationUser{
					User: models.User{
						ID:        createdById,
						Bip39Seed: createdBySeed,
					},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}

			if deadline.Valid {
				tx.Deadline = deadline.Time
			}

			if confirmedAt.Valid {
				tx.ConfirmedAt = confirmedAt.Time
			}

			if commitedAt.Valid {
				tx.CommitedAt = commitedAt.Time
			}

			if cancelledAt.Valid {
				tx.CancelledAt = cancelledAt.Time
			}

			if createdByActivatedAt.Valid {
				tx.CreatedBy.Activated = true
			}

			txs = append(txs, tx)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func (r *repositorySQL) CreateTransaction(ctx context.Context, tx models.Transaction) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		columns := []string{
			"id",
			"description",
			"organization_id",
			"created_by",
			"amount",
			"to_addr",
			"max_fee_allowed",
			"created_at",
			"updated_at",
		}

		values := []any{
			tx.Id,
			tx.Description,
			tx.OrganizationId,
			tx.CreatedBy.ID,
			tx.Amount,
			tx.ToAddr,
			tx.MaxFeeAllowed,
			tx.CreatedAt,
			tx.CreatedAt,
		}

		if !tx.Deadline.IsZero() {
			columns = append(columns, "deadline")
			values = append(values, tx.Deadline)
		}

		query := sq.Insert("transactions").
			Columns(columns...).
			Values(values...).
			PlaceholderFormat(sq.Dollar)

		// todo add optional insertions

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error insert new transaction. %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repositorySQL) UpdateTransaction(ctx context.Context, tx models.Transaction) error {
	return nil
}

func (r *repositorySQL) DeleteTransaction(ctx context.Context, tx models.Transaction) error {
	return nil
}

func (r *repositorySQL) ConfirmTransaction(ctx context.Context, params ConfirmTransactionParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Update("transactions").
			SetMap(sq.Eq{
				"confirmed_at": time.Now(),
				"cancelled_at": nil,
			}).
			Where(sq.Eq{
				"id":              params.TxId,
				"organization_id": params.OrganizationId,
			}).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error update confirmed at. %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repositorySQL) CancelTransaction(ctx context.Context, params CancelTransactionParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Update("transactions").
			SetMap(sq.Eq{
				"cancelled_at": time.Now(),
				"confirmed_at": nil,
			}).
			Where(sq.Eq{
				"id":              params.TxId,
				"organization_id": params.OrganizationId,
			}).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error update confirmed at. %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func buildGetTransactionsQuery(params GetTransactionsParams) sq.SelectBuilder {
	query := sq.Select(
		`t.id,
		t.description,
		t.organization_id,
		t.created_by,
		t.amount,
		t.to_addr,
		t.max_fee_allowed,
		t.deadline,
		t.created_at,
		t.updated_at,
		t.confirmed_at,
		t.cancelled_at,
		t.commited_at,
		
		u.seed,
		u.created_at,
		u.activated_at,
		ou.is_admin`,
	).From("transactions as t").
		InnerJoin("users as u on u.id = t.created_by").
		InnerJoin(
			`organizations_users as ou on 
			u.id = ou.user_id and ou.organization_id = t.organization_id`,
		).
		Where(sq.Eq{
			"t.organization_id": params.OrganizationId,
		}).PlaceholderFormat(sq.Dollar)

	if len(params.Ids) > 0 {
		query = query.Where(sq.Eq{
			"t.id": params.Ids,
		})
	}

	if params.CreatedById != uuid.Nil {
		query = query.Where(sq.Eq{
			"t.created_by": params.CreatedById,
		})
	}

	if params.OrganizationId != uuid.Nil {
		query = query.Where(sq.Eq{
			"t.organization_id": params.OrganizationId,
		})
	}

	if params.To != nil {
		query = query.Where(sq.Eq{
			"t.to_addr": params.To,
		})
	}

	if params.WithExpired {
		query = query.Where(sq.LtOrEq{
			"t.deadline": time.Now(),
		})
	} else {
		query = query.Where(sq.Or{
			sq.GtOrEq{
				"t.deadline": time.Now(),
			},
			sq.Eq{
				"t.deadline": nil,
			},
		})
	}

	if params.WithCancelled {
		query = query.Where(sq.NotEq{
			"t.cancelled_at": nil,
		})
	}

	if params.WithConfirmed {
		query = query.Where(sq.NotEq{
			"t.confirmed_at": nil,
		})
	}

	if params.WithCommited {
		query = query.Where(sq.NotEq{
			"t.commited_at": nil,
		})
	}

	if params.Limit <= 0 || params.Limit > 50 {
		params.Limit = 50
	}

	if params.WithPending {
		query = query.Where(sq.Eq{
			"t.cancelled_at": nil,
			"t.commited_at":  nil,
			"t.confirmed_at": nil,
		})
	}

	if params.CursorId != uuid.Nil {
		query = query.Where(sq.Gt{
			"t.id": params.CursorId,
		})
	}

	query = query.Limit(uint64(params.Limit))

	return query
}