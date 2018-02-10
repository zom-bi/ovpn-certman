package services

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
)

var (
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
	SessionName string
	CookieKey   string
	HttpOnly    bool
	Secure      bool
	Lifetime    time.Duration
}

// Sessions is a wrapped scs.Store in order to implement custom logic
type Sessions struct {
	*scs.Manager
}

// NewSessions populates the default sessions Store
func NewSessions(conf *SessionsConfig) *Sessions {
	store := scs.NewCookieManager(
		conf.CookieKey,
	)
	store.Name(conf.SessionName)
	store.HttpOnly(true)
	store.Lifetime(conf.Lifetime)
	store.Secure(conf.Secure)

	return &Sessions{store}
}

func (store *Sessions) GetUsername(req *http.Request) string {
	if store == nil {
		// if store was not initialized, all requests fail
		log.Println("Nil pointer when checking session for username")
		return ""
	}

	sess := store.Load(req)

	email, err := sess.GetString(UserEmailKey)
	if err != nil {
		// Username found
		return ""

	}

	// User is logged in
	return email
}

func (store *Sessions) SetUsername(w http.ResponseWriter, req *http.Request, username string) {
	if store == nil {
		// if store was not initialized, do nothing
		return
	}

	sess := store.Load(req)

	// renew token to avoid session pinning/fixation attack
	sess.RenewToken(w)

	sess.PutString(w, UserEmailKey, username)

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
	var flashes []Flash

	sess := store.Load(req)

	if err := sess.GetObject(FlashesKey, &flashes); err != nil {
		return err
	}

	flashes = append(flashes, flash)

	return sess.PutObject(w, FlashesKey, flashes)
}

// Flashes returns a slice of flash messages from session data
func (store *Sessions) Flashes(w http.ResponseWriter, req *http.Request) []Flash {
	var flashes []Flash
	sess := store.Load(req)
	sess.PopObject(w, FlashesKey, &flashes)
	return flashes
}
