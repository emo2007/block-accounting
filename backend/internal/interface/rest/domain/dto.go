package domain

import (
	"encoding/json"
	"fmt"
)

type JoinRequest struct {
	Mnemonic string `json:"mnemonic"`
}

func BuildRequest[T any](data []byte) (*T, error) {
	var req T

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("error unmarshal request. %w", err)
	}

	return &req, nil
}
