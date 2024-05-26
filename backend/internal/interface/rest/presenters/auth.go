package presenters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/google/uuid"
)

type AuthPresenter interface {
	ResponseJoin(user *models.User) ([]byte, error)
	ResponseLogin(user *models.User) ([]byte, error)
	ResponseRefresh(tokens jwt.AccessToken) ([]byte, error)
	ResponseNewInvite(
		ctx context.Context,
		organizationID uuid.UUID,
		link string,
	) ([]byte, error)
}

type authPresenter struct {
	jwtInteractor jwt.JWTInteractor
}

func NewAuthPresenter(
	jwtInteractor jwt.JWTInteractor,
) AuthPresenter {
	return &authPresenter{
		jwtInteractor: jwtInteractor,
	}
}

func (p *authPresenter) ResponseJoin(user *models.User) ([]byte, error) {
	tokens, err := p.jwtInteractor.NewToken(user, 24*time.Hour, "")
	if err != nil {
		return nil, fmt.Errorf("error create access token. %w", err)
	}

	out, err := json.Marshal(domain.LoginResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
		ExpiredAt:    tokens.ExpiredAt.UnixMilli(),
		RTExpiredAt:  tokens.RTExpiredAt.UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal join response. %w", err)
	}

	return out, nil
}

func (p *authPresenter) ResponseLogin(user *models.User) ([]byte, error) {
	tokens, err := p.jwtInteractor.NewToken(user, 24*time.Hour, "")
	if err != nil {
		return nil, fmt.Errorf("error create access token. %w", err)
	}

	out, err := json.Marshal(domain.LoginResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
		ExpiredAt:    tokens.ExpiredAt.UnixMilli(),
		RTExpiredAt:  tokens.RTExpiredAt.UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal login response. %w", err)
	}

	return out, nil
}

func (p *authPresenter) ResponseRefresh(tokens jwt.AccessToken) ([]byte, error) {
	out, err := json.Marshal(domain.LoginResponse{
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
		ExpiredAt:    tokens.ExpiredAt.UnixMilli(),
		RTExpiredAt:  tokens.RTExpiredAt.UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal refresh response. %w", err)
	}

	return out, nil
}

func (p *authPresenter) ResponseNewInvite(
	ctx context.Context,
	organizationID uuid.UUID,
	link string,
) ([]byte, error) {
	out, err := json.Marshal(map[string]string{
		"link": "/" + organizationID.String() + "/invite/" + link,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal refresh response. %w", err)
	}

	return out, nil
}
