package models

import (
	"time"

	"github.com/google/uuid"
)

type Multisig struct {
	ID                    uuid.UUID
	Title                 string
	Address               []byte
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

type Payroll struct {
	ID             uuid.UUID
	Title          string
	Address        []byte
	OrganizationID uuid.UUID
	MultisigID     uuid.UUID
}
