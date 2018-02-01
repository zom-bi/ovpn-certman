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
	db.AutoMigrate(models.User{}, models.Client{})
	db.LogMode(conf.Log)

	return &DB{
		gorm: db,
		conf: conf,
	}
}

// CountUsers returns the number of Users in the datastore
func (db *DB) CountUsers() (uint, error) {
	var count uint
	err := db.gorm.Find(&models.User{}).Count(&count).Error
	return count, err
}

// CreateUser inserts a user into the datastore
func (db *DB) CreateUser(user *models.User) error {
	err := db.gorm.Create(&user).Error
	return err
}

// ListUsers returns a slice of 'count' users, starting at 'offset'
func (db *DB) ListUsers(count, offset int) ([]*models.User, error) {
	var users = make([]*models.User, 0)

	err := db.gorm.Find(&users).Limit(count).Offset(offset).Error

	return users, err
}

// GetUserByID returns a single user by ID
func (db *DB) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := db.gorm.Where("id = ?", id).First(&user).Error
	return &user, err
}

// GetUserByEmail returns a single user by email
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.gorm.Where("email = ?", email).First(&user).Error
	return &user, err
}

// DeleteUser removes a user from the datastore
func (db *DB) DeleteUser(id uint) error {
	var user models.User
	err := db.gorm.Where("id = ?", id).Delete(&user).Error
	return err
}

// CreatePasswordReset creates a new password reset token
func (db *DB) CreatePasswordReset(pwReset *models.PasswordReset) error {
	err := db.gorm.Create(&pwReset).Error
	return err
}

// GetPasswordResetByToken retrieves a PasswordReset by token
func (db *DB) GetPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var pwReset models.PasswordReset
	err := db.gorm.Where("token = ?", token).First(&pwReset).Error
	return &pwReset, err
}

// DeletePasswordResetsByUserID deletes all pending password resets for a user
func (db *DB) DeletePasswordResetsByUserID(uid uint) error {
	err := db.gorm.Where("user_id = ?", uid).Delete(&models.PasswordReset{}).Error
	return err
}
