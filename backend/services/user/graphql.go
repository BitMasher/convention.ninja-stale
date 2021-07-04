package user

import (
    "errors"
    "firebase.google.com/go/auth"
    "fmt"
    "github.com/graphql-go/graphql"
)

var contactType = graphql.NewObject(graphql.ObjectConfig{
    Name:        "Contact",
    Description: "A form of contact for the user",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type:        graphql.NewNonNull(graphql.ID),
            Description: "The contacts unique id",
        },
        "contactType": &graphql.Field{
            Type:        graphql.NewNonNull(graphql.String),
            Description: "The type of contact method",
        },
        "value": &graphql.Field{
            Type: graphql.NewNonNull(graphql.String),
        },
    },
})
var userType = graphql.NewObject(graphql.ObjectConfig{
    Name:        "User",
    Description: "A registered user",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type:        graphql.NewNonNull(graphql.ID),
            Description: "The users unique id",
        },
        "firstName": &graphql.Field{
            Type: graphql.String,
        },
        "lastName": &graphql.Field{
            Type: graphql.String,
        },
        "displayName": &graphql.Field{
            Type: graphql.NewNonNull(graphql.String),
        },
        "pronoun": &graphql.Field{
            Type: graphql.NewNonNull(graphql.String),
        },
        "dob": &graphql.Field{
            Type: graphql.DateTime,
        },
        "tosAgreement": &graphql.Field{
            Type: graphql.DateTime,
        },
        "contacts": &graphql.Field{
            Type: graphql.NewList(graphql.NewNonNull(contactType)),
        },
    },
})

var newUserInput = graphql.NewInputObject(graphql.InputObjectConfig{
    Name:        "NewUserInput",
    Fields:      graphql.InputObjectConfigFieldMap{
        "firstName": &graphql.InputObjectFieldConfig {
            Type: graphql.NewNonNull(graphql.String),
        },
        "lastName": &graphql.InputObjectFieldConfig{
            Type: graphql.NewNonNull(graphql.String),
        },
        "displayName": &graphql.InputObjectFieldConfig{
            Type: graphql.String,
        },
        "pronoun": &graphql.InputObjectFieldConfig{
            Type: graphql.Int,
        },
        "tosAgreement": &graphql.InputObjectFieldConfig{
            Type: graphql.Boolean,
        },
    },
})

type UserGql struct {
    Controller *Controller
}

func (t *UserGql) GetQuery() *graphql.Object {
    userSchema := graphql.NewObject(graphql.ObjectConfig{
        Name:        "UserQueryApi",
        Description: "The API for querying user details",
        Fields: graphql.Fields{
            "accessToken": &graphql.Field {
                Name: "accessToken",
                Type: graphql.String,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    idToken, ok := p.Context.Value("idtoken").(*auth.Token)
                    if !ok {
                        return nil, errors.New("unauthorized")
                    }
                    return t.Controller.GetAccessToken(p.Context, idToken)
                },
            },
            "me": &graphql.Field{
                Name: "me",
                Type: userType,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    // TODO: implement
                    return nil, errors.New("not implemented")
                },
                Description: "Gets the logged in user",
            },
            "user": &graphql.Field{
                Name: "user",
                Type: userType,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    // TODO: implement
                    return nil, errors.New("not implemented")
                },
                Description: "Gets a user by their id",
            },
        },
    })
    return userSchema
}

func (t *UserGql) GetMutation() *graphql.Object {
    userSchema := graphql.NewObject(graphql.ObjectConfig{
        Name:        "UserMutationApi",
        Description: "The API for mutating a user",
        Fields: graphql.Fields{
            "newUser": &graphql.Field{
                Type: userType,
                Args: graphql.FieldConfigArgument{
                    "input": &graphql.ArgumentConfig{
                        Type: newUserInput,
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    fmt.Println("should return error")
                    // TODO: implement
                    return nil, errors.New("not implemented")
                },
            },
        },
    })
    return userSchema
}
