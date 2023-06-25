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
	signalCh chan os.Signal
	runCh    chan struct{}
}

func newGracefulShutdown() *gracefulShutdown {
	return &gracefulShutdown{
		signalCh: make(chan os.Signal, 1),
		runCh:    make(chan struct{}),
	}
}

func (gs *gracefulShutdown) Init(ctx context.Context, mediator *mediator) error {
	go func() {
		signal.Notify(gs.signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-gs.signalCh

		close(gs.runCh)
	}()

	mediator.shutdownStartCh = gs.runCh

	return nil
}

func (*gracefulShutdown) SuccessLog() {
	log.Println("setup graceful shutdown")
}

func (gs *gracefulShutdown) Close(ctx context.Context) error {
	close(gs.signalCh)
	return nil
}

func (gs *gracefulShutdown) CloseLog() {
	log.Println("start graceful shutdown")
}

func (*gracefulShutdown) Name() string {
	return gracefulShutdownName
}
