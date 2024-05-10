package models

import (
	"time"

	"github.com/google/uuid"
)

type UserIdentity interface {
	Id() uuid.UUID
	Seed() []byte
}

type User struct {
	ID uuid.UUID

	Name string

	Credentails UserCredentials

	Bip39Seed []byte
	Activated bool
	CreatedAt time.Time
}

type UserCredentials struct {
	Email    string
	Phone    string
	Telegram string
}

func NewUser(
	id uuid.UUID,
	seed []byte,
	activated bool,
	createdAt time.Time,
) *User {
	return &User{
		ID:        id,
		Bip39Seed: seed,
		Activated: activated,
		CreatedAt: createdAt,
	}
}

func (u *User) Id() uuid.UUID {
	return u.ID
}

func (u *User) Seed() []byte {
	return u.Bip39Seed
}

type OrganizationParticipant interface {
	UserIdentity

	IsAdmin() bool
	Position() string
}

type OrganizationUser struct {
	User

	OrgPosition string
	Admin       bool
}

func (u *OrganizationUser) IsAdmin() bool {
	return u.Admin
}

func (u *OrganizationUser) Position() string {
	return u.OrgPosition
}
