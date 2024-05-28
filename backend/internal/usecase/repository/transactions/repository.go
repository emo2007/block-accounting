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
	"github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	"github.com/google/uuid"
)

type GetTransactionsParams struct {
	Ids            uuid.UUIDs
	OrganizationId uuid.UUID
	Limit          int64
	CursorId       uuid.UUID

	Pending        bool
	ReadyToConfirm bool
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

	AddMultisig(ctx context.Context, multisig models.Multisig) error
	ListMultisig(ctx context.Context, params ListMultisigsParams) ([]models.Multisig, error)
	ConfirmMultisig(ctx context.Context, params ConfirmMultisigParams) error

	AddPayrollContract(ctx context.Context, params AddPayrollContract) error
	ListPayrolls(ctx context.Context, params ListPayrollsParams) ([]models.Payroll, error)
}

type repositorySQL struct {
	db      *sql.DB
	orgRepo organizations.Repository
}

func NewRepository(db *sql.DB, orgRepo organizations.Repository) Repository {
	return &repositorySQL{
		db:      db,
		orgRepo: orgRepo,
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
		t.confirmations_required,
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

	if params.OrganizationId != uuid.Nil {
		query = query.Where(sq.Eq{
			"t.organization_id": params.OrganizationId,
		})
	}

	if params.Pending {
		query = query.InnerJoin("multisig_confirmations as mc on mc.multisig_id = t.multisig_id").Where(
			sq.Lt{},
		)
	}

	query = query.Limit(uint64(params.Limit))

	return query
}

func (r *repositorySQL) AddMultisig(
	ctx context.Context,
	multisig models.Multisig,
) error {
	return sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Insert("multisigs").Columns(
			"id",
			"organization_id",
			"title",
			"address",
			"confirmations",
			"created_at",
			"updated_at",
		).Values(
			multisig.ID,
			multisig.OrganizationID,
			multisig.Title,
			multisig.Address,
			multisig.ConfirmationsRequired,
			multisig.CreatedAt,
			multisig.UpdatedAt,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error insert multisig data. %w", err)
		}

		for _, owner := range multisig.Owners {
			addOwnerQuery := sq.Insert("multisig_owners").Columns(
				"multisig_id",
				"owner_id",
				"created_at",
				"updated_at",
			).Values(
				multisig.ID,
				owner.Id(),
				multisig.CreatedAt,
				multisig.UpdatedAt,
			).PlaceholderFormat(sq.Dollar)

			if _, err := addOwnerQuery.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
				return fmt.Errorf("error insert multisig owner data. %w", err)
			}
		}

		return nil
	})
}

type ListMultisigsParams struct {
	IDs            uuid.UUIDs
	OrganizationID uuid.UUID
}

func (r *repositorySQL) ListMultisig(
	ctx context.Context,
	params ListMultisigsParams,
) ([]models.Multisig, error) {
	msgs := make([]models.Multisig, 0)

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Select(
			"id",
			"organization_id",
			"title",
			"address",
			"confirmations",
			"created_at",
			"updated_at",
		).From("multisigs").Where(sq.Eq{
			"organization_id": params.OrganizationID,
		}).PlaceholderFormat(sq.Dollar)

		if len(params.IDs) > 0 {
			query = query.Where(sq.Eq{
				"id": params.IDs,
			})
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch multisigs from database. %w", err)
		}

		defer rows.Close()

		msgsTmp := make([]*models.Multisig, 0)

		for rows.Next() {
			var (
				id             uuid.UUID
				organizationID uuid.UUID
				address        []byte
				title          string
				confirmations  int
				createdAt      time.Time
				updatedAt      time.Time
			)

			if err = rows.Scan(
				&id,
				&organizationID,
				&title,
				&address,
				&confirmations,
				&createdAt,
				&updatedAt,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			msgsTmp = append(msgsTmp, &models.Multisig{
				ID:                    id,
				Title:                 title,
				Address:               address,
				OrganizationID:        organizationID,
				ConfirmationsRequired: confirmations,
				CreatedAt:             createdAt,
				UpdatedAt:             updatedAt,
			})
		}

		for _, m := range msgsTmp {
			owners, err := r.fetchOwners(ctx, fetchOwnersParams{
				OrganizationID: params.OrganizationID,
				MultisigID:     m.ID,
			})
			if err != nil {
				return err
			}

			m.Owners = owners

			msgs = append(msgs, *m)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return msgs, nil
}

type ConfirmMultisigParams struct {
	MultisigID      uuid.UUID
	OrganizationsID uuid.UUID
	CinfirmedBy     *models.OrganizationUser
	ConfirmedAt     time.Time
}

func (r *repositorySQL) ConfirmMultisig(ctx context.Context, params ConfirmMultisigParams) error {
	return sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		deleteOldQuery := sq.Delete("multisig_confirmations").
			Where(sq.Eq{
				"multisig_id": params.MultisigID,
				"owner_id":    params.CinfirmedBy.Id(),
			}).
			PlaceholderFormat(sq.Dollar)

		if _, err := deleteOldQuery.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error delete old multisig confirmation. %w", err)
		}

		query := sq.Insert("multisig_confirmations").
			Columns(
				"multisig_id",
				"owner_id",
				"created_at",
			).
			Values(
				params.MultisigID,
				params.CinfirmedBy.Id(),
				params.ConfirmedAt,
			).
			PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add multisig confirmation. %w", err)
		}

		return nil
	})
}

type AddPayrollContract struct {
	ID             uuid.UUID
	Title          string
	Description    string
	Address        []byte
	Payload        []byte
	OrganizationID uuid.UUID
	MultisigID     uuid.UUID
	CreatedAt      time.Time
}

func (r *repositorySQL) AddPayrollContract(
	ctx context.Context,
	params AddPayrollContract,
) error {
	return sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Insert("payrolls").
			Columns(
				"id",
				"title",
				"description",
				"address",
				"payload",
				"organization_id",
				"multisig_id",
				"created_at",
			).
			Values(
				params.ID,
				params.Title,
				params.Description,
				params.Address,
				params.Payload,
				params.OrganizationID,
				params.MultisigID,
				params.CreatedAt,
			).
			PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add new payroll contract. %w", err)
		}

		return nil
	})
}

type fetchOwnersParams struct {
	MultisigID     uuid.UUID
	OrganizationID uuid.UUID
}

func (r *repositorySQL) fetchOwners(ctx context.Context, params fetchOwnersParams) ([]models.OrganizationParticipant, error) {
	owners := make([]models.OrganizationParticipant, 0)

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Select("owner_id").From("multisig_owners").Where(sq.Eq{
			"multisig_id": params.MultisigID,
		}).PlaceholderFormat(sq.Dollar)

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch multisigs owners from database. %w", err)
		}

		defer rows.Close()

		ids := make(uuid.UUIDs, 0)

		for rows.Next() {
			var ownerId uuid.UUID

			if err = rows.Scan(&ownerId); err != nil {
				return err
			}

			ids = append(ids, ownerId)
		}

		owners, err = r.orgRepo.Participants(ctx, organizations.ParticipantsParams{
			OrganizationId: params.OrganizationID,
			Ids:            ids,
			UsersOnly:      true,
		})
		if err != nil {
			return fmt.Errorf("error fetch owners as participants. %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return owners, nil
}

type ListPayrollsParams struct {
	IDs            []uuid.UUID
	Limit          int64
	OrganizationID uuid.UUID
}

func (r *repositorySQL) ListPayrolls(ctx context.Context, params ListPayrollsParams) ([]models.Payroll, error) {
	payrolls := make([]models.Payroll, 0, len(params.IDs))

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Select(
			"id",
			"title",
			"address",
			"organization_id",
			"multisig_id",
			"created_at",
			"updated_at",
		).Where(sq.Eq{
			"organization_id": params.OrganizationID,
		}).PlaceholderFormat(sq.Dollar)

		if params.Limit <= 0 {
			params.Limit = 100
		}

		if len(params.IDs) > 0 {
			query = query.Where(sq.Eq{
				"id": params.IDs,
			})
		}

		query = query.Limit(uint64(params.Limit))

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch payrolls from database. %w", err)
		}

		defer rows.Close() // todo check error

		for rows.Next() {
			var (
				id             uuid.UUID
				title          string
				address        []byte
				organizationId uuid.UUID
				multisigId     uuid.UUID
				createdAt      time.Time
				updatedAt      time.Time
			)

			if err = rows.Scan(
				&id,
				&title,
				&address,
				&organizationId,
				&multisigId,
				&createdAt,
				&updatedAt,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			payrolls = append(payrolls, models.Payroll{
				ID:             id,
				Title:          title,
				Address:        address,
				OrganizationID: organizationId,
				MultisigID:     multisigId,
				CreatedAt:      createdAt,
				UpdatedAt:      updatedAt,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return payrolls, nil
}
