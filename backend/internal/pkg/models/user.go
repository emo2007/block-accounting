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

	Credentails *UserCredentials

	PK        []byte
	Bip39Seed []byte
	Mnemonic  string
	Activated bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
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
	ParticipantName() string
	Type() OrganizationParticipantType

	GetUser() *OrganizationUser
	GetEmployee() *Employee

	IsAdmin() bool
	IsOwner() bool
	Position() string
	IsActive() bool

	CreatedDate() time.Time
	UpdatedDate() time.Time
	DeletedDate() time.Time
}

type OrganizationUser struct {
	User

	OrgPosition string
	Admin       bool
	Owner       bool

	Employee *Employee

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (u *OrganizationUser) ParticipantName() string {
	return u.Name
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

func (u *OrganizationUser) IsOwner() bool {
	return u.Owner
}

func (u *OrganizationUser) Position() string {
	return u.OrgPosition
}

func (u *OrganizationUser) IsActive() bool {
	return u.Activated
}

func (u *OrganizationUser) CreatedDate() time.Time {
	return u.CreatedAt
}

func (u *OrganizationUser) UpdatedDate() time.Time {
	return u.UpdatedAt
}

func (u *OrganizationUser) DeletedDate() time.Time {
	return u.DeletedAt
}

type Employee struct {
	ID             uuid.UUID
	EmployeeName   string
	UserID         uuid.UUID
	OrganizationId uuid.UUID
	WalletAddress  []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

func (u *Employee) Id() uuid.UUID {
	return u.ID
}

func (u *Employee) ParticipantName() string {
	return u.EmployeeName
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

func (u *Employee) IsOwner() bool {
	return false
}

func (u *Employee) Position() string {
	return "" // todo
}

func (u *Employee) IsActive() bool {
	return u.DeletedAt.IsZero()
}

func (u *Employee) CreatedDate() time.Time {
	return u.CreatedAt
}

func (u *Employee) UpdatedDate() time.Time {
	return u.UpdatedAt
}

func (u *Employee) DeletedDate() time.Time {
	return u.DeletedAt
}
