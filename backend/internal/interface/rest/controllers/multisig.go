package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/chain"
)

type MultisigController interface {
}

type multisigController struct {
	log             *slog.Logger
	chainInteractor chain.ChainInteractor
}

func (c *multisigController) New(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := presenters.CreateRequest[domain.NewMultisigRequest](r)
	if err != nil {
		return nil, fmt.Errorf("error build new multisig request. %w", err)
	}

	c.log.Debug("new_multisig", slog.Any("request", req))

	panic("implement me!")
}
