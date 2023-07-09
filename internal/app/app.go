package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/Pyotr23/the-box/internal/helper"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

const (
	botTokenEnv = "THEBOX_BOTTOKEN"
)

type module interface {
	Name() string
	Init(ctx context.Context, mediator *mediator) error
	SuccessLog()
	Close(ctx context.Context) error
	CloseLog()
}

type (
	mediator struct {
		tunnel          net.Listener
		sockets         []rfcomm.Socket
		shutdownStartCh chan struct{}
		updateCh        chan io.ReadCloser
	}

	App struct {
		modules  []module
		mediator *mediator
		// tunnel          ngrok.Tunnel
		// sockets         []rfcomm.Socket
		// shutdownStart   chan struct{}
		// updateCh        chan *tgbotapi.Update
		// inputMessageCh  chan model.Message
		// outputMessageCh chan model.Message
	}
)

func NewApp(ctx context.Context) (*App, error) {
	a := &App{
		mediator: &mediator{},
	}
	return a, a.init(ctx)
}

func (a *App) Run(ctx context.Context) (chan struct{}, chan error) {
	helper.Logln("run app")

	errCh := make(chan error)

	go func() {
		errCh <- http.Serve(a.mediator.tunnel, http.HandlerFunc(a.handleUpdate))
	}()

	return a.mediator.shutdownStartCh, errCh
}

func (a *App) Exit(ctx context.Context) {
	for i := len(a.modules) - 1; i >= 0; i-- {
		m := a.modules[i]

		if err := m.Close(ctx); err != nil {
			helper.Logln(fmt.Sprintf("failed graceful shutdown of module '%s'", m.Name()))
			continue
		}

		m.CloseLog()
	}
}

func (a *App) init(ctx context.Context) error {
	a.modules = []module{
		// newNgrokTunnel(),
		// newWebhook(),
		// newBotManager(),
		newBluetooth(),
		// newMessage(),
		newGracefulShutdown(),
	}

	for _, module := range a.modules {
		err := module.Init(ctx, a.mediator)
		if err != nil {
			return err
		}

		module.SuccessLog()
	}

	return nil
}

func (a *App) handleUpdate(_ http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.Logln("not POST update method")
		return
	}

	a.mediator.updateCh <- r.Body
}

func closeLog(name string) {
	helper.Logln(fmt.Sprintf("graceful shutdown of module '%s'", name))
}
