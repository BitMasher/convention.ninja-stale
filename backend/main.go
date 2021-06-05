package main

import (
	"convention.ninja/services/user"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
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
	// JWT configs
	TokenSigningKey string `env:"TOKENSIGNINGKEY" env-description:"Your token signing key"`
}

func populateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//r = r.WithContext(context.WithValue(r.Context(), "token", auth.GetToken(r)))
		next.ServeHTTP(w, r)
	})
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
	/*
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
		}*/

	/*userController := users.Controller{Repo: users.Repo{DB: db}}
	  userQueryGql := users.GetQuery(userController)
	  userMutationGql := users.GetMutation(userController)*/

	userQueryGql := user.GetQuery()
	userMutationGql := user.GetMutation()

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:        userQueryGql,
				Description: "The user api",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return userQueryGql, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:        userMutationGql,
				Description: "The user api",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return userMutationGql, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.Handle("/graphql", populateTokenMiddleware(handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})))

	/*auth.JwtSigningKey = config.TokenSigningKey

	  authController := auth.Controller{
	      BaseUri:   config.BaseUri,
	      Validator: userController.GetUserByOauth,
	  }*/

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
