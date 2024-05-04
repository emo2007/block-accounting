package users

import (
	"context"

	"github.com/emochka2007/block-accounting/internal/pkg/models"
)

type GetParams struct {
	Id             string
	OrganizationId string
	Seed           []byte
}

// todo implement
type Repository interface {
	Get(ctx context.Context, params GetParams) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Activate(ctx context.Context, id string) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}
