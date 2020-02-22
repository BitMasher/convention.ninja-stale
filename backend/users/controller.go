package users

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	Id string
}

type Controller struct {
	Repo
}

func (c *Controller) GetUsers(ctx context.Context) ([]User, error) {
	dbUsers, err := c.Repo.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, len(dbUsers))
	for i := range dbUsers {
		users[i] = User{
			dbUsers[i].Id,
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
		dbUser.Id,
	}, nil
}
