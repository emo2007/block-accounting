package presenters

import (
	"encoding/json"
	"fmt"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/domain/hal"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
)

type OrganizationsPresenter interface {
	ResponseCreate(organization *models.Organization) ([]byte, error)
	ResponseList(orgs []*models.Organization, nextCursor string) ([]byte, error)
	Organizations(orgs []*models.Organization) []*hal.Resource
}

type organizationsPresenter struct {
}

func NewOrganizationsPresenter() OrganizationsPresenter {
	return &organizationsPresenter{}
}

func (p *organizationsPresenter) ResponseCreate(o *models.Organization) ([]byte, error) {
	org := domain.Organization{
		Id:        o.ID.String(),
		Name:      o.Name,
		Address:   o.Address,
		CreatedAt: uint64(o.CreatedAt.UnixMilli()),
		UpdatedAt: uint64(o.UpdatedAt.UnixMilli()),
	}

	r := hal.NewResource(
		org,
		"/organizations/"+org.Id,
		hal.WithType("organization"),
	)

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal organization create response. %w", err)
	}

	return out, nil
}

func (p *organizationsPresenter) ResponseList(orgs []*models.Organization, nextCursor string) ([]byte, error) {
	dtoOrgs := domain.Collection[*hal.Resource]{
		Items: p.Organizations(orgs),
		Pagination: domain.Pagination{
			NextCursor: nextCursor,
			TotalItems: uint32(len(orgs)),
		},
	}

	r := hal.NewResource(
		dtoOrgs,
		"/organizations",
		hal.WithType("organizations"),
	)

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshal organizations list response. %w", err)
	}

	return out, nil
}

func (p *organizationsPresenter) Organizations(orgs []*models.Organization) []*hal.Resource {
	out := make([]*hal.Resource, len(orgs))

	for i, o := range orgs {
		org := domain.Organization{
			Id:        o.ID.String(),
			Name:      o.Name,
			Address:   o.Address,
			CreatedAt: uint64(o.CreatedAt.UnixMilli()),
			UpdatedAt: uint64(o.UpdatedAt.UnixMilli()),
		}

		r := hal.NewResource(org, "/organizations/"+org.Id)

		out[i] = r
	}

	return out
}
