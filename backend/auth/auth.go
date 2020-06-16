package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type ProfileFetcher func(*http.Client) (*OauthProfile, error)
type UserValidator func(context.Context, *OauthProfile) (interface{}, error)

type Provider struct {
	Name         string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Profile      ProfileFetcher
	oauth        *oauth2.Config
}

type Controller struct {
	BaseUri   string
	Validator UserValidator
	providers map[string]*Provider
}

type OauthProfile struct {
	Provider string
	Id       string
	Email    string
	Name     string
}

func (c *Controller) AddProvider(provider Provider, endpoint oauth2.Endpoint) *Controller {
	provider.oauth = &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint:     endpoint,
		RedirectURL:  fmt.Sprintf("%s/auth/%s/callback", c.BaseUri, provider.Name),
		Scopes:       provider.Scopes,
	}
	if c.providers == nil {
		c.providers = make(map[string]*Provider)
	}
	c.providers[provider.Name] = &provider
	return c
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerName, ok := vars["provider"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	for k, v := range c.providers {
		if strings.EqualFold(k, providerName) {
			if strings.HasSuffix(r.URL.Path, "callback") {
				// TODO: validate state and CSRF
				token, err := v.oauth.Exchange(r.Context(), r.FormValue("code"))
				if err != nil {
					log.Println(err)
				}
				client := v.oauth.Client(r.Context(), token)
				if v.Profile != nil {
					oauthUser, err := v.Profile(client)
					if err != nil {
						log.Println(err)
					}
					if c.Validator != nil {
						ret, err := c.Validator(r.Context(), oauthUser)
						log.Printf("%+v\n%+v", ret, err)
						if err != nil {
							// no users found for that oauth account, send to register flow
							if errors.Is(err, sql.ErrNoRows) {
								// TODO: set temporary JWT token with oauth details
								http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
								return
							}
							log.Println(err)
						} else {
							// TODO: fix
							http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
						}
					}
					panic(errors.New("no validator configured"))
				}
			} else {
				// TODO: generate a real state parameter for CSRF
				state := rand.Float64()
				http.Redirect(w, r, v.oauth.AuthCodeURL(fmt.Sprintf("%f", state)), http.StatusTemporaryRedirect)
				return
			}
		}
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
