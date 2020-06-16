package users

import (
	"errors"
	"github.com/graphql-go/graphql"
)

func GetSchema(controller Controller) *graphql.Object {

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A registered user",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Name: "Id",
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

	userSchema := graphql.NewObject(graphql.ObjectConfig{
		Name:        "UserApi",
		Description: "The user interaction api",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Name:        "users",
				Description: "Gets the list of active users",
				Type:        graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return controller.GetUsers(p.Context)
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

					return nil, errors.New("invalid id supplied")

					//return controller.GetUser(p.Context, p.Args["id"].(string))
				},
			},
		},
	})

	return userSchema
}
