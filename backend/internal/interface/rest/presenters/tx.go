package presenters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/domain/hal"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrorInvalidHexAddress = errors.New("error invalid hex address")
)

type TransactionsPresenter interface {
	RequestTransaction(ctx context.Context, r *domain.NewTransactionRequest) (models.Transaction, error)
	ResponseTransaction(ctx context.Context, tx *models.Transaction) (*hal.Resource, error)
	ResponseNewTransaction(ctx context.Context, tx *models.Transaction) ([]byte, error)
	ResponseTransactionsArray(ctx context.Context, txs []*models.Transaction) ([]*hal.Resource, error)
	ResponseListTransactions(ctx context.Context, txs []*models.Transaction, cursor string) ([]byte, error)
}

type transactionsPresenter struct {
}

func NewTransactionsPresenter() TransactionsPresenter {
	return &transactionsPresenter{}
}

// RequestTransaction returns a Transaction model WITHOUT CreatedBy user set. CreatedAt set as time.Now()
func (p *transactionsPresenter) RequestTransaction(
	ctx context.Context,
	r *domain.NewTransactionRequest,
) (models.Transaction, error) {
	if !common.IsHexAddress(r.ToAddr) {
		return models.Transaction{}, ErrorInvalidHexAddress
	}

	toAddress := common.HexToAddress(r.ToAddr)

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	tx := models.Transaction{
		OrganizationId: organizationID,
		Description:    r.Description,
		Amount:         r.Amount,
		ToAddr:         toAddress.Bytes(),
		MaxFeeAllowed:  r.MaxFeeAllowed,
		CreatedAt:      time.Now(),
	}

	if r.Deadline > 0 {
		tx.Deadline = time.UnixMilli(r.Deadline)
	}

	return tx, nil
}

func (c *transactionsPresenter) ResponseTransaction(
	ctx context.Context,
	tx *models.Transaction,
) (*hal.Resource, error) {
	r := &domain.Transaction{
		Id:             tx.Id.String(),
		Description:    tx.Description,
		OrganizationId: tx.OrganizationId.String(),
		CreatedBy:      tx.CreatedBy.Id().String(),
		Amount:         tx.Amount,
		MaxFeeAllowed:  tx.MaxFeeAllowed,
		Status:         tx.Status,
		CreatedAt:      tx.CreatedAt.UnixMilli(),
		UpdatedAt:      tx.UpdatedAt.UnixMilli(),
	}

	addr := common.BytesToAddress(tx.ToAddr)

	r.ToAddr = addr.String()

	if !tx.ConfirmedAt.IsZero() {
		r.ConfirmedAt = tx.ConfirmedAt.UnixMilli()
	}

	if !tx.CancelledAt.IsZero() {
		r.CancelledAt = tx.CancelledAt.UnixMilli()
	}

	if !tx.CommitedAt.IsZero() {
		r.CommitedAt = tx.CommitedAt.UnixMilli()
	}

	if !tx.Deadline.IsZero() {
		r.Deadline = tx.Deadline.UnixMilli()
	}

	return hal.NewResource(
		r,
		"/organizations/"+tx.OrganizationId.String()+"/transactions/"+tx.Id.String(),
		hal.WithType("transaction"),
	), nil
}

func (p *transactionsPresenter) ResponseTransactionsArray(
	ctx context.Context,
	txs []*models.Transaction,
) ([]*hal.Resource, error) {
	out := make([]*hal.Resource, len(txs))

	for i, tx := range txs {
		r, err := p.ResponseTransaction(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("error map transaction to hal reource. %w", err)
		}

		out[i] = r
	}

	return out, nil
}

func (p *transactionsPresenter) ResponseListTransactions(
	ctx context.Context,
	txs []*models.Transaction,
	cursor string,
) ([]byte, error) {
	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	arr, err := p.ResponseTransactionsArray(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("error map transactions list to resource array. %w", err)
	}

	txsResource := map[string]any{"transactions": arr}

	if cursor != "" {
		txsResource["next_cursor"] = cursor
	}

	r := hal.NewResource(
		txsResource,
		"/organizations/"+organizationID.String()+"/transactions",
		hal.WithType("organizations"),
	)

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal tx to hal resource. %w", err)
	}

	return out, nil
}

func (c *transactionsPresenter) ResponseNewTransaction(
	ctx context.Context,
	tx *models.Transaction,
) ([]byte, error) {
	dtoTx, err := c.ResponseTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("error map tx to dto. %w", err)
	}

	out, err := json.Marshal(dtoTx)
	if err != nil {
		return nil, fmt.Errorf("error marshal tx to hal resource. %w", err)
	}

	return out, nil
}
