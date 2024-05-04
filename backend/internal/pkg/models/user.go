package models

type UserIdentity interface {
	Id() string
	Mnemonic() string
	IsAdmin() bool
}

type User struct {
	id       string
	mnemonic string
	isAdmin  bool
}

func NewUser(id string, mnemonic string) *User {
	return &User{
		id:       id,
		mnemonic: mnemonic,
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Mnemonic() string {
	return u.mnemonic
}

func (u *User) IsAdmin() bool {
	return u.isAdmin
}

type OrganizationUser struct {
	User
	// add org info
}
