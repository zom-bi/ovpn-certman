package middleware

import (
	"net/http"

	"git.klink.asia/paul/certman/services"
)

// RequireLogin is a middleware that checks for a username in the active
// session, and redirects to `/login` if no username was found.
func RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if username := services.SessionStore.GetUserEmail(req); username == "" {
			http.Redirect(w, req, "/login", http.StatusFound)
		}

		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
