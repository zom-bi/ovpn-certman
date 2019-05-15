package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/zom-bi/ovpn-certman/services"
	"golang.org/x/oauth2"

	"github.com/zom-bi/ovpn-certman/assets"
	"github.com/zom-bi/ovpn-certman/handlers"
	"github.com/zom-bi/ovpn-certman/views"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"

	mw "github.com/zom-bi/ovpn-certman/middleware"
)

var (
	// TODO: make this configurable
	csrfCookieName = "csrf"
	csrfFieldName  = "csrf_token"
	csrfKey        = []byte("7Oj4DllZ9lTsxJnisTuWiiQBGQIzi6gX")
	cookieKey      = []byte("osx70sMD8HZG2ouUl8uKI4wcMugiJ2WH")
)

func HandleRoutes(provider *services.Provider) http.Handler {
	mux := chi.NewMux()

	//mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)             // log requests
	mux.Use(middleware.RealIP)             // use proxy headers
	mux.Use(middleware.RedirectSlashes)    // redirect trailing slashes
	mux.Use(mw.Recoverer)                  // recover on panic
	mux.Use(provider.Sessions.Manager.Use) // use session storage

	// TODO: move this code away from here
	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		Scopes:       []string{"read_user"},
		RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("OAUTH2_AUTH_URL"),
			TokenURL: os.Getenv("OAUTH2_TOKEN_URL"),
		},
	}

	// we are serving the static files directly from the assets package
	// this either means we use the embedded files, or live-load
	// from the file system (if `--tags="dev"` is used).
	fileServer(mux, "/static", assets.Assets)

	mux.Route("/", func(r chi.Router) {
		if os.Getenv("ENVIRONMENT") != "test" {
			r.Use(csrf.Protect(
				csrfKey,
				csrf.Secure(false),
				csrf.CookieName(csrfCookieName),
				csrf.FieldName(csrfFieldName),
				csrf.ErrorHandler(http.HandlerFunc(handlers.CSRFErrorHandler)),
			))
		}

		r.HandleFunc("/", http.RedirectHandler("certs", http.StatusFound).ServeHTTP)

		r.Route("/login", func(r chi.Router) {
			r.Get("/", handlers.GetLoginHandler(provider, oauth2Config))
			r.Get("/oauth2/redirect", handlers.OAuth2Endpoint(provider, oauth2Config))
		})

		r.Route("/certs", func(r chi.Router) {
			r.Use(mw.RequireLogin(provider.Sessions))
			r.Get("/", handlers.ListClientsHandler(provider))
			r.Post("/new", handlers.CreateCertHandler(provider))
			r.HandleFunc("/download/{name}", handlers.DownloadCertHandler(provider))
			r.Post("/delete/{name}", handlers.DeleteCertHandler(provider))
		})

		r.Get("/unconfigured-backend", handlers.NotFoundHandler)
	})

	// what should happen if no route matches
	mux.NotFound(handlers.NotFoundHandler)

	return mux
}

// v is a helper function for quickly displaying a view
func v(template string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		view := views.New(req)
		view.Render(w, template)
	}
}

// fileServer sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	//fs := http.StripPrefix(path, http.FileServer(root))
	fs := http.FileServer(root)

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
