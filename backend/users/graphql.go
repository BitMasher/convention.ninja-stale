package users

import (
	"convention.ninja/auth"
	"errors"
	"github.com/graphql-go/graphql"
	"strings"
	"time"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A registered user",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "The users unique id",
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"displayName": &graphql.Field{
			Type: graphql.String,
		},
		"dob": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

var registrationDetailsType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "UserRegistration",
	Description: "The details for user registration",
	Fields: graphql.InputObjectConfigFieldMap{
		"firstName": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The new users first name",
		},
		"lastName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The new users last name (optional)",
		},
		"displayName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The users display name (optional)",
		},
		"dob": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "The users date of birth",
		},
	},
})

func GetQuery(controller Controller) *graphql.Object {
	userSchema := graphql.NewObject(graphql.ObjectConfig{
		Name:        "UserQueryApi",
		Description: "The user interaction api",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Name:        "me",
				Description: "Gets information about the currently logged in user",
				Type:        userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token := p.Context.Value("token")
					if token != nil {
						tokenData := auth.ValidateToken("api", token.(string))
						if tokenData != nil {
							return controller.GetUser(p.Context, tokenData.Subject())
						}
					}
					return nil, errors.New("invalid privileges")
				},
			},
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
							if details, ok := p.Args["details"].(map[string]interface{}); ok {
								firstName, firstNameOk := details["firstName"].(string)
								if !firstNameOk || len(strings.Trim(firstName, "\r\n\t")) == 0 {
									return nil, errors.New("firstNameError:Name is a required field")
								}
								lastName, lastNameOk := details["lastName"].(string)
								if !lastNameOk {
									lastName = ""
								}
								displayName, displayNameOk := details["displayName"].(string)
								if !displayNameOk {
									displayName = ""
								}
								dob, dobOk := details["dob"].(time.Time)
								if !dobOk {
									return nil, errors.New("dobError:Invalid date of birth received")
								}
								if dob.After(time.Now().Add(-(13 * (time.Hour * 8760)))) {
									return nil, errors.New("dobError:Must be at least 13 years of age to register")
								}
								return controller.Register(p.Context, firstName, lastName, displayName, &dob, jot)
							} else {
								return nil, errors.New("could not parse request")
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
