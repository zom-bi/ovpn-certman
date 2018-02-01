package models

import (
	"errors"
	"time"
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
