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
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/google/uuid"
)

type ParticipantsController interface {
	List(w http.ResponseWriter, r *http.Request) ([]byte, error)
	New(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type participantsController struct {
	log             *slog.Logger
	orgInteractor   organizations.OrganizationsInteractor
	usersInteractor users.UsersInteractor

	presenter presenters.ParticipantsPresenter
}

func NewParticipantsController(
	log *slog.Logger,
	orgInteractor organizations.OrganizationsInteractor,
	usersInteractor users.UsersInteractor,
	presenter presenters.ParticipantsPresenter,
) ParticipantsController {
	return &participantsController{
		log:             log,
		orgInteractor:   orgInteractor,
		usersInteractor: usersInteractor,
		presenter:       presenter,
	}
}

func (c *participantsController) List(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.ListParticipantsRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build list participants request. %w", err)
	}

	user, err := ctxmeta.User(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	participant, err := c.orgInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             user.Id(),
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch user participant. %w", err)
	}

	if !participant.IsActive() {
		return nil, fmt.Errorf("error participant is inactive")
	}

	ids := make(uuid.UUIDs, len(req.IDs))
	for i, id := range req.IDs {
		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("error parse participant id. %w", err)
		}

		ids[i] = uid
	}

	participants, err := c.orgInteractor.Participants(ctx, organizations.ParticipantsParams{
		IDs:            ids,
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch participants. %w", err)
	}

	return c.presenter.ResponseListParticipants(ctx, participants)
}

func (c *participantsController) New(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.AddEmployeeRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build list participants request. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(r.Context())
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	participant, err := c.orgInteractor.AddParticipant(ctx, organizations.AddParticipantParams{
		OrganizationID: organizationID,
		Name:           req.Name,
		Position:       req.Position,
		WalletAddress:  req.WalletAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error create new participant. %w", err)
	}

	return c.presenter.ResponseParticipant(ctx, participant)
}
