package transactions

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

	WithPending bool

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
	List(ctx context.Context, params ListParams) (*ListResult, error)
	Create(ctx context.Context, params CreateParams) (*models.Transaction, error)
	Confirm(ctx context.Context, params ConfirmParams) (*models.Transaction, error)
	Cancel(ctx context.Context, params CancelParams) (*models.Transaction, error)
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

type txsListCursor struct {
	Id uuid.UUID `json:"id"`
}

func newTxsListCursor(id ...uuid.UUID) *txsListCursor {
	if len(id) > 0 {
		return &txsListCursor{id[0]}
	}

	return new(txsListCursor)
}

func (c *txsListCursor) encode() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal reaction id. %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil

}

func (c *txsListCursor) decode(s string) error {
	if c == nil {
		return nil
	}

	token, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("error decode token. %w", err)
	}

	return json.Unmarshal(token, c)
}

type ListResult struct {
	Txs        []*models.Transaction
	NextCursor string
}

func (i *transactionsInteractor) List(ctx context.Context, params ListParams) (*ListResult, error) {
	if params.Limit == 0 {
		params.Limit = 50
	}

	cursor := newTxsListCursor()

	if params.Cursor != "" {
		if err := cursor.decode(params.Cursor); err != nil {
			return nil, fmt.Errorf("error decode cursor value. %w", err) // maybe just log error?
		}
	}

	txs, err := i.txRepo.GetTransactions(ctx, transactions.GetTransactionsParams{
		Ids:            params.IDs,
		OrganizationId: params.OrganizationID,
		CreatedById:    params.CreatedBy,
		To:             params.To,
		Limit:          params.Limit,
		CursorId:       cursor.Id,
		WithCancelled:  params.WithCancelled,
		WithConfirmed:  params.WithConfirmed,
		WithCommited:   params.WithCommited,
		WithExpired:    params.WithExpired,
		WithPending:    params.WithPending,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch transaction from repository. %w", err)
	}

	var nextCursor string

	if len(txs) >= 50 || len(txs) >= int(params.Limit) {
		cursor.Id = txs[len(txs)-1].Id
		if nextCursor, err = cursor.encode(); err != nil {
			return nil, fmt.Errorf("error encode next page token. %w", err) // maybe just log error?
		}
	}

	return &ListResult{
		Txs:        txs,
		NextCursor: nextCursor,
	}, nil
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

	if tx.Id == uuid.Nil {
		tx.Id = uuid.Must(uuid.NewV7())
	}

	participant, err := i.orgInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             user.Id(),
		OrganizationID: params.OrganizationId,
		ActiveOnly:     true,
		UsersOnly:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch actor prticipant. %w", err)
	}

	tx.CreatedBy = participant.GetUser()
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = tx.CreatedAt

	if err = i.txRepo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("error create new tx. %w", err)
	}

	return &tx, nil
}

func (i *transactionsInteractor) Confirm(ctx context.Context, params ConfirmParams) (*models.Transaction, error) {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	participant, err := i.orgInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             user.Id(),
		OrganizationID: params.OrganizationID,
		ActiveOnly:     true,
		UsersOnly:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch actor prticipant. %w", err)
	}

	if !participant.IsAdmin() {
		return nil, fmt.Errorf("error not enouth rights. %w", organizations.ErrorUnauthorizedAccess)
	}

	if err := i.txRepo.ConfirmTransaction(ctx, transactions.ConfirmTransactionParams{
		TxId:           params.TxID,
		OrganizationId: params.OrganizationID,
		UserId:         participant.Id(),
	}); err != nil {
		return nil, fmt.Errorf("error confirm transaction. %w", err)
	}

	tx, err := i.txRepo.GetTransactions(ctx, transactions.GetTransactionsParams{
		Ids:            uuid.UUIDs{params.TxID},
		OrganizationId: params.OrganizationID,
		Limit:          1,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch transaction. %w", err)
	}

	if len(tx) == 0 {
		return nil, fmt.Errorf("error tx not found")
	}

	return tx[0], nil
}

func (i *transactionsInteractor) Cancel(ctx context.Context, params CancelParams) (*models.Transaction, error) {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	participant, err := i.orgInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             user.Id(),
		OrganizationID: params.OrganizationID,
		ActiveOnly:     true,
		UsersOnly:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch actor prticipant. %w", err)
	}

	if !participant.IsAdmin() {
		return nil, fmt.Errorf("error not enouth rights. %w", organizations.ErrorUnauthorizedAccess)
	}

	if err := i.txRepo.CancelTransaction(ctx, transactions.CancelTransactionParams{
		TxId:           params.TxID,
		OrganizationId: params.OrganizationID,
		UserId:         participant.Id(),
	}); err != nil {
		return nil, fmt.Errorf("error confirm transaction. %w", err)
	}

	tx, err := i.txRepo.GetTransactions(ctx, transactions.GetTransactionsParams{
		Ids:            uuid.UUIDs{params.TxID},
		OrganizationId: params.OrganizationID,
		Limit:          1,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch transaction. %w", err)
	}

	if len(tx) == 0 {
		return nil, fmt.Errorf("error tx not found")
	}

	return tx[0], nil
}
