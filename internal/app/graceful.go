package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pyotr23/the-box/internal/helper"
)

const gracefulShutdownName = "graceful shutdown"

type gracefulShutdown struct {
	signalCh chan os.Signal
	runCh    chan struct{}
}

func newGracefulShutdown() *gracefulShutdown {
	return &gracefulShutdown{}
}

func (gs *gracefulShutdown) Init(_ context.Context, mediator *mediator) error {
	gs.signalCh = make(chan os.Signal, 1)
	gs.runCh = make(chan struct{})

	go func() {
		signal.Notify(gs.signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-gs.signalCh

		close(gs.runCh)
	}()

	mediator.shutdownStartCh = gs.runCh

	return nil
}

func (*gracefulShutdown) SuccessLog() {
	helper.Logln("setup graceful shutdown")
}

func (gs *gracefulShutdown) Close(ctx context.Context) error {
	close(gs.signalCh)
	return nil
}

func (gs *gracefulShutdown) CloseLog() {
	helper.Logln("start graceful shutdown")
}

func (*gracefulShutdown) Name() string {
	return gracefulShutdownName
}
