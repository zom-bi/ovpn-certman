package middleware

import (
	"net/http"

	"github.com/zom-bi/ovpn-certman/services"
)

// RequireLogin is a middleware that checks for a username in the active
// session, and redirects to `/login` if no username was found.
func RequireLogin(sessions *services.Sessions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			if username := sessions.GetUsername(req); username == "" {
				http.Redirect(w, req, "/login", http.StatusFound)
			}

			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(fn)
	}
}
