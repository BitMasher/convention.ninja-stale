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
	Id          string
	DisplayName string
	Name        string
	Dob         time.Time
}

func (u *User) GetId() string {
	return u.Id
}

func (u *User) GetDisplayName() string {
	if len(u.DisplayName) > 0 {
		return u.DisplayName
	}
	return u.Name
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
		Name:        dbUser.Name,
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
			Name:        dbUsers[i].Name,
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
		Name:        dbUser.Name,
		Dob:         dbUser.Dob,
	}, nil
}

func (c *Controller) Register(ctx context.Context, name string, displayName string, dob *time.Time, token *jwt.JWT) (*User, error) {
	// TODO: do a better name validation check
	if len(strings.Trim(name, " \r\n\t")) == 0 {
		return nil, errors.New("name is a required field")
	}
	dbUser, err := c.Repo.RegisterOauthUser(ctx, token.Public()["prov"].(string), token.Subject(), name, dob, token.Public()["email"].(string), displayName)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("failed to register user, please try again")
	}
	return &User{
		Id:          dbUser.Id,
		DisplayName: dbUser.DisplayName,
		Name:        dbUser.Name,
		Dob:         dbUser.Dob,
	}, nil
}
