package controllers

import (
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/usecase/interactors/transactions"
)

type TransactionsController interface {
	New(w http.ResponseWriter, r *http.Request) ([]byte, error)
	List(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type transactionsController struct {
	log          *slog.Logger
	txInteractor transactions.TransactionsInteractor
}

func NewTransactionsController(
	log *slog.Logger,
	txInteractor transactions.TransactionsInteractor,
) TransactionsController {
	return &transactionsController{
		log:          log,
		txInteractor: txInteractor,
	}
}

func (c *transactionsController) New(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	panic("implement me!")
}

func (c *transactionsController) List(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	panic("implement me!")
}
