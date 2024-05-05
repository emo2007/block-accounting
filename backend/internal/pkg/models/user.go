package models

import (
	"github.com/google/uuid"
)

type UserIdentity interface {
	Id() uuid.UUID
	Seed() []byte
	IsAdmin() bool
}

type User struct {
	ID        uuid.UUID
	Bip32Seed []byte
	Admin     bool
	Activated bool
}

func NewUser(id uuid.UUID, seed []byte) *User {
	return &User{
		ID:        id,
		Bip32Seed: seed,
	}
}

func (u *User) Id() uuid.UUID {
	return u.ID
}

func (u *User) Seed() []byte {
	return u.Bip32Seed
}

func (u *User) IsAdmin() bool {
	return u.Admin
}

type OrganizationUser struct {
	User
	// add org info
}
