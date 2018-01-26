package services

import (
	"git.klink.asia/paul/certman/settings"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var Sessions sessions.Store

func InitSession() {
	store := sessions.NewCookieStore(
		securecookie.GenerateRandomKey(32), // signing key
		securecookie.GenerateRandomKey(32), // encryption key
	)
	store.Options.HttpOnly = true
	store.Options.MaxAge = 7 * 24 * 60 * 60 // 1 Week
	store.Options.Secure = settings.Get("ENVIRONMENT", "") == "production"

	Sessions = store
}
