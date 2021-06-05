package user

import (
    "errors"
    "time"
)

type Controller struct {
    Repo *Repo
}

type User struct {
    Id int
    FirstName string
    LastName string
    DisplayName string
    Pronoun string
    Dob time.Time
    TosAgreement time.Time
    Contacts []Contact
}

type Contact struct {
    Id int
    ContactType string
    Value string
}

type NewUserParameters struct {
    FirstName string
    LastName string
    DisplayName string
    Pronoun int
    TosAgreement bool
    IdpId string
}

func (c *Controller) GetUser(id int) (*User, error) {
    return nil, errors.New("not implemented")
}

func (c *Controller) CreateUser(parameters NewUserParameters) (*User, error) {
    return nil, errors.New("not implemented")
}

