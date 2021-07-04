package user

import (
    "context"
    "database/sql"
    "errors"
    "time"
)

type DbUser struct {
    Id          int
    IdpId       string
    FirstName   string
    LastName    string
    Dob         time.Time
    DisplayName string
    Pronoun     int
}

type DbPronoun struct {
    Id    int
    Label string
}

type DbTosAgreement struct {
    Id     int
    UserId int
    Agreed bool
    Ts     time.Time
}

type DbUserProfile struct {
    User         DbUser
    Pronoun      DbPronoun
    TosAgreement *time.Time
}

type Repo struct {
    *sql.DB
}

func (r *Repo) GetUser(ctx context.Context, userId int) (*DbUser, error) {
    // TODO: implement
    return nil, errors.New("not implemented")
}

func (r *Repo) GetUserByIdp(ctx context.Context, idpId string) (*int, error) {
    rows, err := r.DB.QueryContext(ctx, "SELECT id FROM users WHERE idp_id = $1", idpId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if rows.Next() {
        var userId int
        if err := rows.Scan(&userId); err != nil {
            return nil, err
        }
        return &userId, nil
    }
    return nil, sql.ErrNoRows
}

func (r *Repo) GetUserProfile(ctx context.Context, userId int) (*DbUserProfile, error) {
    // TODO: implement
    return nil, errors.New("not implemented")
}

func (r *Repo) AddUser(ctx context.Context, idpId string) (*int, error) {
    rows, err := r.DB.QueryContext(ctx, "INSERT INTO users (idp_id) VALUES($1) ON CONFLICT DO NOTHING RETURNING id", idpId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if rows.Next() {
        var userId int
        if err := rows.Scan(&userId); err != nil {
            return nil, err
        }
        return &userId, nil
    }
    return nil, sql.ErrNoRows
}

func (r *Repo) AgreeTos(ctx context.Context, userId int) error {
    // TODO: implement
    return errors.New("not implemented")
}
