package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/hdwallet"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/users"
	"github.com/google/uuid"
)

var (
	ErrorUsersNotFound = errors.New("users not found")
)

type CreateParams struct {
	Mnemonic string
	IsAdmin  bool
	Activate bool
}

type GetParams struct {
	Ids            uuid.UUIDs
	OrganizationId uuid.UUID
	Mnemonic       string
	Seed           []byte
}

type DeleteParams struct {
	Id             uuid.UUID
	OrganizationId uuid.UUID
}

type ActivateParams struct {
	Id             uuid.UUID
	OrganizationId uuid.UUID
}

type UsersInteractor interface {
	Create(ctx context.Context, params CreateParams) (*models.User, error)
	Update(ctx context.Context, newState models.User) error
	Activate(ctx context.Context, params ActivateParams) error
	Get(ctx context.Context, params GetParams) ([]*models.User, error)
	Delete(ctx context.Context, params DeleteParams) error
}

type usersInteractor struct {
	log       *slog.Logger
	usersRepo users.Repository
}

func NewUsersInteractor(
	log *slog.Logger,
	usersRepo users.Repository,
) UsersInteractor {
	return &usersInteractor{
		log:       log,
		usersRepo: usersRepo,
	}
}

func (i *usersInteractor) Create(ctx context.Context, params CreateParams) (*models.User, error) {
	seed, err := hdwallet.NewSeedFromMnemonic(params.Mnemonic)
	if err != nil {
		return nil, fmt.Errorf("error convert mnemonic into a bytes. %w", err)
	}

	user := models.NewUser(
		uuid.New(),
		seed,
		params.IsAdmin,
		params.Activate,
		time.Now(),
	)

	if err = i.usersRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error create new user. %w", err)
	}

	return user, nil
}

func (i *usersInteractor) Update(ctx context.Context, newState models.User) error {
	return nil
}

func (i *usersInteractor) Activate(ctx context.Context, params ActivateParams) error {
	return nil
}

func (i *usersInteractor) Get(ctx context.Context, params GetParams) ([]*models.User, error) {
	users, err := i.usersRepo.Get(ctx, users.GetParams{
		Ids:            params.Ids,
		OrganizationId: params.OrganizationId,
		Seed:           params.Seed,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch users from repository. %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("error empty users set. %w", ErrorUsersNotFound)
	}

	return users, nil
}

func (i *usersInteractor) Delete(ctx context.Context, params DeleteParams) error {
	return nil
}
