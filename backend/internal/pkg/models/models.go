package models

import (
	"time"

	"github.com/google/uuid"
)

type Multisig struct {
	ID                    uuid.UUID
	Title                 string
	OrganizationID        uuid.UUID
	Owners                []OrganizationParticipant
	ConfirmationsRequired int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type MultisigConfirmation struct {
	MultisigID uuid.UUID
	Owner      OrganizationParticipant
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
