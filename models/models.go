package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotImplemented gets thrown if some action was not attempted,
	// because it is not implemented in the code yet.
	ErrNotImplemented = errors.New("Not implemented")
)

// Model is a base model definition, including helpful fields for dealing with
// models in a database
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// User represents a User of the system which is able to log in
type User struct {
	Model
	Email          string
	EmailValid     bool
	DisplayName    string
	HashedPassword []byte
	IsAdmin        bool
}

// SetPassword sets the password of an user struct, but does not save it yet
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = bytes
	return nil
}

// CheckPassword compares a supplied plain text password with the internally
// stored password hash, returns error=nil on success.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}

type UserProvider interface {
	CountUsers() (uint, error)
	CreateUser(*User) (*User, error)
	ListUsers(count, offset int) ([]*User, error)
	GetUserByID(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	DeleteUser(id uint) error
}

// Client represent the OpenVPN client configuration
type Client struct {
	Model
	Name       string
	User       User
	UserID     uint
	Cert       []byte
	PrivateKey []byte
}

type ClientProvider interface {
	CountClients() (uint, error)
	CreateClient(*User) (*User, error)
	ListClients(count, offset int) ([]*User, error)
	GetClientByID(id uint) (*User, error)
	DeleteClient(id uint) error
}
