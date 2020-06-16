package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbrlsnchs/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
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

type AuthUser interface {
	GetId() string
	GetDisplayName() string
}

type Controller struct {
	BaseUri   string
	Validator UserValidator
	providers map[string]*Provider
}

var JwtSigningKey = ""

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
						if err != nil {
							// no users found for that oauth account, send to register flow
							if errors.Is(err, sql.ErrNoRows) {
								jwtOpts := &jwt.Options{
									JWTID:          ksuid.New().String(),
									Timestamp:      true,
									ExpirationTime: time.Now().Add(time.Hour),
									Subject:        oauthUser.Id,
									Audience:       "reg",
									Issuer:         "reg",
									KeyID:          "1",
									Public: map[string]interface{}{
										"prov":  oauthUser.Provider,
										"name":  oauthUser.Name,
										"email": oauthUser.Email,
									},
								}
								sig := jwt.HS512(JwtSigningKey)
								signedToken, err := jwt.Sign(sig, jwtOpts)
								if err != nil {
									// TODO: do something with this error
									log.Println(err)
									http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
									return
								}
								http.SetCookie(w, &http.Cookie{
									Name:     "token",
									Value:    signedToken,
									MaxAge:   3600,
									Path:     "/",
									Secure:   false,
									HttpOnly: false,
									SameSite: http.SameSiteStrictMode,
								})
								http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
								return
							} else {
								http.Redirect(w, r, "/login#error", http.StatusTemporaryRedirect)
								return
							}
						} else {
							if authUser, ok := ret.(AuthUser); ok {
								// TODO: handle successful login
								fmt.Printf("Got user %s, %s\n", authUser.GetId(), authUser.GetDisplayName())
								jwtOpts := &jwt.Options{
									JWTID:          ksuid.New().String(),
									Timestamp:      true,
									ExpirationTime: time.Now().Add(time.Hour * 5),
									Subject:        authUser.GetId(),
									Audience:       "api",
									Issuer:         "api",
									KeyID:          "1",
									Public: map[string]interface{}{
										"displayName": authUser.GetDisplayName(),
									},
								}
								sig := jwt.HS512(JwtSigningKey)
								signedToken, err := jwt.Sign(sig, jwtOpts)
								if err != nil {
									// TODO: do something with this error
									log.Println(err)
									http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
									return
								}
								http.SetCookie(w, &http.Cookie{
									Name:     "token",
									Value:    signedToken,
									MaxAge:   3600,
									Path:     "/",
									Secure:   false,
									HttpOnly: false,
									SameSite: http.SameSiteStrictMode,
								})
								http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
								return
							} else {
								http.Redirect(w, r, "/login#error", http.StatusTemporaryRedirect)
								return
							}
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

func GetToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	token := ""
	if strings.HasPrefix(auth, "Bearer") {
		authSplit := strings.SplitN(auth, " ", 2)
		if len(authSplit) == 2 {
			token = authSplit[1]
		}
	}
	if token == "" {
		if cookie, err := r.Cookie("token"); err == nil {
			token = cookie.Value
		}
	}
	return token
}

func ValidateToken(aud string, token string) *jwt.JWT {
	jot, err := jwt.FromString(token)
	if err != nil {
		// TODO: do something with error
		return nil
	}
	sig := jwt.HS512(JwtSigningKey)
	err = jot.Verify(sig)
	if err != nil {
		// token is invalid
		return nil
	}
	algValidator := jwt.AlgorithmValidator(jwt.MethodHS512)
	expValidator := jwt.ExpirationTimeValidator(time.Now())
	audValidator := jwt.AudienceValidator(aud)
	if err := jot.Validate(algValidator, expValidator, audValidator); err != nil {
		return nil
	}
	return jot
}
