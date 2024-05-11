package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
)

type OrganizationsController interface {
	NewOrganization(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type organizationsController struct {
	log           *slog.Logger
	orgInteractor organizations.OrganizationsInteractor
}

func NewOrganizationsController(
	log *slog.Logger,
	orgInteractor organizations.OrganizationsInteractor,
) OrganizationsController {
	return &organizationsController{
		log:           log,
		orgInteractor: orgInteractor,
	}
}

func (c *organizationsController) NewOrganization(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	_, err := presenters.CreateRequest[domain.NewOrganizationRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	// todo call int.Create

	// todo build response

	return nil, nil
}
