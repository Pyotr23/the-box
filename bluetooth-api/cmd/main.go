package main

import (
	"context"
	"log"
	"os"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/app"
)

type IApp interface {
	Run(ctx context.Context) (chan os.Signal, chan error)
	Exit(ctx context.Context)
}

func main() {
	var (
		a   IApp
		err error
		ctx = context.Background()
	)
	a, err = app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	shutdownCh, errCh := a.Run(ctx)

	select {
	case err := <-errCh:
		log.Printf("http serve: %s\n", err.Error())
	case <-shutdownCh:
		log.Println("shutdown signal")
	}

	a.Exit(ctx)
}
