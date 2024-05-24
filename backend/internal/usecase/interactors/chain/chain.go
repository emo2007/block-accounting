package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/transactions"
	"github.com/ethereum/go-ethereum/common"
)

type ChainInteractor interface {
	NewMultisig(ctx context.Context, params NewMultisigParams) error
}

type chainInteractor struct {
	log             *slog.Logger
	config          config.Config
	txRepository    transactions.Repository
	usersInteractor users.UsersInteractor
}

func NewChainInteractor(
	log *slog.Logger,
	config config.Config,
	txRepository transactions.Repository,
	usersInteractor users.UsersInteractor,
) ChainInteractor {
	return &chainInteractor{
		log:             log,
		config:          config,
		txRepository:    txRepository,
		usersInteractor: usersInteractor,
	}
}

type NewMultisigParams struct {
	OwnersPKs     []string
	Confirmations int
}

func (i *chainInteractor) NewMultisig(ctx context.Context, params NewMultisigParams) error {
	deployAddr := i.config.ChainAPI.Host + "/multi-sig/deploy"

	i.log.Debug(
		"deploy multisig",
		slog.String("endpoint", deployAddr),
		slog.Any("params", params),
	)

	requestBody, err := json.Marshal(map[string]any{
		"owners":        params.OwnersPKs,
		"confirmations": params.Confirmations,
	})
	if err != nil {
		return fmt.Errorf("error marshal request body. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	doneCh := make(chan struct{})

	errCh := make(chan error)

	go func() {
		resp, err := http.Post(http.MethodPost, deployAddr, body)
		if err != nil {
			i.log.Error(
				"error send deploy multisig request",
				slog.String("endpoint", deployAddr),
				slog.Any("params", params),
			)

			errCh <- fmt.Errorf("error build new multisig request. %w", err)
			return
		}

		defer resp.Body.Close()

		i.log.Debug(
			"deploy multisig response",
			slog.Int("code", resp.StatusCode),
		)

		if _, ok := <-doneCh; ok {
			doneCh <- struct{}{}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (i *chainInteractor) PubKey(ctx context.Context, user *models.User) ([]byte, error) {
	pubAddr := i.config.ChainAPI.Host + "/address-from-seed"

	doneCh := make(chan struct{})
	errCh := make(chan error)

	requestBody, err := json.Marshal(map[string]any{
		"seedPhrase": user.Mnemonic,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal request body. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	var pubKeyStr string

	go func() {
		resp, err := http.Post(pubAddr, "application/json", body)
		if err != nil {
			errCh <- fmt.Errorf("error fetch pub address. %w", err)
			return
		}

		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			errCh <- fmt.Errorf("error read resp body. %w", err)
			return
		}

		pubKeyStr = string(respBody)

		doneCh <- struct{}{}
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-doneCh:
		return common.Hex2Bytes(pubKeyStr), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
