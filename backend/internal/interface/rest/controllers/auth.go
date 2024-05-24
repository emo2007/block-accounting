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
	"github.com/emochka2007/block-accounting/internal/pkg/bip39"
	"github.com/emochka2007/block-accounting/internal/pkg/hdwallet"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
)

var (
	ErrorAuthInvalidMnemonic = errors.New("invalid mnemonic")
	ErrorTokenRequired       = errors.New("token required")
)

type AuthController interface {
	Join(w http.ResponseWriter, req *http.Request) ([]byte, error)
	JoinWithInvite(w http.ResponseWriter, req *http.Request) ([]byte, error)
	Login(w http.ResponseWriter, req *http.Request) ([]byte, error)
	Invite(w http.ResponseWriter, req *http.Request) ([]byte, error)
	Refresh(w http.ResponseWriter, req *http.Request) ([]byte, error)
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
	jwtInteractor jwt.JWTInteractor,
) AuthController {
	return &authController{
		log:             log,
		presenter:       presenter,
		usersInteractor: usersInteractor,
		jwtInteractor:   jwtInteractor,
	}
}

func (c *authController) Join(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	request, err := presenters.CreateRequest[domain.JoinRequest](req)
	if err != nil {
		return nil, fmt.Errorf("error create join request. %w", err)
	}

	c.log.Debug("join request", slog.String("mnemonic", request.Mnemonic))

	if !bip39.IsMnemonicValid(request.Mnemonic) {
		return nil, fmt.Errorf("error invalid mnemonic. %w", ErrorAuthInvalidMnemonic)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	user, err := c.usersInteractor.Create(ctx, users.CreateParams{
		Name:     request.Name,
		Email:    request.Credentals.Email,
		Phone:    request.Credentals.Phone,
		Tg:       request.Credentals.Telegram,
		Mnemonic: request.Mnemonic,
		Activate: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error create new user. %w", err)
	}

	c.log.Debug("join request", slog.String("user id", user.ID.String()))

	return c.presenter.ResponseJoin(user)
}

// NIT: wrap with idempotent action handler
func (c *authController) Login(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	request, err := presenters.CreateRequest[domain.LoginRequest](req)
	if err != nil {
		return nil, fmt.Errorf("error create login request. %w", err)
	}

	c.log.Debug("login request", slog.String("mnemonic", request.Mnemonic))

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	seed, err := hdwallet.NewSeedFromMnemonic(request.Mnemonic)
	if err != nil {
		return nil, fmt.Errorf("error create seed from mnemonic. %w", err)
	}

	users, err := c.usersInteractor.Get(ctx, users.GetParams{
		Seed: seed,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch user by seed. %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("error empty users set")
	}

	c.log.Debug("login request", slog.String("user id", users[0].ID.String()))

	return c.presenter.ResponseLogin(users[0])
}

func (c *authController) Refresh(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	request, err := presenters.CreateRequest[domain.RefreshRequest](req)
	if err != nil {
		return nil, fmt.Errorf("error create refresh request. %w", err)
	}

	c.log.Debug(
		"refresh request",
		slog.String("token", request.Token),
		slog.String("refresh_token", request.RefreshToken),
	)

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	newTokens, err := c.jwtInteractor.RefreshToken(ctx, request.Token, request.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error refresh access token. %w", err)
	}

	return c.presenter.ResponseRefresh(newTokens)
}

// const mnemonicEntropyBitSize int = 256

func (c *authController) Invite(w http.ResponseWriter, req *http.Request) ([]byte, error) {

	return nil, nil
}

func (c *authController) JoinWithInvite(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	return nil, nil // implement
}
