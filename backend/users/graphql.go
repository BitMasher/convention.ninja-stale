package users

import (
	"convention.ninja/auth"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"time"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A registered user",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "id",
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(DbUser); ok {
					return user.Id, nil
				}
				return nil, nil
			},
			Description: "The users unique id",
		},
	},
})

var registrationDetailsType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "UserRegistration",
	Description: "The details for user registration",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The new users name",
		},
		"displayName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The users display name",
		},
		"dob": &graphql.InputObjectFieldConfig{
			Type:        graphql.DateTime,
			Description: "The users date of birth",
		},
	},
})

func GetQuery(controller Controller) *graphql.Object {
	userSchema := graphql.NewObject(graphql.ObjectConfig{
		Name:        "UserQueryApi",
		Description: "The user interaction api",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Name:        "users",
				Description: "Gets the list of active users",
				Type:        graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token := p.Context.Value("token")
					if token != nil && auth.ValidateToken("api", token.(string)) != nil {
						return controller.GetUsers(p.Context)
					}
					return nil, errors.New("invalid privileges")
				},
			},
			"user": &graphql.Field{
				Name:        "user",
				Description: "Get a specific user by id",
				Type:        userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The id of the user to look up",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token := p.Context.Value("token")
					if token != nil && auth.ValidateToken("api", token.(string)) != nil {
						return controller.GetUser(p.Context, p.Args["id"].(string))
					}
					return nil, errors.New("invalid privileges")
				},
			},
		},
	})

	return userSchema
}

func GetMutation(controller Controller) *graphql.Object {
	userSchema := graphql.NewObject(graphql.ObjectConfig{
		Name:        "UserMutationApi",
		Description: "The user interaction api",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Name:        "register",
				Description: "Registers a new user",
				Type:        userType,
				Args: graphql.FieldConfigArgument{
					"details": &graphql.ArgumentConfig{
						Type:        registrationDetailsType,
						Description: "The details of the registration",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token := p.Context.Value("token")
					if token != nil {
						jot := auth.ValidateToken("reg", token.(string))
						if jot != nil {
							fmt.Printf("details: %+v\n", p.Args)
							if details, ok := p.Args["details"].(map[string]interface{}); ok {
								name, nameOk := details["name"].(string)
								if !nameOk || len(name) == 0 {
									return nil, errors.New("nameError:Name is a required field")
								}
								displayName, displayNameOk := details["displayName"].(string)
								if !displayNameOk {
									displayName = ""
								}
								dob, dobOk := details["dob"].(time.Time)
								if !dobOk {
									return nil, errors.New("dobError:Invalid date of birth received")
								}
								return controller.Register(p.Context, name, displayName, &dob, jot)
							}
						}
					}
					return nil, errors.New("invalid privileges")
				},
			},
		},
	})

	return userSchema
}
