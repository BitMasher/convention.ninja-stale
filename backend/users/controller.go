package users

import (
	"context"
	"convention.ninja/auth"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbrlsnchs/jwt"
	"strings"
	"time"
)

type User struct {
	Id          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Dob         time.Time `json:"dob"`
}

func (u *User) GetId() string {
	return u.Id
}

func (u *User) GetDisplayName() string {
	if len(u.DisplayName) > 0 {
		return u.DisplayName
	}
	if len(u.LastName) > 0 {
		return u.FirstName + " " + string([]rune(u.LastName)[0]) + "."
	}
	return u.FirstName
}

type Controller struct {
	Repo
}

func (c *Controller) GetUserByOauth(ctx context.Context, profile *auth.OauthProfile) (interface{}, error) {
	dbUser, err := c.Repo.GetUserByProvider(ctx, profile.Provider, profile.Id)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:          dbUser.Id,
		DisplayName: dbUser.DisplayName,
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		Dob:         dbUser.Dob,
	}, nil
}

func (c *Controller) GetUsers(ctx context.Context) ([]User, error) {
	dbUsers, err := c.Repo.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, len(dbUsers))
	for i := range dbUsers {
		users[i] = User{
			Id:          dbUsers[i].Id,
			DisplayName: dbUsers[i].DisplayName,
			FirstName:   dbUsers[i].FirstName,
			LastName:    dbUsers[i].LastName,
			Dob:         dbUsers[i].Dob,
		}
	}
	return users, nil
}

func (c *Controller) GetUser(ctx context.Context, id string) (*User, error) {
	dbUser, err := c.Repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &User{
		Id:          dbUser.Id,
		DisplayName: dbUser.DisplayName,
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		Dob:         dbUser.Dob,
	}, nil
}

func (c *Controller) Register(ctx context.Context, firstName string, lastName string, displayName string, dob *time.Time, token *jwt.JWT) (*User, error) {
	// TODO: do a better name validation check
	// TODO: block restricted names such as admin?
	if len(strings.Trim(firstName, " \r\n\t")) == 0 {
		return nil, errors.New("name is a required field")
	}
	dbUser, err := c.Repo.RegisterOauthUser(ctx, token.Public()["prov"].(string), token.Subject(), firstName, lastName, dob, token.Public()["email"].(string), displayName)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("failed to register user, please try again")
	}
	return &User{
		Id:          dbUser.Id,
		DisplayName: dbUser.DisplayName,
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		Dob:         dbUser.Dob,
	}, nil
}
