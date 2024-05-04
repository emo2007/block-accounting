package presenters

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
)

type AuthPresenter interface {
	ResponseJoin(w http.ResponseWriter, mnemonic string) error
}

type authPresenter struct{}

func NewAuthPresenter() AuthPresenter {
	return &authPresenter{}
}

func (p *authPresenter) ResponseJoin(w http.ResponseWriter, mnemonic string) error {
	out, err := json.Marshal(domain.JoinResponse{
		Mnemonic: mnemonic,
	})
	if err != nil {
		return fmt.Errorf("error marshal join response. %w", err)
	}

	if _, err = w.Write(out); err != nil {
		return fmt.Errorf("error write response. %w", err)
	}

	return nil
}
