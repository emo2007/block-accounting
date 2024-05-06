package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/bip32"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
)

var (
	ErrorAuthInvalidMnemonic = errors.New("Invalid Mnemonic")
)

type AuthController interface {
	Join(w http.ResponseWriter, req *http.Request) error
	JoinWithInvite(w http.ResponseWriter, req *http.Request) error
	Login(w http.ResponseWriter, req *http.Request) error
	Invite(w http.ResponseWriter, req *http.Request) error
}

type authController struct {
	log             *slog.Logger
	presenter       presenters.AuthPresenter
	usersInteractor users.UsersInteractor
}

func NewAuthController(
	log *slog.Logger,
	presenter presenters.AuthPresenter,
	usersInteractor users.UsersInteractor,
) AuthController {
	return &authController{
		log:             log,
		presenter:       presenter,
		usersInteractor: usersInteractor,
	}
}

func (c *authController) Join(w http.ResponseWriter, req *http.Request) error {
	request, err := c.presenter.CreateJoinRequest(req)
	if err != nil {
		return fmt.Errorf("error create join request. %w", err)
	}

	c.log.Debug("join request", slog.String("mnemonic", request.Mnemonic))

	if !bip32.IsMnemonicValid(request.Mnemonic) {
		return fmt.Errorf("error invalid mnemonic. %w", ErrorAuthInvalidMnemonic)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if _, err = c.usersInteractor.Create(ctx, users.CreateParams{
		Mnemonic: request.Mnemonic,
		IsAdmin:  true,
		Activate: true,
	}); err != nil {
		return fmt.Errorf("error create new user. %w", err)
	}

	return nil
}

func (c *authController) JoinWithInvite(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}

func (c *authController) Login(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}

const mnemonicEntropyBitSize int = 256

func (c *authController) Invite(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}
