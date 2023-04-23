package query

import (
	"fmt"
	"log"

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
}

func NewQueryHandler(c model.Command) QueryHandler {
	return QueryHandler{
		code:   c.Code,
		socket: c.Socket,
		botAPI: c.BotAPI,
		chatID: c.ChatID,
	}
}

func (h QueryHandler) Handle() {
	answer, err := h.socket.Query(h.code)
	if err != nil {
		answer = fmt.Sprintf("query: %s", err.Error())
		log.Println(answer)
	}

	message := tgapi.NewMessage(h.chatID, answer)

	_, err = h.botAPI.Send(message)
	if err != nil {
		log.Printf("send: %s\n", err.Error())
	}
}
