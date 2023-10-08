package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Pyotr23/the-box/bot/internal/app/module"
)

type appModule interface {
	Name() string
	Init(ctx context.Context, app any) error
	SuccessLog()
	Close(ctx context.Context) error
	CloseLog()
}

type (
	App struct {
		modules         []appModule
		tunnel          net.Listener
		shutdownStartCh chan struct{}
	}
)

func NewApp(ctx context.Context) (*App, error) {
	var a = new(App)
	return a, a.init(ctx)
}

func (a *App) Run(ctx context.Context) (chan struct{}, chan error) {
	log.Println("run app")

	errCh := make(chan error)

	// go func() {
	// 	errCh <- http.Serve(a.mediator.tunnel, http.HandlerFunc(a.handleUpdate))
	// }()

	return a.shutdownStartCh, errCh
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

func (a *App) SetTunnel(tunnel net.Listener) {
	a.tunnel = tunnel
}

func (a *App) GetAddr() string {
	return a.tunnel.Addr().String()
}

func (a *App) SetShutdownStartChannel(ch chan struct{}) {
	a.shutdownStartCh = ch
}

func (a *App) init(ctx context.Context) error {
	a.modules = []appModule{
		module.NewNgrokTunnel(),
		module.NewWebhook(),
		// newBotManager(),
		// newBluetoothClient(),
		// newMessage(),
		module.NewGracefulShutdown(),
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

// func (a *App) handleUpdate(_ http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		log.Print("not POST update method")
// 		return
// 	}

// 	a.mediator.updateCh <- r.Body
// }
