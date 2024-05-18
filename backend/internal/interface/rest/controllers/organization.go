package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
)

type OrganizationsController interface {
	NewOrganization(w http.ResponseWriter, r *http.Request) ([]byte, error)
	ListOrganizations(w http.ResponseWriter, r *http.Request) ([]byte, error)
	// todo delete
	// todo update
}

type organizationsController struct {
	log           *slog.Logger
	orgInteractor organizations.OrganizationsInteractor
	presenter     presenters.OrganizationsPresenter
}

func NewOrganizationsController(
	log *slog.Logger,
	orgInteractor organizations.OrganizationsInteractor,
	presenter presenters.OrganizationsPresenter,
) OrganizationsController {
	return &organizationsController{
		log:           log,
		orgInteractor: orgInteractor,
		presenter:     presenter,
	}
}

func (c *organizationsController) NewOrganization(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.NewOrganizationRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	org, err := c.orgInteractor.Create(ctx, organizations.CreateParams{
		Name:           req.Name,
		Address:        req.Address,
		WalletMnemonic: req.WalletMnemonic,
	})
	if err != nil {
		return nil, fmt.Errorf("error create new organization. %w", err)
	}

	return c.presenter.ResponseCreate(org)
}

func (c *organizationsController) ListOrganizations(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.ListOrganizationsRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := c.orgInteractor.List(ctx, organizations.ListParams{
		Cursor:     req.Cursor,
		Limit:      req.Limit,
		OffsetDate: time.UnixMilli(req.OffsetDate),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch organizations list. %w", err)
	}

	return c.presenter.ResponseList(resp.Organizations, resp.NextCursor)
}
