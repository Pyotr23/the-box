package app

import (
	"context"
	"log"
	"strconv"

	"github.com/Pyotr23/the-box/internal/cqrs/command"
	"github.com/Pyotr23/the-box/internal/cqrs/query"
	"github.com/Pyotr23/the-box/internal/cqrs/unknown"
	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	"github.com/Pyotr23/the-box/internal/model"
	"github.com/Pyotr23/the-box/internal/model/enum"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type updateHandler struct {
	botAPI    *tgbotapi.BotAPI
	inputCh   chan byte
	waitInput bool
}

type botCommandHandler interface {
	Handle()
}

func (h *updateHandler) handle(update *tgbotapi.Update, socket rfcomm.Socket) {
	if update == nil {
		log.Println("nil update")
		return
	}
	if update.Message == nil {
		log.Println("nil message")
		return
	}
	if update.Message.Text == "" {
		log.Println("empty message")
		return
	}

	text := update.Message.Text

	log.Printf("text from bot '%s'\n", text)

	if h.waitInput {
		n, _ := strconv.ParseUint(text, 10, 8)
		log.Printf("text from bot '%d'\n", n)
		h.inputCh <- byte(n)
		h.waitInput = false
		return
	}

	c := model.Command{
		Code:   enum.GetCode(text),
		Socket: socket,
		BotAPI: h.botAPI,
		ChatID: update.Message.Chat.ID,
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
		h.waitInput = true
		handler = command.NewCallbackCommandHandler(c, h.inputCh)
	case enum.GetIDCode:
		handler = query.NewQueryHandler(c)
	case enum.UnknownCode:
		log.Printf("unknown command '%s'\n", text)
		handler = unknown.NewUnknownCommandHandler(h.botAPI, update.Message.Chat.ID)
	}

	handler.Handle()
}

func (h *updateHandler) Init(ctx context.Context, a *App) error {
	a.updateHandler = &updateHandler{
		botAPI:  *&a.botAPI,
		inputCh: make(chan byte),
	}
	return nil
}

func (n *updateHandler) SuccessLog() {
	log.Println("setup update handler")
}
