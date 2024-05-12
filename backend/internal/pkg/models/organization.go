package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"addess"`
	WalletSeed []byte    `json:"wallet_seed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (i Organization) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
