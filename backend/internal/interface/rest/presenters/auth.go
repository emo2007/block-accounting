package presenters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
)

type AuthPresenter interface {
	CreateJoinRequest(r *http.Request) (*domain.JoinRequest, error)
	ResponseJoin(w http.ResponseWriter, user *models.User, err error) error
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

func (p *authPresenter) CreateJoinRequest(r *http.Request) (*domain.JoinRequest, error) {
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error read request body. %w", err)
	}

	var request domain.JoinRequest

	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("error unmarshal join request. %w", err)
	}

	return &request, nil
}

func (p *authPresenter) ResponseJoin(w http.ResponseWriter, user *models.User, err error) error {
	resp := new(domain.JoinResponse)

	if err != nil {
		// todo map error
	} else {
		token, err := p.jwtInteractor.NewToken(user, 24*time.Hour)
		if err != nil {
			return fmt.Errorf("error create access token. %w", err)
		}

		resp.Ok = true
		resp.Token = token
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error marshal join response. %w", err)
	}

	if _, err = w.Write(out); err != nil {
		return fmt.Errorf("error write response. %w", err)
	}

	return nil
}
