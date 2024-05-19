package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id uuid.UUID

	Description    string
	OrganizationId uuid.UUID
	CreatedBy      *OrganizationUser
	Amount         float64

	ToAddr []byte

	MaxFeeAllowed float64
	Deadline      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	ConfirmedAt time.Time
	CancelledAt time.Time

	CommitedAt time.Time
}

type TransactionConfirmation struct {
	TxId           uuid.UUID
	User           *User
	OrganizationId uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Confirmed      bool
}
