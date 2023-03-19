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
	err := http.Serve(a.tunnel, http.HandlerFunc(a.updateHandler))
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

func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	update, err := a.botAPI.HandleUpdate(r)
	if err != nil {
		log.Printf("handle update: %s", err.Error())
		return
	}

	log.Printf("message text: %s", update.Message.Text)

	answer, err := a.sockets[0].Write(update.Message.Text)
	if err != nil {
		log.Printf("write: %s", err.Error())
		return
	}

	log.Printf("arduino answer: %s", answer)
}
