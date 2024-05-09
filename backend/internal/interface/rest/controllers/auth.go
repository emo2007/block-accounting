package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/interface/rest/domain"
	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/bip32"
	"github.com/emochka2007/block-accounting/internal/pkg/hdwallet"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
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
	jwtInteractor   jwt.JWTInteractor
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
	request, err := presenters.CreateRequest[domain.JoinRequest](req)
	if err != nil {
		return fmt.Errorf("error create join request. %w", err)
	}

	c.log.Debug("join request", slog.String("mnemonic", request.Mnemonic))

	if !bip32.IsMnemonicValid(request.Mnemonic) {
		return fmt.Errorf("error invalid mnemonic. %w", ErrorAuthInvalidMnemonic)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user, err := c.usersInteractor.Create(ctx, users.CreateParams{
		Mnemonic: request.Mnemonic,
		IsAdmin:  true,
		Activate: true,
	})
	if err != nil {
		return fmt.Errorf("error create new user. %w", err)
	}

	c.log.Debug("join request", slog.String("user id", user.ID.String()))

	return c.presenter.ResponseJoin(w, user)
}

// NIT: wrap with idempotent action handler
func (c *authController) Login(w http.ResponseWriter, req *http.Request) error {
	request, err := presenters.CreateRequest[domain.LoginRequest](req)
	if err != nil {
		return fmt.Errorf("error create login request. %w", err)
	}

	c.log.Debug("login request", slog.String("mnemonic", request.Mnemonic))

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	seed, err := hdwallet.NewSeedFromMnemonic(request.Mnemonic)
	if err != nil {
		return fmt.Errorf("error create seed from mnemonic. %w", err)
	}

	users, err := c.usersInteractor.Get(ctx, users.GetParams{
		Seed: seed,
	})
	if err != nil {
		return fmt.Errorf("error fetch user by seed. %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("error empty users set")
	}

	c.log.Debug("login request", slog.String("user id", users[0].ID.String()))

	return c.presenter.ResponseLogin(w, users[0])
}

// const mnemonicEntropyBitSize int = 256

func (c *authController) Invite(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (c *authController) JoinWithInvite(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}
