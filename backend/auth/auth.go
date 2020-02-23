package auth

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

type AuthProvider struct {
	Name         string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

type AuthController struct {
	BaseUri   string
	providers map[string]*oauth2.Config
}

func (c *AuthController) AddProvider(provider AuthProvider, endpoint oauth2.Endpoint) *AuthController {

	c.providers[provider.Name] = &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint:     endpoint,
		RedirectURL:  fmt.Sprintf("%s/auth/%s/callback", c.BaseUri, provider.Name),
		Scopes:       provider.Scopes,
	}
	return c
}

func (c *AuthController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerName, ok := vars["provider"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	for k, v := range c.providers {
		if strings.EqualFold(k, providerName) {
			if strings.HasSuffix(r.URL.Path, "callback") {
				// TODO: validate state and CSRF
				token, err := v.Exchange(r.Context(), r.FormValue("code"))
				if err != nil {
					log.Println(err)
				}
				client := v.Client(r.Context(), token)
				client.
			} else {
				// TODO: generate a real state parameter for CSRF
				http.Redirect(w, r, v.AuthCodeURL(""), http.StatusTemporaryRedirect)
			}
		}
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
