package users

import (
	"context"
	"database/sql"
	"github.com/segmentio/ksuid"
	"time"
)

type Repo struct {
	*sql.DB
}

type DbUser struct {
	Id          string
	DisplayName string
	Name        string
	Dob         time.Time
}

func (repo *Repo) GetActiveUsers(ctx context.Context) ([]DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id, u.display_name, u.name FROM users u")
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]DbUser, 0)
	for rows.Next() {
		var id string
		var displayName string
		var name string
		if err = rows.Scan(&id, &displayName, &name); err != nil {
			return nil, err
		}
		users = append(users, DbUser{
			Id:          id,
			DisplayName: displayName,
			Name:        name,
		})
	}

	return users, nil
}

func (repo *Repo) GetUserById(ctx context.Context, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.display_name, u.name, u.dob FROM users u WHERE u.Id = $1", id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var displayName string
		var name string
		var dob time.Time
		if err = rows.Scan(&displayName, &name, &dob); err != nil {
			return nil, err
		}
		return &DbUser{
			Id:          id,
			DisplayName: displayName,
			Name:        name,
			Dob:         dob,
		}, nil
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}

func (repo *Repo) GetUserByProvider(ctx context.Context, provider string, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id, u.display_name, u.name, u.dob FROM users u INNER JOIN user_oauth_providers p ON p.user_id = u.id AND p.provider = $1 AND p.id = $2", provider, id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var userId string
		var displayName string
		var name string
		var dob time.Time
		if err = rows.Scan(&userId, &displayName, &name, &dob); err != nil {
			return nil, err
		}
		return &DbUser{
			Id:          userId,
			DisplayName: displayName,
			Name:        name,
			Dob:         dob,
		}, nil
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}

func (repo *Repo) RegisterOauthUser(ctx context.Context, provider string, providerId string, name string, dob *time.Time, email string, displayName string) (*DbUser, error) {
	userId := ksuid.New().String()
	// TODO: check for id conflicts
	row := repo.DB.QueryRowContext(ctx, "INSERT INTO users (id,display_name,name,dob) VALUES($1, $2, $3, $4) RETURNING users.*", userId, displayName, name, dob)
	var dbUserId string
	var userDisplayName string
	var userName string
	var userDob time.Time
	if err := row.Scan(&dbUserId, &userDisplayName, &userName, &userDob); err != nil {
		return nil, err
	}
	// TODO: check for id conflicts
	_, err := repo.DB.ExecContext(ctx, "INSERT INTO user_oauth_providers (id, provider, user_id) VALUES($1, $2, $3)", providerId, provider, userId)
	if err != nil {
		return nil, err
	}
	return &DbUser{
		Id:          userId,
		DisplayName: userDisplayName,
		Name:        userName,
		Dob:         userDob,
	}, nil
}
