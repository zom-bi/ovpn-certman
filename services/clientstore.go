package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"git.klink.asia/paul/certman/models"
)

var (
	ErrNilCertificate  = errors.New("Trying to store nil certificate")
	ErrDuplicate       = errors.New("Client with that name already exists")
	ErrUserNotExists   = errors.New("User does not exist")
	ErrClientNotExists = errors.New("Client does not exist")
)

type ClientCollection struct {
	sync.RWMutex
	path string

	Clients   map[uint]*models.Client
	UserIndex map[string]map[string]uint
	LastID    uint
}

func NewClientCollection(path string) *ClientCollection {
	// empty collection
	var clientCollection = ClientCollection{
		path:      path,
		Clients:   make(map[uint]*models.Client),
		UserIndex: make(map[string]map[string]uint),
		LastID:    0,
	}

	raw, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return &clientCollection
	} else if err != nil {
		log.Println(err)
		return &clientCollection
	}

	if err := json.Unmarshal(raw, &clientCollection); err != nil {
		log.Println(err)
	}

	return &clientCollection
}

// CreateClient inserts a client into the datastore
func (db *ClientCollection) CreateClient(client *models.Client) error {
	db.Lock()
	defer db.Unlock()

	if client == nil {
		return ErrNilCertificate
	}

	db.LastID++ // increment Id
	client.ID = db.LastID

	userIndex, exists := db.UserIndex[client.User]
	if !exists {
		// create user index if not exists
		db.UserIndex[client.User] = make(map[string]uint)
		userIndex = db.UserIndex[client.User]
	}

	if _, exists = userIndex[client.Name]; exists {
		return ErrDuplicate
	}

	// if all went well, add client and set the index
	db.Clients[client.ID] = client
	userIndex[client.Name] = client.ID
	db.UserIndex[client.User] = userIndex

	return db.save()
}

// ListClientsForUser returns a slice of 'count' client for user 'user', starting at 'offset'
func (db *ClientCollection) ListClientsForUser(user string) ([]*models.Client, error) {
	db.RLock()
	defer db.RUnlock()

	var clients = make([]*models.Client, 0)

	userIndex, exists := db.UserIndex[user]
	if !exists {
		return nil, errors.New("user does not exist")
	}

	for _, clientID := range userIndex {
		clients = append(clients, db.Clients[clientID])
	}

	return clients, nil
}

// GetClientByID returns a single client by ID
func (db *ClientCollection) GetClientByID(id uint) (*models.Client, error) {

	client, exists := db.Clients[id]
	if !exists {
		return nil, ErrClientNotExists
	}

	return client, nil
}

// GetClientByNameUser returns a single client by ID
func (db *ClientCollection) GetClientByNameUser(name, user string) (*models.Client, error) {
	db.RLock()
	defer db.RUnlock()

	userIndex, exists := db.UserIndex[user]
	if !exists {
		return nil, ErrUserNotExists
	}

	clientID, exists := userIndex[name]
	if !exists {
		return nil, ErrClientNotExists
	}

	client, exists := db.Clients[clientID]
	if !exists {
		return nil, ErrClientNotExists
	}

	return client, nil
}

// DeleteClient removes a client from the datastore
func (db *ClientCollection) DeleteClient(id uint) error {
	db.Lock()
	defer db.Unlock()

	client, exists := db.Clients[id]
	if !exists {
		return nil // nothing to delete
	}

	userIndex, exists := db.UserIndex[client.User]
	if !exists {
		return ErrUserNotExists
	}

	delete(userIndex, client.Name) // delete client index

	// if index is now empty, delete the user entry
	if len(userIndex) == 0 {
		delete(db.UserIndex, client.User)
	}

	// finally delete the client
	delete(db.Clients, id)

	return db.save()
}

func (c *ClientCollection) save() error {
	collectionJSON, _ := json.Marshal(c)
	return ioutil.WriteFile(c.path, collectionJSON, 0600)
}
