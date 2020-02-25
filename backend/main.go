package main

import (
	"convention.ninja/auth"
	facebookFetch "convention.ninja/auth/facebook"
	googleFetch "convention.ninja/auth/google"
	"convention.ninja/users"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	// HTTP configs
	ServerPort int    `env:"PORT" env-default:"3000" env-description:"The port to run the HTTP server on"`
	BaseUri    string `env:"BASEURI" env-default:"http://localhost:3000" env-description:"The base url to use for oauth redirection"`
	// Database configs
	DbConnString string `env:"DBCONNSTRING" env-description:"The connection string for database"`
	DbMaxConn    int    `env:"DBMAXCONN" env-default:"5" env-description:"The maximum number of pooled database connections"`
	// OAUTH configs
	GoogleClientId       string `env:"GOOGLECLIENTID" env-description:"Your google oauth client id"`
	GoogleClientSecret   string `env:"GOOGLECLIENTSECRET" env-description:"Your google oauth client secret"`
	FacebookClientId     string `env:"FACEBOOKCLIENTID" env-description:"Your facebook oauth client id"`
	FacebookClientSecret string `env:"FACEBOOKCLIENTSECRET" env-description:"Your facebook oauth client secret"`
}

func main() {

	var config ServerConfig

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		panic(err)
	}

	// TODO: provide better error for invalid port numbers
	// TODO: check more invalid port numbers than just 0
	if config.ServerPort == 0 {
		panic("Invalid port number supplied")
	}

	// TODO: provide proper feedback for invalid database connection strings
	if len(config.DbConnString) == 0 {
		panic("Invalid database connection string provided")
	}

	db, err := sql.Open("postgres", config.DbConnString)
	// TODO: provide better feedback on sql init failure
	if err != nil {
		log.Print(err)
		panic("Failed to initialize database driver")
	}
	db.SetMaxOpenConns(config.DbMaxConn)

	err = db.Ping()
	if err != nil {
		log.Print(err)
		panic("Failed to connect to database")
	}

	userController := users.Controller{Repo: users.Repo{DB: db}}
	userGql := users.GetSchema(userController)

	// TODO: create me resolver
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:        userGql,
				Description: "The user api",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return userGql, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        rootQuery,
		Mutation:     nil,
		Subscription: nil,
		Types:        nil,
		Directives:   nil,
		Extensions:   nil,
	})

	router := mux.NewRouter()

	router.Handle("/graphql", handler.New(&handler.Config{
		Schema:           &schema,
		Pretty:           true,
		GraphiQL:         true,
		Playground:       true,
		RootObjectFn:     nil,
		ResultCallbackFn: nil,
		FormatErrorFn:    nil,
	}))

	authController := auth.Controller{
		BaseUri: config.BaseUri,
		Validator: userController.GetUserByOauth,
	}
	authController.AddProvider(auth.Provider{
		Name:         "google",
		ClientID:     config.GoogleClientId,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{
			"profile",
			"email",
		},
		Profile: googleFetch.FetchProfile,
	}, google.Endpoint).AddProvider(auth.Provider{
		Name:         "facebook",
		ClientID:     config.FacebookClientId,
		ClientSecret: config.FacebookClientSecret,
		Scopes: []string{
			"email",
			"public_profile",
		},
		Profile: facebookFetch.FetchProfile,
	}, facebook.Endpoint)

	router.PathPrefix("/auth/{provider}").Handler(&authController)
	staticFs := RedirectingFileSystem{http.Dir("static")}
	router.PathPrefix("/").Handler(http.FileServer(staticFs))

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Starting HTTP server")
	log.Fatal(srv.ListenAndServe())
}
