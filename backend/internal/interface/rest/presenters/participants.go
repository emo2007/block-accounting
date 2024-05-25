package presenters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/domain/hal"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
)

type ParticipantsPresenter interface {
	ResponseListParticipants(
		ctx context.Context,
		participants []models.OrganizationParticipant,
	) ([]byte, error)
	ResponseParticipant(
		ctx context.Context,
		participant models.OrganizationParticipant,
	) ([]byte, error)
}

type participantsPresenter struct{}

func NewParticipantsPresenter() ParticipantsPresenter {
	return new(participantsPresenter)
}

func (p *participantsPresenter) responseParticipant(
	ctx context.Context,
	participant models.OrganizationParticipant,
) (*hal.Resource, error) {
	domainParticipant := &domain.Participant{
		ID:        participant.Id().String(),
		Name:      participant.ParticipantName(),
		Position:  participant.Position(),
		CreatedAt: participant.CreatedDate().UnixMilli(),
		UpdatedAt: participant.UpdatedDate().UnixMilli(),
	}

	if !participant.DeletedDate().IsZero() {
		domainParticipant.DeletedAt = participant.DeletedDate().UnixMilli()
	}

	if user := participant.GetUser(); user != nil {
		if user.Credentails != nil {
			domainParticipant.Credentials = &domain.UserParticipantCredentials{
				Email:    user.Credentails.Email,
				Phone:    user.Credentails.Phone,
				Telegram: user.Credentails.Telegram,
			}
		}

		domainParticipant.IsUser = true
		domainParticipant.IsAdmin = user.IsAdmin()
		domainParticipant.IsOwner = user.IsOwner()
		domainParticipant.IsActive = user.Activated

	}

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	r := hal.NewResource(
		domainParticipant,
		"/organizations/"+organizationID.String()+"/participants"+domainParticipant.ID,
		hal.WithType("participant"),
	)

	return r, nil
}

func (p *participantsPresenter) responseParticipants(
	ctx context.Context,
	participants []models.OrganizationParticipant,
) (*hal.Resource, error) {
	resources := make([]*hal.Resource, len(participants))

	for i, pt := range participants {
		r, err := p.responseParticipant(ctx, pt)
		if err != nil {
			return nil, fmt.Errorf("error map participant to hal resource. %w", err)
		}

		resources[i] = r
	}

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch organization id from context. %w", err)
	}

	r := hal.NewResource(
		map[string][]*hal.Resource{
			"participants": resources,
		},
		"/organizations/"+organizationID.String()+"/participants",
		hal.WithType("participants"),
	)

	return r, nil
}

func (p *participantsPresenter) ResponseListParticipants(
	ctx context.Context,
	participants []models.OrganizationParticipant,
) ([]byte, error) {
	r, err := p.responseParticipants(ctx, participants)
	if err != nil {
		return nil, fmt.Errorf("error map participants to hal. %w", err)
	}

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal organization create response. %w", err)
	}

	return out, nil
}

func (p *participantsPresenter) ResponseParticipant(
	ctx context.Context,
	participant models.OrganizationParticipant,
) ([]byte, error) {
	r, err := p.responseParticipant(ctx, participant)
	if err != nil {
		return nil, fmt.Errorf("error map participant to hal resource. %w", err)
	}

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal organization create response. %w", err)
	}

	return out, nil
}
