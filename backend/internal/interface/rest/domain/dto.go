package domain

import (
	"encoding/json"
	"fmt"
)

type JoinRequest struct {
	Name       string `json:"name,omitempty"`
	Credentals struct {
		Email    string `json:"email,omitempty"`
		Phone    string `json:"phone,omitempty"`
		Telegram string `json:"telegram,omitempty"`
	} `json:"credentals,omitempty"`

	Mnemonic string `json:"mnemonic"`
}

type JoinResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Mnemonic string `json:"mnemonic"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type NewOrganizationRequest struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	WalletMnemonic string `json:"wallet_mnemonic,omitempty"`
}

func BuildRequest[T any](data []byte) (*T, error) {
	var req T

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("error unmarshal request. %w", err)
	}

	return &req, nil
}
