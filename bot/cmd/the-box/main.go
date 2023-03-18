package main

import (
	"context"
	"log"

	"github.com/Pyotr23/the-box/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
