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

	ResponseMultisigs(ctx context.Context, msgs []models.Multisig) ([]byte, error)

	ResponsePayrolls(ctx context.Context, payrolls []models.Payroll) ([]byte, error)
}

type transactionsPresenter struct {
	participantsPresenter ParticipantsPresenter
}

func NewTransactionsPresenter() TransactionsPresenter {
	return &transactionsPresenter{
		participantsPresenter: NewParticipantsPresenter(),
	}
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
		CreatedAt:      time.Now(),
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
		hal.WithType("transactions"),
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

type Multisig struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Owners *hal.Resource `json:"owners"`
}

func (c *transactionsPresenter) ResponseMultisigs(ctx context.Context, msgs []models.Multisig) ([]byte, error) {
	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	outArray := make([]Multisig, len(msgs))

	for i, m := range msgs {
		mout := Multisig{
			ID:    m.ID.String(),
			Title: m.Title,
		}

		partOut, err := c.participantsPresenter.ResponseParticipantsHal(ctx, m.Owners)
		if err != nil {
			return nil, err
		}

		mout.Owners = partOut

		outArray[i] = mout
	}

	txsResource := map[string]any{"multisigs": outArray}

	r := hal.NewResource(
		txsResource,
		"/organizations/"+organizationID.String()+"/multisig",
		hal.WithType("multisigs"),
	)

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal multisigs to hal resource. %w", err)
	}

	return out, nil
}

type Payroll struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (c *transactionsPresenter) ResponsePayrolls(
	ctx context.Context,
	payrolls []models.Payroll,
) ([]byte, error) {
	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	outArray := make([]Payroll, len(payrolls))

	for i, pr := range payrolls {
		outArray[i] = Payroll{
			ID:        pr.ID.String(),
			Title:     pr.Title,
			CreatedAt: pr.CreatedAt.UnixMilli(),
			UpdatedAt: pr.UpdatedAt.UnixMilli(),
		}
	}

	txsResource := map[string]any{"payrolls": outArray}

	r := hal.NewResource(
		txsResource,
		"/organizations/"+organizationID.String()+"/payrolls",
		hal.WithType("paurolls"),
	)

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal payrolls to hal resource. %w", err)
	}

	return out, nil
}
