package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const gracefulShutdownName = "graceful shutdown"

type gracefulShutdown struct {
	c chan os.Signal
}

func newGracefulShutdown() *gracefulShutdown {
	return &gracefulShutdown{
		c: make(chan os.Signal, 1),
	}
}

func (gs *gracefulShutdown) Init(ctx context.Context, a *App) error {
	a.shutdownStart = make(chan struct{})

	go func() {
		signal.Notify(gs.c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-gs.c

		close(a.shutdownStart)
	}()

	return nil
}

func (*gracefulShutdown) SuccessLog() {
	log.Println("setup graceful shutdown")
}

func (gs *gracefulShutdown) Close(ctx context.Context, a *App) error {
	close(gs.c)
	return nil
}

func (gs *gracefulShutdown) CloseLog() {
	log.Println("start graceful shutdown")
}

func (*gracefulShutdown) Name() string {
	return gracefulShutdownName
}
