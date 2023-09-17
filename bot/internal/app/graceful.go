package app

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const gracefulShutdownName = "graceful shutdown"

type shutdownStartChSetter interface {
	setShutdownStartChannel(ch chan struct{})
}

type gracefulShutdown struct {
	signalCh chan os.Signal
	runCh    chan struct{}
}

func newGracefulShutdown() *gracefulShutdown {
	return &gracefulShutdown{}
}

func (gs *gracefulShutdown) Init(_ context.Context, app interface{}) error {
	setter, ok := app.(shutdownStartChSetter)
	if !ok {
		return errors.New("app not implements shutdown start channel setter")
	}

	gs.signalCh = make(chan os.Signal, 1)
	gs.runCh = make(chan struct{})

	go func() {
		signal.Notify(gs.signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-gs.signalCh

		close(gs.runCh)
	}()

	setter.setShutdownStartChannel(gs.runCh)

	return nil
}

func (*gracefulShutdown) SuccessLog() {
	log.Print("setup graceful shutdown")
}

func (gs *gracefulShutdown) Close(ctx context.Context) error {
	close(gs.signalCh)
	return nil
}

func (gs *gracefulShutdown) CloseLog() {
	log.Print("start graceful shutdown")
}

func (*gracefulShutdown) Name() string {
	return gracefulShutdownName
}
