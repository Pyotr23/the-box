package app

import (
	"context"
	"errors"
	"log"

	"github.com/Pyotr23/the-box/internal/cqrs/command"
	"github.com/Pyotr23/the-box/internal/cqrs/query"
	"github.com/Pyotr23/the-box/internal/cqrs/unknown"
	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	"github.com/Pyotr23/the-box/internal/model"
	"github.com/Pyotr23/the-box/internal/model/enum"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type messageHandler struct {
	botAPI      *tgbotapi.BotAPI
	inputCh     chan string
	waitInputCh chan struct{}
	errCh       chan model.ErrorChatID
	socket      rfcomm.Socket
}

type botCommandHandler interface {
	Handle()
}

func (h *messageHandler) handle(update *tgbotapi.Update) {
	var err error
	if update == nil {
		err = errors.New("nil update")
	} else if update.Message == nil {
		err = errors.New("nil message")
	} else if update.Message.Text == "" {
		err = errors.New("empty message")
	}

	if err != nil {
		log.Printf("not valid update: %s", err.Error())
		return
	}

	text := update.Message.Text
	log.Printf("text from bot '%s'\n", text)

	select {
	case <-h.waitInputCh:
		h.inputCh <- text
		return
	default:
		break
	}

	c := model.Command{
		Code:    enum.GetCode(text),
		Socket:  h.socket,
		BotAPI:  h.botAPI,
		ChatID:  update.Message.Chat.ID,
		ErrorCh: h.errCh,
	}

	var handler botCommandHandler
	switch c.Code {
	case enum.TemperatureCode:
		handler = query.NewQueryHandler(c)
	case enum.RelayOnCode:
		handler = command.NewCommandHandler(c)
	case enum.RelayOffCode:
		handler = command.NewCommandHandler(c)
	case enum.SetIDCode:
		handler = command.NewCallbackCommandHandler(c, h.inputCh, h.waitInputCh)
	case enum.GetIDCode:
		handler = query.NewQueryHandler(c)
	case enum.UnknownCode:
		log.Printf("unknown command '%s'\n", text)
		handler = unknown.NewUnknownCommandHandler(h.botAPI, update.Message.Chat.ID)
	}

	handler.Handle()
}

func (h *messageHandler) Init(ctx context.Context, a *App) error {
	a.messageHandler = &messageHandler{
		botAPI:      *&a.botAPI,
		inputCh:     make(chan string),
		waitInputCh: make(chan struct{}),
		socket:      a.sockets[0],
		errCh:       make(chan model.ErrorChatID),
	}

	go func() {
		for eid := range h.errCh {
			log.Println(eid.Err.Error())

			message := tgbotapi.NewMessage(eid.ChatID, eid.Err.Error())

			_, err := h.botAPI.Send(message)
			if err != nil {
				log.Printf("send fail: %s\n", err.Error())
			}
		}
	}()

	return nil
}

func (n *messageHandler) SuccessLog() {
	log.Println("setup update handler")
}
