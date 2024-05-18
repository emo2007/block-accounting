package transactions

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/transactions"
	"github.com/google/uuid"
)

type ListParams struct {
	IDs            uuid.UUIDs
	OrganizationID uuid.UUID
	CreatedBy      uuid.UUID

	To []byte

	Limit  int64
	Cursor string

	WithCancelled bool
	WithConfirmed bool
	WithCommited  bool
	WithExpired   bool

	WithConfirmations bool
}

type CreateParams struct {
	Tx             models.Transaction
	OrganizationId uuid.UUID
}

type ConfirmParams struct {
	TxID           uuid.UUID
	OrganizationID uuid.UUID
}

type CancelParams struct {
	TxID           uuid.UUID
	OrganizationID uuid.UUID
	Cause          string
}

type TransactionsInteractor interface {
	List(ctx context.Context, params ListParams) ([]*models.Transaction, error)
	Create(ctx context.Context, params CreateParams) (*models.Transaction, error)
	Confirm(ctx context.Context, params ConfirmParams) (*models.Transaction, error)
	Cancel(ctx context.Context, params CancelParams) (*models.Transaction, error)
	// TODO delete
	// TODO update
}

type transactionsInteractor struct {
	log           *slog.Logger
	txRepo        transactions.Repository
	orgInteractor organizations.OrganizationsInteractor
}

func NewTransactionsInteractor(
	log *slog.Logger,
	txRepo transactions.Repository,
	orgInteractor organizations.OrganizationsInteractor,
) TransactionsInteractor {
	return &transactionsInteractor{
		log:           log,
		txRepo:        txRepo,
		orgInteractor: orgInteractor,
	}
}

func (i *transactionsInteractor) List(ctx context.Context, params ListParams) ([]*models.Transaction, error) {
	filters := make([]transactions.GetTransactionsFilter, 0)

	if params.WithCancelled {
		filters = append(filters, transactions.GetFilterCancelled)
	}

	if params.WithConfirmed {
		filters = append(filters, transactions.GetFilterConfirmed)
	}

	if params.WithCommited {
		filters = append(filters, transactions.GetFilterCommited)
	}

	txs, err := i.txRepo.GetTransactions(ctx, transactions.GetTransactionsParams{
		Ids:            params.IDs,
		OrganizationId: params.OrganizationID,
		CreatedById:    params.CreatedBy,
		Filters:        filters,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch transaction from repository. %w", err)
	}

	return txs, nil
}

func (i *transactionsInteractor) Create(
	ctx context.Context,
	params CreateParams,
) (*models.Transaction, error) {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	tx := params.Tx

	participant, err := i.orgInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:         user.Id(),
		ActiveOnly: true,
		UsersOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch actor prticipant. %w", err)
	}

	tx.CreatedBy = participant.GetUser()
	tx.CreatedAt = time.Now()

	if err = i.txRepo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("error create new tx. %w", err)
	}

	return &tx, nil
}

func (i *transactionsInteractor) Confirm(ctx context.Context, params ConfirmParams) (*models.Transaction, error) {
	panic("implement me!")
}
func (i *transactionsInteractor) Cancel(ctx context.Context, params CancelParams) (*models.Transaction, error) {
	panic("implement me!")
}
