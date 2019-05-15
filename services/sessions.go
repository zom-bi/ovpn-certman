package services

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
)

const (
	// FlashesKey is the key used for the flashes in the cookie
	FlashesKey = "_flashes"
	// UserEmailKey is the key used to reference usernames
	UserEmailKey = "_user_email"
)

func init() {
	// Register the Flash message type, so gob can serialize it
	gob.Register(Flash{})
}

type SessionsConfig struct {
	HTTPOnly bool
	Secure   bool
	Lifetime time.Duration
}

// Sessions is a wrapped scs.Store in order to implement custom logic
type Sessions struct {
	*scs.Session
}

// NewSessions populates the default sessions Store
func NewSessions(conf *SessionsConfig) *Sessions {
	session := scs.NewSession()
	session.Lifetime = conf.Lifetime
	session.Cookie.HttpOnly = true
	session.Cookie.Secure = conf.Secure

	return &Sessions{session}
}

func (store *Sessions) GetUsername(req *http.Request) string {
	if store == nil {
		// if store was not initialized, all requests fail
		log.Println("Nil pointer when checking session for username")
		return ""
	}

	email := store.GetString(req.Context(), UserEmailKey)
	return email // "" if no user is logged in
}

func (store *Sessions) SetUsername(w http.ResponseWriter, req *http.Request, username string) {
	if store == nil {
		// if store was not initialized, do nothing
		return
	}

	// renew token to avoid session pinning/fixation attack
	store.RenewToken(req.Context())

	store.Put(req.Context(), UserEmailKey, username)
}

type Flash struct {
	Message template.HTML
	Type    string
}

// Render renders the flash message as a notification box
func (flash Flash) Render() template.HTML {
	return template.HTML(
		fmt.Sprintf(
			"<div class=\"notification is-radiusless is-%s\"><div class=\"container has-text-centered\">%s</div></div>",
			flash.Type, flash.Message,
		),
	)
}

// Flash add flash message to session data
func (store *Sessions) Flash(w http.ResponseWriter, req *http.Request, flash Flash) error {
	flashes, ok := store.Get(req.Context(), FlashesKey).([]Flash)
	if !ok {
		return errors.New("Could not get flashes")
	}

	flashes = append(flashes, flash)
	return nil
}

// Flashes returns a slice of flash messages from session data
func (store *Sessions) Flashes(w http.ResponseWriter, req *http.Request) []Flash {
	flashes, ok := store.Pop(req.Context(), FlashesKey).([]Flash)
	if !ok {
		return nil
	}
	return flashes
}
