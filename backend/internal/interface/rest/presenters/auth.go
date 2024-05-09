package presenters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
)

type AuthPresenter interface {
	ResponseJoin(w http.ResponseWriter, user *models.User) error
	ResponseLogin(w http.ResponseWriter, user *models.User) error
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

func (p *authPresenter) ResponseJoin(w http.ResponseWriter, user *models.User) error {
	resp := new(domain.JoinResponse)

	token, err := p.jwtInteractor.NewToken(user, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("error create access token. %w", err)
	}

	resp.Token = token

	out, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error marshal join response. %w", err)
	}

	if _, err = w.Write(out); err != nil {
		return fmt.Errorf("error write response. %w", err)
	}

	return nil
}

func (p *authPresenter) ResponseLogin(w http.ResponseWriter, user *models.User) error {
	resp := new(domain.LoginResponse)

	token, err := p.jwtInteractor.NewToken(user, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("error create access token. %w", err)
	}

	resp.Token = token

	out, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error marshal login response. %w", err)
	}

	if _, err = w.Write(out); err != nil {
		return fmt.Errorf("error write response. %w", err)
	}

	return nil
}
