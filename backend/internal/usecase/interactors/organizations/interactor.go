package organizations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/hdwallet"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	"github.com/google/uuid"
)

type CreateParams struct {
	Name           string
	Address        string
	WalletMnemonic string
}

type OrganizationsInteractor interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (*models.Organization, error)
}

type organizationsInteractor struct {
	log           *slog.Logger
	orgRepository organizations.Repository
}

func NewOrganizationsInteractor(
	log *slog.Logger,
	orgRepository organizations.Repository,
) OrganizationsInteractor {
	return &organizationsInteractor{
		log:           log,
		orgRepository: orgRepository,
	}
}

func (i *organizationsInteractor) Create(
	ctx context.Context,
	params CreateParams,
) (*models.Organization, error) {
	var walletSeed []byte

	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	if params.WalletMnemonic == "" {
		walletSeed = user.Seed()
	} else {
		seed, err := hdwallet.NewSeedFromMnemonic(params.WalletMnemonic)
		if err != nil {
			return nil, fmt.Errorf("error convert organization wallet mnemonic into a seed. %w", err)
		}

		walletSeed = seed
	}

	org := models.Organization{
		ID:         uuid.New(),
		Name:       params.Name,
		Address:    params.Address,
		WalletSeed: walletSeed,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := i.orgRepository.CreateAndAdd(ctx, org, user); err != nil {
		return nil, fmt.Errorf("error create new organization. %w", err)
	}

	return &org, nil
}
