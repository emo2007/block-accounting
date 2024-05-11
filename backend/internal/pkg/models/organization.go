package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID         uuid.UUID
	Name       string
	Address    string
	WalletSeed []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
