package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/handler/command"
	"github.com/Pyotr23/the-box/internal/handler/model"
	"github.com/Pyotr23/the-box/internal/handler/query"
	"github.com/Pyotr23/the-box/internal/rfcomm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const updateHandlerName = "update handler"

type updateHandler struct {
	updateCh     chan *tgbotapi.Update
	inputCh      chan string
	outputTextCh chan model.TextChatID
	waitInputCh  chan struct{}
	socket       rfcomm.Socket
}

type botCommandHandler interface {
	Handle()
}

func newUpdateHandler() *updateHandler {
	return &updateHandler{
		updateCh:     make(chan *tgbotapi.Update),
		inputCh:      make(chan string),
		outputTextCh: make(chan model.TextChatID),
		waitInputCh:  make(chan struct{}),
	}
}

func (*updateHandler) Name() string {
	return updateHandlerName
}

func (h *updateHandler) Init(ctx context.Context, a *App) error {
	a.updateCh = h.updateCh
	h.socket = a.sockets[0]

	go func() {
		for update := range h.updateCh {
			h.handle(update)
		}
	}()

	go func() {
		for tid := range h.outputTextCh {
			message := tgbotapi.NewMessage(tid.ChatID, tid.Text)
			_, err := a.botAPI.Send(message)
			if err != nil {
				log.Printf("send fail: %s\n", err.Error())
			}
		}
	}()

	return nil
}

func (*updateHandler) SuccessLog() {
	log.Println("setup update handler")
}

func (h *updateHandler) Close(ctx context.Context, a *App) error {
	close(h.updateCh)
	close(h.outputTextCh)
	close(h.inputCh)
	close(h.waitInputCh)
	return nil
}

func (*updateHandler) CloseLog() {
	closeLog(updateHandlerName)
}

func (h *updateHandler) handle(update *tgbotapi.Update) {
	if err := validate(update); err != nil {
		log.Printf("not valid update: %s", err.Error())
		return
	}

	text := update.Message.Text

	select {
	case <-h.waitInputCh:
		h.inputCh <- text
		return
	default:
		break
	}

	log.Printf("update text from bot '%s'\n", text)

	c := h.createCommand(update.Message)

	var handler botCommandHandler
	switch c.Code {
	case enum.TemperatureCode:
		handler = query.NewQueryHandler(c)
	case enum.RelayOnCode:
		handler = command.NewCommand(c)
	case enum.RelayOffCode:
		handler = command.NewCommand(c)
	case enum.SetIDCode:
		handler = command.NewSetIDCallbackCommand(c, h.inputCh, h.waitInputCh)
	case enum.GetIDCode:
		handler = query.NewQueryHandler(c)
	case enum.UnknownCode:
		h.outputTextCh <- model.TextChatID{
			Text:   fmt.Sprintf("unknown command '%s'", text),
			ChatID: c.ChatID,
		}
		return
	}

	handler.Handle()
}

func (h *updateHandler) createCommand(msg *tgbotapi.Message) model.Command {
	return model.Command{
		Code:         enum.GetCode(msg.Text),
		Socket:       h.socket,
		ChatID:       msg.Chat.ID,
		OutputTextCh: h.outputTextCh,
	}
}

func validate(update *tgbotapi.Update) error {
	if update == nil {
		return errors.New("nil update")
	}
	if update.Message == nil {
		return errors.New("nil message")
	}
	if update.Message.Text == "" {
		return errors.New("empty message")
	}
	return nil
}
