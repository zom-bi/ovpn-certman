package services

import (
	"errors"
	"log"

	"git.klink.asia/paul/certman/models"
	"github.com/jinzhu/gorm"
)

// Error Definitions
var (
	ErrNotImplemented = errors.New("Not implemented")
)

type DBConfig struct {
	Type string
	DSN  string
	Log  bool
}

// DB is a wrapper around gorm.DB to provide custom methods
type DB struct {
	gorm *gorm.DB

	conf *DBConfig
}

func NewDB(conf *DBConfig) *DB {
	// Establish connection
	db, err := gorm.Open(conf.Type, conf.DSN)
	if err != nil {
		log.Fatalf("Could not open database: %s", err.Error())
	}

	// Migrate models
	db.AutoMigrate(models.Client{})
	db.LogMode(conf.Log)

	return &DB{
		gorm: db,
		conf: conf,
	}
}

// CountClients returns the number of clients in the datastore
func (db *DB) CountClients() (uint, error) {
	var count uint
	err := db.gorm.Find(&models.Client{}).Count(&count).Error
	return count, err
}

// CreateClient inserts a client into the datastore
func (db *DB) CreateClient(client *models.Client) error {
	err := db.gorm.Create(&client).Error
	return err
}

// ListClients returns a slice of 'count' client, starting at 'offset'
func (db *DB) ListClients(count, offset int) ([]*models.Client, error) {
	var clients = make([]*models.Client, 0)

	err := db.gorm.Find(&clients).Limit(count).Offset(offset).Error

	return clients, err
}

// ListClientsForUser returns a slice of 'count' client for user 'user', starting at 'offset'
func (db *DB) ListClientsForUser(user string, count, offset int) ([]*models.Client, error) {
	var clients = make([]*models.Client, 0)

	err := db.gorm.Find(&clients).Where("user = ?", user).Limit(count).Offset(offset).Error

	return clients, err
}

// GetClientByID returns a single client by ID
func (db *DB) GetClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := db.gorm.Where("id = ?", id).First(&client).Error
	return &client, err
}

// GetClientByNameUser returns a single client by ID
func (db *DB) GetClientByNameUser(name, user string) (*models.Client, error) {
	var client models.Client
	err := db.gorm.Where("name = ?", name).Where("user = ?", user).First(&client).Error
	return &client, err
}

// DeleteClient removes a client from the datastore
func (db *DB) DeleteClient(id uint) error {
	err := db.gorm.Where("id = ?", id).Delete(&models.Client{}).Error
	return err
}
