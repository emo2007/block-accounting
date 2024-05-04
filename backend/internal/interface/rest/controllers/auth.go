package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/presenters"
	"github.com/emochka2007/block-accounting/internal/pkg/bip32"
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
	entropy, err := bip32.NewEntropy(mnemonicEntropyBitSize)
	if err != nil {
		return fmt.Errorf("error generate new entropy. %w", err)
	}

	mnemonic, err := bip32.NewMnemonic(entropy)
	if err != nil {
		return fmt.Errorf("error generate mnemonic from entropy. %w", err)
	}

	// todo create user

	return c.presenter.ResponseJoin(w, mnemonic)
}

func (c *authController) Login(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}

func (c *authController) Invite(w http.ResponseWriter, req *http.Request) error {
	return nil // implement
}
