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
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/cache"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	"github.com/ethereum/go-ethereum/common"
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

	Cursor string
	Limit  uint8 // Max limit is 50 (may change)
}

type ParticipantParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID

	UsersOnly     bool
	ActiveOnly    bool
	EmployeesOnly bool
	OwnerOnly     bool
}

type ParticipantsParams struct {
	IDs            uuid.UUIDs
	OrganizationID uuid.UUID
	PKs            [][]byte

	UsersOnly     bool
	ActiveOnly    bool
	EmployeesOnly bool
	OwnerOnly     bool
}

type OrganizationsInteractor interface {
	Create(ctx context.Context, params CreateParams) (*models.Organization, error)
	List(ctx context.Context, params ListParams) (*ListResponse, error)

	Participant(ctx context.Context, params ParticipantParams) (models.OrganizationParticipant, error)
	Participants(ctx context.Context, params ParticipantsParams) ([]models.OrganizationParticipant, error)
	AddEmployee(ctx context.Context, params AddParticipantParams) (models.OrganizationParticipant, error)
	AddUser(ctx context.Context, params AddUserParams) error
}

type organizationsInteractor struct {
	log           *slog.Logger
	orgRepository organizations.Repository
	cache         cache.Cache
}

func NewOrganizationsInteractor(
	log *slog.Logger,
	orgRepository organizations.Repository,
	cache cache.Cache,
) OrganizationsInteractor {
	return &organizationsInteractor{
		log:           log,
		orgRepository: orgRepository,
		cache:         cache,
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
	Organizations models.Organizations
	NextCursor    string
}

func (i ListResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
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

	out := new(ListResponse)

	// BUG: empty org set fetched from cache
	// if err := i.cache.Get(ctx, params, out); err != nil && errors.Is(err, redis.Nil) {
	// 	i.log.Error("no cache hit!", logger.Err(err))
	// } else {
	// 	i.log.Debug("cache hit!", slog.AnyValue(out))
	// 	return out, nil
	// }

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
		UserId:   params.UserId,
		Ids:      params.Ids,
		Limit:    int64(params.Limit),
		CursorId: cursor.Id,
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

	out = &ListResponse{
		Organizations: orgs,
		NextCursor:    nextCursor,
	}

	if err = i.cache.Cache(ctx, params, *out, time.Hour*1); err != nil {
		i.log.Error("error add cache record", logger.Err(err))
	}

	return out, nil
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

func (i *organizationsInteractor) Participant(
	ctx context.Context,
	params ParticipantParams,
) (models.OrganizationParticipant, error) {
	participants, err := i.Participants(ctx, ParticipantsParams{
		IDs:            uuid.UUIDs{params.ID},
		OrganizationID: params.OrganizationID,
		ActiveOnly:     params.ActiveOnly,
		UsersOnly:      params.UsersOnly,
		EmployeesOnly:  params.EmployeesOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch organization participant. %w", err)
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("error organization participant empty. %w", err)
	}

	return participants[0], nil
}

func (i *organizationsInteractor) Participants(
	ctx context.Context,
	params ParticipantsParams,
) ([]models.OrganizationParticipant, error) {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	_, err = i.orgRepository.Participants(ctx, organizations.ParticipantsParams{
		Ids:            uuid.UUIDs{user.Id()},
		OrganizationId: params.OrganizationID,
		ActiveOnly:     params.ActiveOnly,
		UsersOnly:      true,
	})
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("error fetch organization user. %w", err),
			ErrorUnauthorizedAccess,
		)
	}

	participants, err := i.orgRepository.Participants(ctx, organizations.ParticipantsParams{
		Ids:            params.IDs,
		OrganizationId: params.OrganizationID,
		PKs:            params.PKs,
		UsersOnly:      params.UsersOnly,
		EmployeesOnly:  params.EmployeesOnly,
		ActiveOnly:     params.ActiveOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch organization participants. %w", err)
	}

	return participants, nil
}

type AddParticipantParams struct {
	OrganizationID uuid.UUID
	EmployeeUserID uuid.UUID
	Name           string
	Position       string
	WalletAddress  string
}

func (i *organizationsInteractor) AddEmployee(
	ctx context.Context,
	params AddParticipantParams,
) (models.OrganizationParticipant, error) {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetch user from context. %w", err)
	}

	actor, err := i.Participant(ctx, ParticipantParams{
		ID:             user.Id(),
		OrganizationID: params.OrganizationID,
		ActiveOnly:     true,
		UsersOnly:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch actor. %w", err)
	}

	if !actor.IsAdmin() || !actor.IsOwner() {
		return nil, fmt.Errorf("error actor not an owner")
	}

	if !common.IsHexAddress(params.WalletAddress) {
		return nil, fmt.Errorf("error invalid address")
	}

	participantID := uuid.Must(uuid.NewV7())

	empl := models.Employee{
		ID:             participantID,
		EmployeeName:   params.Name,
		UserID:         params.EmployeeUserID,
		OrganizationId: params.OrganizationID,
		WalletAddress:  common.FromHex(params.WalletAddress),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err = i.orgRepository.AddEmployee(ctx, empl); err != nil {
		return nil, fmt.Errorf("error add new employee. %w", err)
	}

	return &empl, nil
}

type AddUserParams struct {
	User           *models.User
	IsAdmin        bool
	IsOwner        bool
	OrganizationID uuid.UUID
	SkipRights     bool
}

func (i *organizationsInteractor) AddUser(ctx context.Context, params AddUserParams) error {
	if !params.SkipRights {
		user, err := ctxmeta.User(ctx)
		if err != nil {
			return fmt.Errorf("error fetch user from context. %w", err)
		}

		actor, err := i.Participant(ctx, ParticipantParams{
			ID:             user.Id(),
			OrganizationID: params.OrganizationID,
			ActiveOnly:     true,
			UsersOnly:      true,
		})
		if err != nil {
			return fmt.Errorf("error fetch actor. %w", err)
		}

		if !actor.IsAdmin() || !actor.IsOwner() {
			return fmt.Errorf("error actor not an owner")
		}
	}

	i.log.Debug(
		"add user",
		slog.Any("params", params),
	)

	if err := i.orgRepository.AddParticipant(ctx, organizations.AddParticipantParams{
		OrganizationId: params.OrganizationID,
		UserId:         params.User.Id(),
		IsAdmin:        params.IsAdmin,
		IsOwner:        params.IsOwner,
	}); err != nil {
		return fmt.Errorf("error add user into organization. %w", err)
	}

	return nil
}
