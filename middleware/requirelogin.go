package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"git.klink.asia/paul/certman/handlers"
)

func RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Println(rvr)
				log.Println(string(debug.Stack()))
				handlers.ErrorHandler(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
