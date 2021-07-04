package user

import (
    kms "cloud.google.com/go/kms/apiv1"
    "context"
    "convention.ninja/utilities"
    "errors"
    "firebase.google.com/go/auth"
    "github.com/SermoDigital/jose/jws"
    kms2 "google.golang.org/genproto/googleapis/cloud/kms/v1"
    "time"
)

type Controller struct {
    Repo *Repo
    KeyPath string
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

func (c *Controller) GetAccessToken(ctx context.Context, idToken *auth.Token) ([]byte, error) {
    userId, err := c.Repo.GetUserByIdp(ctx, idToken.UID)
    if err != nil {
        return nil, errors.New("temporary error try again")
    }
    if userId == nil {
        userId, err = c.Repo.AddUser(ctx, idToken.UID)
        if err != nil {
            return nil, errors.New("temporary error try again")
        }
    }
    // TODO: create access token
    token := jws.NewJWT(jws.Claims{}, utilities.SigningMethodKRsaPss)
    token.Claims().Set("sub", userId)
    token.Claims().Set("iss", "convention.ninja")

    ctxBackground := context.Background()
    kmsclient, err := kms.NewKeyManagementClient(ctxBackground)
    if err != nil {
        return nil, errors.New("kms failure")
    }
    defer kmsclient.Close()
    keyVersionIt := kmsclient.ListCryptoKeyVersions(ctxBackground, &kms2.ListCryptoKeyVersionsRequest{
        Parent:    c.KeyPath,
        PageSize:  1,
        View:      0,
        Filter:    "state=\"ENABLED\"",
        OrderBy:   "name desc",
    })
    keyVersion, err := keyVersionIt.Next()
    if err != nil {
        return nil, errors.New("kms failure")
    }
    tokenSerialized, err := token.Serialize(&kms2.AsymmetricSignRequest{
        Name:         keyVersion.Name,
    })
    if err != nil {
        return nil, errors.New("token generation failed")
    }
    return tokenSerialized, nil
}

func (c *Controller) GetUser(ctx context.Context, id int) (*User, error) {
    return nil, errors.New("not implemented")
}

func (c *Controller) CreateUser(ctx context.Context, parameters NewUserParameters) (*User, error) {
    return nil, errors.New("not implemented")
}

