package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
	CreateUser(*User) error
	ListUsers(count, offset int) ([]*User, error)
	GetUserByID(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	DeleteUser(id uint) error
}

type PasswordReset struct {
	Model
	User       *User
	UserID     uint
	Token      string
	ValidUntil time.Time
}

type PasswordResetProvider interface {
	CreatePasswordReset(*PasswordReset) error
	GetPasswordResetByToken(token string) (*PasswordReset, error)
}
