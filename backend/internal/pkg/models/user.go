package models

type UserIdentity interface {
	Id() string
	Mnemonic() string
}
