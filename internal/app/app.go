package app

import (
	"log"
	"net/http"

	"github.com/randnull/posts-service/internal/config"
	"github.com/randnull/posts-service/internal/repository"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/randnull/posts-service/internal/graph"
)


type App struct {
	repo            graph.PostsRepository
	resolver        *graph.Resolver

	cfg 			*config.Config
}


func NewApp(cfg *config.Config) *App {
	var repo graph.PostsRepository

	if cfg.RepoType == "in-memory" {
		repo = repository.NewInMemoryRepository(cfg)
	} else {
		repo = repository.NewRepository(cfg)
	}

	resolver := graph.NewResolver(repo)

	return &App{
		repo: repo,
		resolver: resolver,
		cfg: cfg,
	}
}


func (a *App) Run() {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: a.resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", a.cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+a.cfg.ServerPort, nil))
}
