package query

import (
	"fmt"

	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	"github.com/Pyotr23/the-box/internal/model"
	"github.com/Pyotr23/the-box/internal/model/enum"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QueryHandler struct {
	code   enum.Code
	socket rfcomm.Socket
	botAPI *tgapi.BotAPI
	chatID int64
	errCh  chan<- model.ErrorChatID
}

func NewQueryHandler(c model.Command) QueryHandler {
	return QueryHandler{
		code:   c.Code,
		socket: c.Socket,
		botAPI: c.BotAPI,
		chatID: c.ChatID,
		errCh:  c.ErrorCh,
	}
}

func (h QueryHandler) Handle() {
	var err error
	defer func() {
		if err != nil {
			h.errCh <- model.ErrorChatID{
				Err:    fmt.Errorf("query handle: %w", err),
				ChatID: h.chatID,
			}
		}
	}()

	answer, err := h.socket.Query(h.code)
	if err != nil {
		err = fmt.Errorf("query: %w", err)
		return
	}

	message := tgapi.NewMessage(h.chatID, answer)

	_, err = h.botAPI.Send(message)
	if err != nil {
		err = fmt.Errorf("send: %w", err)
		return
	}
}
