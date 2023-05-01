package command

import (
	"fmt"
	"strconv"

	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	"github.com/Pyotr23/the-box/internal/model"
	"github.com/Pyotr23/the-box/internal/model/enum"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	code   enum.Code
	socket rfcomm.Socket
	botAPI *tgapi.BotAPI
	chatID int64
	errCh  chan model.ErrorChatID
}

func NewCommandHandler(c model.Command) CommandHandler {
	return CommandHandler{
		code:   c.Code,
		socket: c.Socket,
		botAPI: c.BotAPI,
		chatID: c.ChatID,
		errCh:  c.ErrorCh,
	}
}

func (h CommandHandler) Handle() {
	err := h.socket.Command(h.code)
	go func() {
		h.errCh <- model.ErrorChatID{
			Err:    fmt.Errorf("command handle: %w", err),
			ChatID: h.chatID,
		}
	}()
}

type CallbackCommandHandler struct {
	code        enum.Code
	socket      rfcomm.Socket
	botAPI      *tgapi.BotAPI
	chatID      int64
	inputCh     <-chan string
	errCh       chan<- model.ErrorChatID
	waitInputCh chan struct{}
}

func NewCallbackCommandHandler(c model.Command, inputCh <-chan string, waitInputCh chan struct{}) CallbackCommandHandler {
	return CallbackCommandHandler{
		code:        c.Code,
		socket:      c.Socket,
		botAPI:      c.BotAPI,
		chatID:      c.ChatID,
		inputCh:     inputCh,
		waitInputCh: waitInputCh,
		errCh:       c.ErrorCh,
	}
}

func (h CallbackCommandHandler) Handle() {
	go func() {
		var err error
		defer func() {
			if err != nil {
				h.errCh <- model.ErrorChatID{
					Err:    fmt.Errorf("callback command handle: %w", err),
					ChatID: h.chatID,
				}
			}
		}()

		h.waitInputCh <- struct{}{}

		n, err := strconv.ParseUint(<-h.inputCh, 10, 8)
		if err != nil {
			err = fmt.Errorf("parse uint: %w", err)
			return
		}

		err = h.socket.SendText(h.code, []byte{byte(n)})
		if err != nil {
			err = fmt.Errorf("send text %d: %w", byte(n), err)
			return
		}
	}()
}
