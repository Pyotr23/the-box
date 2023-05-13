package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Pyotr23/the-box/internal/handler/model"
	"github.com/Pyotr23/the-box/internal/rfcomm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.ngrok.com/ngrok"
)

const (
	botTokenEnv = "THEBOX_BOTTOKEN"
)

type IApp interface {
	Run(ctx context.Context) (chan struct{}, chan error)
	Exit(ctx context.Context)
}

type Module interface {
	Name() string
	Init(ctx context.Context, app *App) error
	SuccessLog()
	Close(ctx context.Context, app *App) error
	CloseLog()
}

type App struct {
	modules         []Module
	tunnel          ngrok.Tunnel
	sockets         []rfcomm.Socket
	shutdownStart   chan struct{}
	updateCh        chan *tgbotapi.Update
	inputMessageCh  chan model.Message
	outputMessageCh chan model.Message
}

func NewApp(ctx context.Context) (IApp, error) {
	a := &App{}
	return a, a.init(ctx)
}

func (a *App) Run(ctx context.Context) (chan struct{}, chan error) {
	log.Println("run app")

	errCh := make(chan error)

	go func() {
		errCh <- http.Serve(a.tunnel, http.HandlerFunc(a.handleUpdate))
	}()

	return a.shutdownStart, errCh
}

func (a *App) Exit(ctx context.Context) {
	for i := len(a.modules) - 1; i >= 0; i-- {
		m := a.modules[i]

		if err := m.Close(ctx, a); err != nil {
			log.Printf("failed graceful shutdown of module '%s'", m.Name())
			continue
		}

		m.CloseLog()
	}
}

func (a *App) init(ctx context.Context) error {
	a.modules = []Module{
		newNgrokTunnel(),
		newWebhook(),
		newBot(),
		newBluetooth(),
		newMessage(),
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
