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
	UpdatedAt time.Time
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

type OrganizationParticipantType int

const (
	OrganizationParticipantTypeUser OrganizationParticipantType = iota
	OrganizationParticipantTypeEmployee
)

type OrganizationParticipant interface {
	Id() uuid.UUID

	Type() OrganizationParticipantType

	GetUser() *OrganizationUser
	GetEmployee() *Employee

	IsAdmin() bool
	Position() string
}

type OrganizationUser struct {
	User

	OrgPosition string
	Admin       bool

	Employee *Employee
}

func (u *OrganizationUser) Type() OrganizationParticipantType {
	return OrganizationParticipantTypeUser
}

func (u *OrganizationUser) GetUser() *OrganizationUser {
	return u
}

func (u *OrganizationUser) GetEmployee() *Employee {
	return u.Employee
}

func (u *OrganizationUser) IsAdmin() bool {
	return u.Admin
}

func (u *OrganizationUser) Position() string {
	return u.OrgPosition
}

type Employee struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationId uuid.UUID
	WalletAddress  []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (u *Employee) Id() uuid.UUID {
	return u.ID
}

func (u *Employee) UserId() uuid.UUID {
	return u.UserID
}

func (u *Employee) Type() OrganizationParticipantType {
	return OrganizationParticipantTypeEmployee
}

func (u *Employee) GetUser() *OrganizationUser {
	return nil
}

func (u *Employee) GetEmployee() *Employee {
	return u
}

func (u *Employee) IsAdmin() bool {
	return false
}

func (u *Employee) Position() string {
	return "" // todo
}
