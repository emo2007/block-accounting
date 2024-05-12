package organizations

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/hdwallet"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	"github.com/google/uuid"
)

var (
	ErrorUnauthorizedAccess = errors.New("unauthorized access")
)

type CreateParams struct {
	Name           string
	Address        string
	WalletMnemonic string
}

type ListParams struct {
	Ids    uuid.UUIDs
	UserId uuid.UUID

	Cursor     string
	OffsetDate time.Time
	Limit      uint8 // Max limit is 50 (may change)
}

type OrganizationsInteractor interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (*models.Organization, error)
	List(
		ctx context.Context,
		params ListParams,
	) (*ListResponse, error)
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

type organizationsListCursor struct {
	Id uuid.UUID `json:"id"`
}

func newOrganizationsListCursor(id ...uuid.UUID) *organizationsListCursor {
	if len(id) > 0 {
		return &organizationsListCursor{id[0]}
	}

	return new(organizationsListCursor)
}

func (c *organizationsListCursor) encode() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal reaction id. %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil

}

func (c *organizationsListCursor) decode(s string) error {
	if c == nil {
		return nil
	}

	token, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("error decode token. %w", err)
	}

	return json.Unmarshal(token, c)
}

type ListResponse struct {
	Organizations []*models.Organization
	NextCursor    string
}

func (i *organizationsInteractor) List(
	ctx context.Context,
	params ListParams,
) (*ListResponse, error) {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	if params.UserId != uuid.Nil {
		if params.UserId != user.Id() {
			return nil, fmt.Errorf("error unauthorized organizations list access. %w", ErrorUnauthorizedAccess)
		}
	} else {
		params.UserId = user.Id()
	}

	if params.Limit <= 0 || params.Limit > 50 {
		params.Limit = 50
	}

	cursor := newOrganizationsListCursor()

	if params.Cursor != "" {
		if err := cursor.decode(params.Cursor); err != nil {
			return nil, fmt.Errorf("error decode cursor value. %w", err) // maybe just log error?
		}
	}

	i.log.Debug(
		"organizations_list",
		slog.String("cursor", params.Cursor),
		slog.Int("limit", int(params.Limit)),
		slog.Any("cursor-id", cursor.Id),
		slog.Any("ids", params.Ids),
		slog.Any("user_id", params.UserId),
	)

	orgs, err := i.orgRepository.Get(ctx, organizations.GetParams{
		UserId:     params.UserId,
		Ids:        params.Ids,
		OffsetDate: params.OffsetDate,
		Limit:      int64(params.Limit),
		CursorId:   cursor.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch organizations. %w", err)
	}

	var nextCursor string

	if len(orgs) >= 50 || len(orgs) >= int(params.Limit) {
		cursor.Id = orgs[len(orgs)-1].ID
		if nextCursor, err = cursor.encode(); err != nil {
			return nil, fmt.Errorf("error encode next page token. %w", err) // maybe just log error?
		}
	}

	return &ListResponse{
		Organizations: orgs,
		NextCursor:    nextCursor,
	}, nil
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
		ID:         uuid.Must(uuid.NewV7()),
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
