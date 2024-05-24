package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/transactions"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionsController interface {
	New(w http.ResponseWriter, r *http.Request) ([]byte, error)
	List(w http.ResponseWriter, r *http.Request) ([]byte, error)
	UpdateStatus(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type transactionsController struct {
	log          *slog.Logger
	txInteractor transactions.TransactionsInteractor
	txPresenter  presenters.TransactionsPresenter
}

func NewTransactionsController(
	log *slog.Logger,
	txInteractor transactions.TransactionsInteractor,
	txPresenter presenters.TransactionsPresenter,
) TransactionsController {
	return &transactionsController{
		log:          log,
		txInteractor: txInteractor,
		txPresenter:  txPresenter,
	}
}

func (c *transactionsController) New(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.NewTransactionRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build new transaction request. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch organization ID from context. %w", err)
	}

	requestTx, err := c.txPresenter.RequestTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error build transaction model from request. %w", err)
	}

	c.log.Debug(
		"new reuqest",
		slog.Any("req", req),
		slog.String("org id", organizationID.String()),
	)

	tx, err := c.txInteractor.Create(ctx, transactions.CreateParams{
		OrganizationId: organizationID,
		Tx:             requestTx,
	})
	if err != nil {
		return nil, fmt.Errorf("error create new transaction. %w", err)
	}

	return c.txPresenter.ResponseNewTransaction(ctx, tx)
}

func (c *transactionsController) List(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.ListTransactionsRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build new transaction request. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch organization ID from context. %w", err)
	}

	ids := make(uuid.UUIDs, len(req.IDs))

	for i, id := range req.IDs {
		txUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("error parse tx id. %w", err)
		}

		ids[i] = txUUID
	}

	var toAddr []byte

	if req.To != "" {
		toAddr = common.HexToAddress(req.To).Bytes()
	}

	var createdBy uuid.UUID
	if req.CreatedBy != "" {
		createdBy, err = uuid.Parse(req.CreatedBy)
		if err != nil {
			return nil, fmt.Errorf("error parse created by id. %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	txs, err := c.txInteractor.List(ctx, transactions.ListParams{
		IDs:            ids,
		OrganizationID: organizationID,
		To:             toAddr,

		CreatedBy: createdBy,

		Limit:  int64(req.Limit),
		Cursor: req.Cursor,

		WithCancelled: req.Cancelled,
		WithConfirmed: req.Confirmed,
		WithCommited:  req.Commited,
		WithExpired:   req.Expired,
		WithPending:   req.Pending,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch organizations list. %w", err)
	}

	return c.txPresenter.ResponseListTransactions(ctx, txs.Txs, txs.NextCursor)
}

func (c *transactionsController) UpdateStatus(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.UpdateTransactionStatusRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build new transaction request. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch organization ID from context. %w", err)
	}

	txIDstr := chi.URLParam(r, "tx_id")

	var txID uuid.UUID

	if txID, err = uuid.Parse(txIDstr); err != nil {
		return nil, fmt.Errorf("error parse tx id. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var tx *models.Transaction

	if req.Cancel {
		tx, err = c.txInteractor.Cancel(ctx, transactions.CancelParams{
			TxID:           txID,
			OrganizationID: organizationID,
		})
		if err != nil {
			return nil, fmt.Errorf("error cancel transaction. %w", err)
		}
	} else if req.Confirm {
		tx, err = c.txInteractor.Confirm(ctx, transactions.ConfirmParams{
			TxID:           txID,
			OrganizationID: organizationID,
		})
		if err != nil {
			return nil, fmt.Errorf("error cancel transaction. %w", err)
		}
	} else {
		return nil, fmt.Errorf("error new status required")
	}

	return c.txPresenter.ResponseNewTransaction(ctx, tx)
}

// todo creates a new payout
func (c *transactionsController) NewPayout(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return nil, nil
}
