package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"git.klink.asia/paul/certman/views"
	"golang.org/x/oauth2"

	"git.klink.asia/paul/certman/services"
)

var GitlabConfig = &oauth2.Config{
	ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
	Scopes:       []string{"read_user"},
	RedirectURL:  os.Getenv("HOST") + "/login/oauth2/redirect",
	Endpoint: oauth2.Endpoint{
		AuthURL:  os.Getenv("OAUTH2_AUTH_URL"),
		TokenURL: os.Getenv("OAUTH2_TOKEN_URL"),
	},
}

func OAuth2Endpoint(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.NewWithSession(req, p.Sessions)

		code := req.FormValue("code")

		// exchange code for token
		accessToken, err := GitlabConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, req)
			return
		}

		if accessToken.Valid() {
			// generate a client using the access token
			httpClient := GitlabConfig.Client(oauth2.NoContext, accessToken)

			apiRequest, err := http.NewRequest("GET", "https://git.klink.asia/api/v4/user", nil)
			if err != nil {
				v.RenderError(w, http.StatusNotFound)
				return
			}

			resp, err := httpClient.Do(apiRequest)
			if err != nil {
				fmt.Println(err.Error())
				v.RenderError(w, http.StatusInternalServerError)
				return
			}

			var user struct {
				Username string `json:"username"`
			}

			err = json.NewDecoder(resp.Body).Decode(&user)
			if err != nil {
				fmt.Println(err.Error())
				v.RenderError(w, http.StatusInternalServerError)
				return
			}

			if user.Username != "" {
				p.Sessions.SetUsername(w, req, user.Username)
				http.Redirect(w, req, "/certs", http.StatusFound)
				return
			}

			fmt.Println(err.Error())
			v.RenderError(w, http.StatusInternalServerError)
			return
		}
	}
}

func GetLoginHandler(p *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authURL := GitlabConfig.AuthCodeURL("", oauth2.AccessTypeOnline)
		http.Redirect(w, req, authURL, http.StatusFound)
	}
}
