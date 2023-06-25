package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Pyotr23/the-box/internal/rfcomm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.ngrok.com/ngrok"
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
		tunnel          ngrok.Tunnel
		sockets         []rfcomm.Socket
		shutdownStartCh chan struct{}
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
	a := &App{}
	return a, a.init(ctx)
}

func (a *App) Run(ctx context.Context) (chan struct{}, chan error) {
	log.Println("run app")

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
			log.Printf("failed graceful shutdown of module '%s'", m.Name())
			continue
		}

		m.CloseLog()
	}
}

func (a *App) init(ctx context.Context) error {
	a.modules = []module{
		newNgrokTunnel(),
		newWebhook(),
		newBot(),
		newBluetooth(),
		newMessage(),
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
		log.Println("not POST update method")
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("decode: %s\n", err.Error())
		return
	}

	a.updateCh <- &update
}

func closeLog(name string) {
	log.Printf("graceful shutdown of module '%s'\n", name)
}
