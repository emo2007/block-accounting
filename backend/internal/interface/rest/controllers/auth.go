package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/bip32"
)

var (
	ErrorAuthInvalidMnemonic = errors.New("Invalid Mnemonic")
)

type AuthController interface {
	Join(w http.ResponseWriter, req *http.Request) error
	Login(w http.ResponseWriter, req *http.Request) error
	Invite(w http.ResponseWriter, req *http.Request) error
}

type authController struct {
	log       *slog.Logger
	presenter presenters.AuthPresenter
	// interactors ...
}

func NewAuthController(
	log *slog.Logger,
	presenter presenters.AuthPresenter,
) AuthController {
	return &authController{
		log:       log,
		presenter: presenter,
	}
}

const mnemonicEntropyBitSize int = 256

func (c *authController) Join(w http.ResponseWriter, req *http.Request) error {
	request, err := c.presenter.CreateJoinRequest(req)
	if err != nil {
		return fmt.Errorf("error create join request. %w", err)
	}

	c.log.Debug("join request", slog.String("mnemonic", request.Mnemonic))

	if !bip32.IsMnemonicValid(request.Mnemonic) {
		return fmt.Errorf("error invalid mnemonic. %w", ErrorAuthInvalidMnemonic)
	}

	// todo create user

	return nil
}

func (c *authController) Login(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}

func (c *authController) Invite(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}
