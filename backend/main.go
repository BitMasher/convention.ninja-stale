package main

import (
	"convention.ninja/users"
	"database/sql"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	ServerPort   int    `env:"PORT" env-default:"3000" env-description:"The port to run the HTTP server on"`
	DbConnString string `env:"DBCONNSTRING" env-description:"The connection string for database"`
	DbMaxConn    int    `env:"DBMAXCONN" env-default:"5" env-description:"The maximum number of pooled database connections"`
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

	userGql := users.GetSchema(users.Controller{Repo: users.Repo{DB: db}})

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

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.ServerPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	http.Handle("/graphql/", handler.New(&handler.Config{
		Schema:           &schema,
		Pretty:           true,
		GraphiQL:         true,
		Playground:       true,
		RootObjectFn:     nil,
		ResultCallbackFn: nil,
		FormatErrorFn:    nil,
	}))

	staticFs := RedirectingFileSystem{http.Dir("static")}
	http.Handle("/", http.FileServer(staticFs))

	log.Print("Starting HTTP server")
	log.Fatal(srv.ListenAndServe())
}
