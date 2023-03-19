package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type gracefulShutdown struct{}

func (*gracefulShutdown) Init(ctx context.Context, a *App) error {
	go func() {
		done := make(chan os.Signal, 1)

		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		close(done)

		close(a.shutdownStart)

		fmt.Println()
		log.Println("start graceful shutdown")

		err := a.tunnel.CloseWithContext(ctx)
		if err != nil {
			log.Printf("tunnel close: %s", err.Error())
			return
		}

		for _, s := range a.sockets {
			s.Close()
		}
		log.Println("close sockets")

		close(a.shutdownFinish)
	}()

	return nil
}

func (*gracefulShutdown) SuccessLog() {
	log.Println("setup graceful shutdown")
}
