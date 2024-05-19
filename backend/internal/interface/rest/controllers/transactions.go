package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/transactions"
)

type TransactionsController interface {
	New(w http.ResponseWriter, r *http.Request) ([]byte, error)
	List(w http.ResponseWriter, r *http.Request) ([]byte, error)
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

	ctx, cancel := context.WithTimeout(r.Context(), 30000*time.Second)
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
	panic("implement me!")
}
