package presenters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
)

type AuthPresenter interface {
	CreateJoinRequest(r *http.Request) (*domain.JoinRequest, error)
	// ResponseJoin(w http.ResponseWriter, mnemonic string) error
}

type authPresenter struct{}

func NewAuthPresenter() AuthPresenter {
	return &authPresenter{}
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

// func (p *authPresenter) ResponseJoin(w http.ResponseWriter, mnemonic string) error {
// 	out, err := json.Marshal(domain.JoinResponse{
// 		Mnemonic: mnemonic,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error marshal join response. %w", err)
// 	}

// 	if _, err = w.Write(out); err != nil {
// 		return fmt.Errorf("error write response. %w", err)
// 	}

// 	return nil
// }
