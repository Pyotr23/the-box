package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	b "github.com/Pyotr23/the-box/bot/internal/client/bluetooth"
)

type module interface {
	Name() string
	Init(ctx context.Context, app interface{}) error
	SuccessLog()
	Close(ctx context.Context) error
	CloseLog()
}

type (
	mediator struct {
		// tunnel          net.Listener
		shutdownStartCh chan struct{}
		updateCh        chan io.ReadCloser
		bluetoothClient b.BluetoothClient
	}

	App struct {
		modules         []module
		mediator        *mediator
		tunnel          net.Listener
		shutdownStartCh chan struct{}
	}
)

func NewApp(ctx context.Context) (*App, error) {
	a := &App{
		mediator: &mediator{},
	}
	return a, a.init(ctx)
}

func (a *App) Run(ctx context.Context) (chan struct{}, chan error) {
	log.Println("run app")

	errCh := make(chan error)

	// go func() {
	// 	errCh <- http.Serve(a.mediator.tunnel, http.HandlerFunc(a.handleUpdate))
	// }()

	return a.mediator.shutdownStartCh, errCh
}

func (a *App) Exit(ctx context.Context) {
	for i := len(a.modules) - 1; i >= 0; i-- {
		m := a.modules[i]

		if err := m.Close(ctx); err != nil {
			log.Print(fmt.Sprintf("failed graceful shutdown of module '%s'\n", m.Name()))
			continue
		}

		m.CloseLog()
	}
}

func (a *App) setTunnel(tunnel net.Listener) {
	a.tunnel = tunnel
}

func (a *App) getTunnel() net.Listener {
	return a.tunnel
}

func (a *App) setShutdownStartChannel(ch chan struct{}) {
	a.shutdownStartCh = ch
}

func (a *App) init(ctx context.Context) error {
	a.modules = []module{
		newNgrokTunnel(),
		newWebhook(),
		// newBotManager(),
		// newBluetoothClient(),
		// newMessage(),
		newGracefulShutdown(),
	}

	for _, module := range a.modules {
		err := module.Init(ctx, a)
		if err != nil {
			return err
		}

		module.SuccessLog()
	}

	return nil
}

func (a *App) handleUpdate(_ http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("not POST update method")
		return
	}

	a.mediator.updateCh <- r.Body
}

func closeLog(name string) {
	log.Print(fmt.Sprintf("graceful shutdown of module '%s'\n", name))
}
