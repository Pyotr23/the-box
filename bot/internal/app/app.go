package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Pyotr23/the-box/bot/internal/pkg/module"
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
		updateCh        chan *json.Decoder
	}
)

func NewApp(ctx context.Context) (*App, error) {
	var a = new(App)
	return a, a.init(ctx)
}

func (a *App) Run(ctx context.Context) (chan struct{}, chan error) {
	log.Println("run app")

	errCh := make(chan error)

	go func() {
		errCh <- http.Serve(a.tunnel, http.HandlerFunc(a.handleUpdate))
	}()

	return a.shutdownStartCh, errCh
}

func (a *App) Exit(ctx context.Context) {
	for _, mod := range a.modules {
		if err := mod.Close(ctx); err != nil {
			log.Print(fmt.Sprintf("failed graceful shutdown of module '%s'\n", mod.Name()))
			continue
		}

		mod.CloseLog()
	}
}

func (a *App) SetUpdateChannel(ch chan *json.Decoder) {
	a.updateCh = ch
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
		module.NewGracefulShutdown(),
		module.NewNgrokTunnel(),
		module.NewWebhook(),
		module.NewBotManager(),
		// newBluetoothClient(),
		// newMessage(),
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

	a.updateCh <- json.NewDecoder(r.Body)

	time.Sleep(time.Second)
}
