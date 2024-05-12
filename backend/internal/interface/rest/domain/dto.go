package domain

import (
	"encoding/json"
	"fmt"
)

// Generic

type Collection[T any] struct {
	Items      []T        `json:"items,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	NextCursor string `json:"next_cursor,omitempty"`
	TotalItems uint32 `json:"total_items,omitempty"`
}

// Auth related DTO's

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

// Organizations

type NewOrganizationRequest struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	WalletMnemonic string `json:"wallet_mnemonic,omitempty"`
}

type NewOrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type ListOrganizationsRequest struct {
	Cursor     string `json:"cursor,omitempty"`
	Limit      uint8  `json:"limit,omitempty"`       // Default: 50, Max: 50
	OffsetDate int64  `json:"offset_date,omitempty"` // List organizations, updated since the date
}

type ListOrganizationsResponse struct {
	Collection[Organization]
}

func BuildRequest[T any](data []byte) (*T, error) {
	var req T

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("error unmarshal request. %w", err)
	}

	return &req, nil
}
