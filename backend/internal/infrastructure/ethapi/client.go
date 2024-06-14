package ethapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/ethereum/go-ethereum/common"
)

type EthAPIClient interface {
	DeployMultisig(
		ctx context.Context,
		ownersPubKeys []string,
		confirmations int,
	) (string, error)
}

type ethAPIClient struct {
	host string
	c    *http.Client
	log  *slog.Logger
}

func NewEthAPIClient(
	address string,
	log *slog.Logger,
) EthAPIClient {
	return &ethAPIClient{
		host: address,
		c:    &http.Client{},
		log:  log,
	}
}

func (c *ethAPIClient) DeployMultisig(
	ctx context.Context,
	ownersPubKeys []string,
	confirmations int,
) (string, error) {
	endpoint := c.host + "/multi-sig/deploy"

	user, err := ctxmeta.User(ctx)
	if err != nil {
		return "", fmt.Errorf("error fetch user from context. %w", err)
	}

	requestBody, err := json.Marshal(map[string]any{
		"owners":        ownersPubKeys,
		"confirmations": confirmations,
	})
	if err != nil {
		return "", fmt.Errorf("error marshal request body. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return "", fmt.Errorf("error create a new request. %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error send deploy multisig request. %w", err)
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error read response body. %w", err)
	}

	c.log.Debug(
		"deploy multisig response",
		slog.Int("code", resp.StatusCode),
		slog.String("body", string(raw)),
	)

	respBody := make(map[string]string, 1)

	if err := json.Unmarshal(raw, &respBody); err != nil {
		return "", fmt.Errorf("error parse chain-api response body. %w", err)
	}

	if address, ok := respBody["address"]; ok {
		if address == "" {
			return "", fmt.Errorf("error empty address")
		}

		return address, nil
	}

	return "", fmt.Errorf("error deploy multisig. %+v", respBody)
}
