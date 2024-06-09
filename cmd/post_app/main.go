package main

import (
	"log"

	"github.com/randnull/posts-service/internal/app"
	"github.com/randnull/posts-service/internal/config"
)

func main() {
	NewConfig, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(NewConfig)
	a.Run()
}
