package users

import (
	"context"
	"database/sql"
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
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u")
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]DbUser, 0)
	for rows.Next() {
		var dbUser DbUser
		if err = rows.Scan(&dbUser); err != nil {
			users = append(users, dbUser)
		}
	}

	return users, nil
}

func (repo *Repo) GetUserById(ctx context.Context, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u WHERE u.Id = $1", id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var dbUser DbUser
		if err = rows.Scan(&dbUser); err != nil {
			return &dbUser, nil
		}
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}

func (repo *Repo) GetUserByProvider(ctx context.Context, provider string, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u INNER JOIN user_oauth_providers p ON p.user_id = u.id AND p.provider = $1 AND p.id = $2", provider, id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var dbUser DbUser
		if err = rows.Scan(&dbUser); err != nil {
			return &dbUser, nil
		}
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}

func (repo *Repo) GetUserByFacebook(ctx context.Context, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u WHERE u.FacebookId = $1", id)
	// TODO: return better errors on sql failure
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var dbUser DbUser
		if err = rows.Scan(&dbUser); err != nil {
			return &dbUser, nil
		}
	}
	// TODO: need a not found error
	return nil, sql.ErrNoRows
}
