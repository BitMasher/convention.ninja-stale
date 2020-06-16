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
	FirstName   string
	LastName    string
	Dob         time.Time
}

func (repo *Repo) GetActiveUsers(ctx context.Context) ([]DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id, u.display_name, u.first_name, u.last_name, dob FROM users u")
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]DbUser, 0)
	for rows.Next() {
		var id string
		var displayName string
		var firstName string
		var lastName string
		var dob time.Time
		if err = rows.Scan(&id, &displayName, &firstName, &lastName, &dob); err != nil {
			return nil, err
		}
		users = append(users, DbUser{
			Id:          id,
			DisplayName: displayName,
			FirstName:   firstName,
			LastName:    lastName,
			Dob:         dob,
		})
	}

	return users, nil
}

func (repo *Repo) GetUserById(ctx context.Context, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.display_name, u.first_name, u.last_name, u.dob FROM users u WHERE u.Id = $1", id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var displayName string
		var firstName string
		var lastName string
		var dob time.Time
		if err = rows.Scan(&displayName, &firstName, &lastName, &dob); err != nil {
			return nil, err
		}
		return &DbUser{
			Id:          id,
			DisplayName: displayName,
			FirstName:   firstName,
			LastName:    lastName,
			Dob:         dob,
		}, nil
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}

func (repo *Repo) GetUserByProvider(ctx context.Context, provider string, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id, u.display_name, u.first_name, u.last_name, u.dob FROM users u INNER JOIN user_oauth_providers p ON p.user_id = u.id AND p.provider = $1 AND p.id = $2", provider, id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var userId string
		var displayName string
		var firstName string
		var lastName string
		var dob time.Time
		if err = rows.Scan(&userId, &displayName, &firstName, &lastName, &dob); err != nil {
			return nil, err
		}
		return &DbUser{
			Id:          userId,
			DisplayName: displayName,
			FirstName:   firstName,
			LastName:    lastName,
			Dob:         dob,
		}, nil
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}

func (repo *Repo) RegisterOauthUser(ctx context.Context, provider string, providerId string, firstName string, lastName string, dob *time.Time, email string, displayName string) (*DbUser, error) {
	userId := ksuid.New().String()
	// TODO: check for id conflicts
	row := repo.DB.QueryRowContext(ctx, "INSERT INTO users (id,display_name,first_name,last_name,dob) VALUES($1, $2, $3, $4, $5) RETURNING users.*", userId, displayName, firstName, lastName, dob)
	var dbUserId string
	var userDisplayName string
	var userFirstName string
	var userLastName string
	var userDob time.Time
	if err := row.Scan(&dbUserId, &userDisplayName, &userFirstName, &userLastName, &userDob); err != nil {
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
		FirstName:   userFirstName,
		LastName:    userLastName,
		Dob:         userDob,
	}, nil
}
