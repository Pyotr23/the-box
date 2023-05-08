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
	go func() {
		a.shutdownStart = make(chan struct{})

		signal.Notify(gs.c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-gs.c

		a.shutdownStart <- struct{}{}
	}()

	return nil
}

func (*gracefulShutdown) SuccessLog() {
	log.Println("setup graceful shutdown")
}

func (gs *gracefulShutdown) Close(ctx context.Context, a *App) error {
	close(a.shutdownStart)
	close(gs.c)
	return nil
}

func (gs *gracefulShutdown) CloseLog() {
	log.Println("start graceful shutdown")
}

func (*gracefulShutdown) Name() string {
	return gracefulShutdownName
}
