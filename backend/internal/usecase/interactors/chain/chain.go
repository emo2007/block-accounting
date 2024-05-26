package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/transactions"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ChainInteractor interface {
	NewMultisig(ctx context.Context, params NewMultisigParams) error
	PubKey(ctx context.Context, user *models.User) ([]byte, error)
	SalaryDeploy(ctx context.Context, firtsAdmin models.OrganizationParticipant) error
}

type chainInteractor struct {
	log          *slog.Logger
	config       config.Config
	txRepository transactions.Repository
}

func NewChainInteractor(
	log *slog.Logger,
	config config.Config,
	txRepository transactions.Repository,
) ChainInteractor {
	return &chainInteractor{
		log:          log,
		config:       config,
		txRepository: txRepository,
	}
}

type NewMultisigParams struct {
	Title         string
	Owners        []models.OrganizationParticipant
	Confirmations int
}

type newMultisigChainResponse struct {
	Address string `json:"address"`
}

func (i *chainInteractor) NewMultisig(ctx context.Context, params NewMultisigParams) error {
	endpoint := i.config.ChainAPI.Host + "/multi-sig/deploy"

	i.log.Debug(
		"deploy multisig",
		slog.String("endpoint", endpoint),
		slog.Any("params", params),
	)

	pks := make([]string, len(params.Owners))

	for i, owner := range params.Owners {
		if owner.GetUser() == nil {
			return fmt.Errorf("error invalis owners set")
		}

		pks[i] = "0x" + common.Bytes2Hex(owner.GetUser().PublicKey())
	}

	requestBody, err := json.Marshal(map[string]any{
		"owners":        pks,
		"confirmations": params.Confirmations,
	})
	if err != nil {
		return fmt.Errorf("error marshal request body. %w", err)
	}

	user, err := ctxmeta.User(ctx)
	if err != nil {
		return fmt.Errorf("error fetch user from context. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return fmt.Errorf("error fetch organization id from context. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return fmt.Errorf("error build request. %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		i.log.Error(
			"error send deploy multisig request",
			slog.String("endpoint", endpoint),
			slog.Any("params", params),
		)

		return fmt.Errorf("error build new multisig request. %w", err)
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read body. %w", err)
	}

	respObject := new(newMultisigChainResponse)

	if err := json.Unmarshal(raw, &respObject); err != nil {
		return fmt.Errorf("error parse chain-api response body. %w", err)
	}

	multisigAddress := common.Hex2Bytes(respObject.Address)

	createdAt := time.Now()

	if err := i.txRepository.AddMultisig(ctx, models.Multisig{
		ID:                    uuid.Must(uuid.NewV7()),
		Title:                 params.Title,
		Address:               multisigAddress,
		OrganizationID:        organizationID,
		Owners:                params.Owners,
		ConfirmationsRequired: params.Confirmations,
		CreatedAt:             createdAt,
		UpdatedAt:             createdAt,
	}); err != nil {
		return fmt.Errorf("error add new multisig. %w", err)
	}

	i.log.Debug(
		"deploy multisig response",
		slog.Int("code", resp.StatusCode),
	)

	return nil
}

func (i *chainInteractor) PubKey(ctx context.Context, user *models.User) ([]byte, error) {
	pubAddr := i.config.ChainAPI.Host + "/address-from-seed"

	requestBody, err := json.Marshal(map[string]any{
		"seedPhrase": user.Mnemonic,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal request body. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pubAddr, body)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetch pub address. %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read resp body. %w", err)
	}

	pubKeyStr := string(respBody)[2:]

	return common.Hex2Bytes(pubKeyStr), nil
}

func (i *chainInteractor) SalaryDeploy(ctx context.Context, firtsAdmin models.OrganizationParticipant) error {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return fmt.Errorf("error fetch user from context. %w", err)
	}

	if user.Id() != firtsAdmin.Id() || firtsAdmin.GetUser() == nil {
		return fmt.Errorf("error unauthorized access")
	}

	requestBody, err := json.Marshal(map[string]any{
		"authorizedWallet": common.Bytes2Hex(user.Seed()),
	})
	if err != nil {
		return fmt.Errorf("error marshal request body. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	endpoint := i.config.ChainAPI.Host + "/salaries/deploy"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return fmt.Errorf("error build request. %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error fetch deploy salary contract. %w", err)
	}

	defer resp.Body.Close()

	return nil
}

// func (i *chainInteractor)
