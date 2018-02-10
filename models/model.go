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

// Client represent the OpenVPN client configuration
type Client struct {
	ID         uint
	CreatedAt  time.Time
	Name       string
	User       string
	Cert       []byte
	PrivateKey []byte
}

type ClientProvider interface {
	CountClients() (uint, error)
	CreateClient(*Client) (*Client, error)
	ListClients(count, offset int) ([]*Client, error)
	ListClientsForUser(user string, count, offset int) ([]*Client, error)
	GetClientByID(id uint) (*Client, error)
	DeleteClient(id uint) error
}
