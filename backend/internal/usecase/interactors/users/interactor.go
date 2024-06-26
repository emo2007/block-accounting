package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/hdwallet"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/chain"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/users"
	"github.com/google/uuid"
)

var (
	ErrorUsersNotFound = errors.New("users not found")
)

type CreateParams struct {
	Name     string
	Email    string
	Phone    string
	Tg       string
	Mnemonic string
	Activate bool
	Owner    bool
	Admin    bool
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
	log             *slog.Logger
	usersRepo       users.Repository
	chainInteractor chain.ChainInteractor
}

func NewUsersInteractor(
	log *slog.Logger,
	usersRepo users.Repository,
	chainInteractor chain.ChainInteractor,
) UsersInteractor {
	return &usersInteractor{
		log:             log,
		usersRepo:       usersRepo,
		chainInteractor: chainInteractor,
	}
}

func (i *usersInteractor) Create(ctx context.Context, params CreateParams) (*models.User, error) {
	seed, err := hdwallet.NewSeedFromMnemonic(params.Mnemonic)
	if err != nil {
		return nil, fmt.Errorf("error convert mnemonic into a seed. %w", err)
	}

	user := models.NewUser(
		uuid.Must(uuid.NewV7()),
		seed,
		params.Activate,
		time.Now(),
	)

	user.Name = params.Name
	user.Mnemonic = params.Mnemonic

	user.Credentails = &models.UserCredentials{
		Email:    params.Email,
		Phone:    params.Phone,
		Telegram: params.Tg,
	}

	pk, err := i.chainInteractor.PubKey(ctx, user)
	if err != nil {
		// todo пока мокнуть
		return nil, fmt.Errorf("error fetch user pub key. %w", err)
	}

	user.PK = pk

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
