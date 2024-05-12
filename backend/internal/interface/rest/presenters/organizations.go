package presenters

import (
	"encoding/json"
	"fmt"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
)

type OrganizationsPresenter interface {
	ResponseCreate(organization *models.Organization) ([]byte, error)
	ResponseList(orgs []*models.Organization, nextCursor string) ([]byte, error)
	Organizations(orgs []*models.Organization) []domain.Organization
}

type organizationsPresenter struct {
}

func NewOrganizationsPresenter() OrganizationsPresenter {
	return &organizationsPresenter{}
}

func (p *organizationsPresenter) ResponseCreate(o *models.Organization) ([]byte, error) {
	resp := &domain.NewOrganizationResponse{
		Organization: domain.Organization{
			Id:        o.ID.String(),
			Name:      o.Name,
			Address:   o.Address,
			CreatedAt: uint64(o.CreatedAt.UnixMilli()),
			UpdatedAt: uint64(o.UpdatedAt.UnixMilli()),
		},
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshal organization create response. %w", err)
	}

	return out, nil
}

func (p *organizationsPresenter) ResponseList(orgs []*models.Organization, nextCursor string) ([]byte, error) {
	resp := &domain.ListOrganizationsResponse{
		Collection: domain.Collection[domain.Organization]{
			Items: p.Organizations(orgs),
			Pagination: domain.Pagination{
				NextCursor: nextCursor,
				TotalItems: uint32(len(orgs)),
			},
		},
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshal organizations list response. %w", err)
	}

	return out, nil
}

func (p *organizationsPresenter) Organizations(orgs []*models.Organization) []domain.Organization {
	out := make([]domain.Organization, len(orgs))

	for i, o := range orgs {
		out[i] = domain.Organization{
			Id:        o.ID.String(),
			Name:      o.Name,
			Address:   o.Address,
			CreatedAt: uint64(o.CreatedAt.UnixMilli()),
			UpdatedAt: uint64(o.UpdatedAt.UnixMilli()),
		}
	}

	return out
}
