package google

import (
	"convention.ninja/auth"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type googleProfile struct {
	Id string `json:"sub"`
	Email string `json:"email"`
	EmailVerified bool `json:"email_verified"`
	Name string `json:"name"`
}

func FetchProfile(client *http.Client) (*auth.OauthProfile, error) {
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	buff, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	/*
	{
	"sub": "account id",
	"name": "full name",
	"email": "email",
	"email_verified": true,
	}
	*/
	// TODO: parse into unified profile object
	var profile googleProfile
	err = json.Unmarshal(buff, &profile)
	if err != nil {
		return nil, err
	}
	log.Printf("Got profile: %+v", profile)
	oauthProfile := &auth.OauthProfile{
		Provider: "google",
		Id:       profile.Id,
		Name:     profile.Name,
	}
	if profile.EmailVerified {
		oauthProfile.Email = profile.Email
	} else {
		oauthProfile.Email = ""
	}
	return oauthProfile, nil
}