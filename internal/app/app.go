package app

import (
	"context"
	"log"
	"net/http"

	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.ngrok.com/ngrok"
)

const (
	botTokenEnv = "THEBOX_BOTTOKEN"
)

type Runner interface {
	Run(ctx context.Context) error
}

type Initter interface {
	Init(ctx context.Context, app *App) error
	SuccessLog()
}

type App struct {
	tunnel         ngrok.Tunnel
	botAPI         *tgapi.BotAPI
	sockets        []rfcomm.Socket
	updateHandler  *updateHandler
	shutdownStart  chan struct{}
	shutdownFinish chan struct{}
}

func NewApp(ctx context.Context) (Runner, error) {
	a := &App{
		shutdownStart:  make(chan struct{}),
		shutdownFinish: make(chan struct{}),
	}
	err := a.init(ctx)
	return a, err
}

func (a *App) Run(ctx context.Context) error {
	log.Println("run app")
	err := http.Serve(a.tunnel, http.HandlerFunc(a.handleUpdate))
	if _, opened := <-a.shutdownStart; opened {
		return err
	}
	<-a.shutdownFinish
	log.Println("graceful stop listen bot")
	return nil
}

func (a *App) init(ctx context.Context) error {
	inits := []Initter{
		&ngrokTunnel{},
		&webhook{},
		&bot{},
		&bluetooth{},
		&updateHandler{},
		&gracefulShutdown{},
	}

	for _, initter := range inits {
		err := initter.Init(ctx, a)
		if err != nil {
			return err
		}
		initter.SuccessLog()
	}

	return nil
}

func (a *App) handleUpdate(w http.ResponseWriter, r *http.Request) {
	update, err := a.botAPI.HandleUpdate(r)
	if err != nil {
		log.Printf("handle update: %s", err.Error())
		return
	}

	a.updateHandler.handle(update, a.sockets[0])
}
