package user

import (
    "database/sql"
    "errors"
    "time"
)

type DbUser struct {
    Id int
    IdpId string
    FirstName string
    LastName string
    Dob time.Time
    DisplayName string
    Pronoun int
}

type DbPronoun struct {
    Id int
    Label string
}

type DbTosAgreement struct {
    Id int
    UserId int
    Agreed bool
    Ts time.Time
}

type DbUserProfile struct {
    User DbUser
    Pronoun DbPronoun
    TosAgreement *time.Time
}

type Repo struct {
    *sql.DB
}

func (r *Repo) GetUser(userId int) (*DbUser, error) {
    // TODO: implement
    return nil, errors.New("not implemented")
}

func (r *Repo) GetUserByIdp(idpId string) (*DbUser, error) {
    // TODO: implement
    return nil, errors.New("not implemented")
}

func (r *Repo) GetUserProfile(userId int) (*DbUserProfile, error) {
    // TODO: implement
    return nil, errors.New("not implemented")
}

func (r *Repo) AddUser(idpId string, user DbUser) (*DbUser, error) {
    // TODO: implement
    return nil, errors.New("not implemented")
}

func (r *Repo) AgreeTos(userId int) error {
    // TODO: implement
    return errors.New("not implemented")
}