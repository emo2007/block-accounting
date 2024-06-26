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
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/chain"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/transactions"
	"github.com/ethereum/go-ethereum/common"
)

// TODO по хорошему это уебищу надо разносить, но ни времени ни сил пока нет
// в рамках рефакторинка не забыть
// TransactionsController | ChainController
type TransactionsController interface {
	New(w http.ResponseWriter, r *http.Request) ([]byte, error)
	List(w http.ResponseWriter, r *http.Request) ([]byte, error)
	UpdateStatus(w http.ResponseWriter, r *http.Request) ([]byte, error)

	NewPayroll(w http.ResponseWriter, r *http.Request) ([]byte, error)
	ConfirmPayroll(w http.ResponseWriter, r *http.Request) ([]byte, error)
	ListPayrolls(w http.ResponseWriter, r *http.Request) ([]byte, error)

	NewMultisig(w http.ResponseWriter, r *http.Request) ([]byte, error)
	ListMultisigs(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type transactionsController struct {
	log                     *slog.Logger
	txInteractor            transactions.TransactionsInteractor
	txPresenter             presenters.TransactionsPresenter
	chainInteractor         chain.ChainInteractor
	organizationsInteractor organizations.OrganizationsInteractor
}

func NewTransactionsController(
	log *slog.Logger,
	txInteractor transactions.TransactionsInteractor,
	txPresenter presenters.TransactionsPresenter,
	chainInteractor chain.ChainInteractor,
	organizationsInteractor organizations.OrganizationsInteractor,
) TransactionsController {
	return &transactionsController{
		log:                     log,
		txInteractor:            txInteractor,
		txPresenter:             txPresenter,
		chainInteractor:         chainInteractor,
		organizationsInteractor: organizationsInteractor,
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

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	txs, err := c.txInteractor.List(ctx, transactions.ListParams{
		OrganizationID: organizationID,
		Limit:          int64(req.Limit),
		Cursor:         req.Cursor,
		Pending:        req.Pending,
		ReadyToConfirm: req.ReadyToConfirm,
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

func (c *transactionsController) NewMultisig(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.NewMultisigRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build new transaction request. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch organization ID from context. %w", err)
	}

	c.log.Debug(
		"new multisig request",
		slog.String("org id", organizationID.String()),
		slog.Any("req", req),
	)

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	ownersPKs := make([][]byte, len(req.Owners))

	for i, pk := range req.Owners {
		ownersPKs[i] = common.Hex2Bytes(pk.PublicKey[2:])
	}

	if req.Confirmations <= 0 {
		req.Confirmations = 1
	}

	participants, err := c.organizationsInteractor.Participants(ctx, organizations.ParticipantsParams{
		PKs:            ownersPKs,
		OrganizationID: organizationID,
		UsersOnly:      true,
		ActiveOnly:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch participants by pks. %w", err)
	}

	if err := c.chainInteractor.NewMultisig(ctx, chain.NewMultisigParams{
		Title:         req.Title,
		Owners:        participants,
		Confirmations: req.Confirmations,
	}); err != nil {
		return nil, fmt.Errorf("error deploy multisig. %w", err)
	}

	return presenters.ResponseOK()
}

func (s *transactionsController) ListMultisigs(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	organizationID, err := ctxmeta.OrganizationId(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch organization ID from context. %w", err)
	}

	msgs, err := s.chainInteractor.ListMultisigs(r.Context(), chain.ListMultisigsParams{
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, err
	}

	return s.txPresenter.ResponseMultisigs(r.Context(), msgs)
}

func (c *transactionsController) NewPayroll(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.NewPayrollRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("erropr fetch organization id from context. %w", err)
	}

	c.log.Debug(
		"NewPayrollRequest",
		slog.Any("req", req),
		slog.String("org id", organizationID.String()),
	)

	multisigID, err := uuid.Parse(req.MultisigID)
	if err != nil {
		return nil, fmt.Errorf("error invalid ")
	}

	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	userParticipant, err := c.organizationsInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             user.Id(),
		OrganizationID: organizationID,
		// TODO fetch REAL first admin
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch user protocpant. %w", err)
	}

	if !userParticipant.IsOwner() {
		return nil, fmt.Errorf("only owner can create payrolls")
	}

	firstAdmin, err := c.organizationsInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             user.Id(),
		OrganizationID: organizationID,
		// TODO fetch REAL first admin
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch first admin. %w", err)
	}

	if !firstAdmin.IsOwner() {
		return nil, fmt.Errorf("invalid first admin. not owner")
	}

	err = c.chainInteractor.PayrollDeploy(ctx, chain.PayrollDeployParams{
		MultisigID: multisigID,
		FirstAdmin: firstAdmin,
		Title:      req.Title,
	})
	if err != nil {
		return nil, fmt.Errorf("error create new payroll contract. %w", err)
	}

	return presenters.ResponseOK()
}

func (c *transactionsController) ListPayrolls(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.ListPayrollsRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("erropr fetch organization id from context. %w", err)
	}

	ids := make(uuid.UUIDs, len(req.IDs))
	for i, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("error parse payroll id. %w", err)
		}

		ids[i] = id
	}

	payrolls, err := c.chainInteractor.ListPayrolls(ctx, chain.ListPayrollsParams{
		OrganizationID: organizationID,
		IDs:            ids,
		Limit:          int(req.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch payrolls. %w", err)
	}

	return c.txPresenter.ResponsePayrolls(ctx, payrolls)
}

func (c *transactionsController) SetSalary(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	// req, err := presenters.CreateRequest[domain.]()

	return presenters.ResponseOK()
}

func (c *transactionsController) ConfirmPayroll(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	return nil, nil
}
