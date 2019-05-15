package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/zom-bi/ovpn-certman/views"
	"golang.org/x/oauth2"

	"github.com/zom-bi/ovpn-certman/services"
)

func OAuth2Endpoint(p *services.Provider, config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.NewWithSession(req, p.Sessions)

		code := req.FormValue("code")

		// exchange code for token
		accessToken, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, req)
			return
		}

		if accessToken.Valid() {
			// generate a client using the access token
			httpClient := config.Client(oauth2.NoContext, accessToken)

			apiRequest, err := http.NewRequest("GET", os.Getenv("USER_ENDPOINT"), nil)
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

func GetLoginHandler(p *services.Provider, config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authURL := config.AuthCodeURL("", oauth2.AccessTypeOnline)
		http.Redirect(w, req, authURL, http.StatusFound)
	}
}
