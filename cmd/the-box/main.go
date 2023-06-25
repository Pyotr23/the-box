package main

import (
	"context"
	"log"

	"github.com/Pyotr23/the-box/internal/app"
)

type IApp interface {
	Run(ctx context.Context) (chan struct{}, chan error)
	Exit(ctx context.Context)
}

func main() {
	var (
		a   IApp
		err error
		ctx = context.Background()
	)
	a, err = app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	shutdownCh, errCh := a.Run(ctx)

	select {
	case err := <-errCh:
		log.Printf("http serve: %s\n", err.Error())
	case <-shutdownCh:
		log.Println(" <- get shutdown signal")
	}

	log.Println("was select")

	a.Exit(ctx)
}
