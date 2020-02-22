package users

import (
	"context"
	"database/sql"
)

type Repo struct {
	*sql.DB
}

type DbUser struct {
	Id         string
	GoogleId   string
	FacebookId string
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
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u WHERE u.Id = ?", id)
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

func (repo *Repo) GetUserByGoogle(ctx context.Context, id string) (*DbUser, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u WHERE u.GoogleId = ?", id)
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
	rows, err := repo.DB.QueryContext(ctx, "SELECT u.Id FROM users u WHERE u.FacebookId = ?", id)
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
