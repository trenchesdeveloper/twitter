package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/trenchesdeveloper/tweeter/domain"
	"github.com/trenchesdeveloper/tweeter/graph"
	"log"
	"net/http"
	"time"

	"github.com/trenchesdeveloper/tweeter/config"
	"github.com/trenchesdeveloper/tweeter/postgres"
)

func main() {
	ctx := context.Background()

	// load config
	config.LoadEnv(".env")

	// create config
	conf := config.New()

	// create db
	db := postgres.New(ctx, conf)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("migrations ran successfully")

	// create router
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	// REPOS
	userRepo := postgres.NewUserRepo(db)

	// SERVICES
	authService := domain.AuthService{
		userRepo,
	}

	//set graphql playground
	router.Get("/", playground.Handler("Twitter Clone", "/query"))

	router.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService: &authService,
				},
			},
		),
	))

	log.Fatal(http.ListenAndServe(":8080", router))
}
