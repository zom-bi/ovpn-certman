package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	// ErrNotImplemented gets thrown if some action was not attempted,
	// because it is not implemented in the code yet.
	ErrNotImplemented = errors.New("Not implemented")
)

// User represents a User of the system which is able to log in
type User struct {
	gorm.Model
	Username       string
	HashedPassword []byte
	IsAdmin        bool
}

// SetPassword sets the password of an user struct, but does not save it yet
func (u *User) SetPassword(password string) error {
	return ErrNotImplemented
}

// CheckPassword compares a supplied plain text password with the internally
// stored password hash, returns error=nil on success.
func (u *User) CheckPassword(password string) error {
	return ErrNotImplemented
}

// ClientConf represent the OpenVPN client configuration
type ClientConf struct {
	gorm.Model
	Name       string
	User       User
	Cert       []byte
	PublicKey  []byte
	PrivateKey []byte
}
