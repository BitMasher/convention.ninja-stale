package facebook

import (
	"convention.ninja/auth"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type facebookProfile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}

func FetchProfile(client *http.Client) (*auth.OauthProfile, error) {
	res, err := client.Get("https://graph.facebook.com/v3.2/me?fields=name,email")
	if err != nil {
		return nil, err
	}

	buff, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var profile facebookProfile
	err = json.Unmarshal(buff, &profile)
	if err != nil {
		return nil, err
	}
	log.Printf("Got profile: %+v", profile)
	return &auth.OauthProfile{
		Provider: "facebook",
		Id:       profile.Id,
		Email:    profile.Email,
		Name:     profile.Name,
	}, nil
}
