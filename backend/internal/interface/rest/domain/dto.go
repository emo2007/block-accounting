package domain

import (
	"encoding/json"
	"fmt"
)

type JoinRequest struct {
	Mnemonic string `json:"mnemonic"`
}

type JoinResponse struct {
	Ok    bool   `json:"ok"`
	Token string `json:"token,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type LoginRequest struct {
	Mnemonc string `json:"mnemonic"`
}

type LoginResponse struct {
	Ok    bool   `json:"ok"`
	Token string `json:"token,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func BuildRequest[T any](data []byte) (*T, error) {
	var req T

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("error unmarshal request. %w", err)
	}

	return &req, nil
}
