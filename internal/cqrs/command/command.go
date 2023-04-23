package command

import (
	"fmt"
	"log"

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
}

func NewCommandHandler(c model.Command) CommandHandler {
	return CommandHandler{
		code:   c.Code,
		socket: c.Socket,
		botAPI: c.BotAPI,
		chatID: c.ChatID,
	}
}

func (h CommandHandler) Handle() {
	err := h.socket.Command(h.code)
	if err == nil {
		return
	}

	errMsg := fmt.Sprintf("query: %s", err.Error())
	log.Println(errMsg)

	message := tgapi.NewMessage(h.chatID, errMsg)

	_, err = h.botAPI.Send(message)
	if err != nil {
		log.Printf("send: %s\n", err.Error())
	}
}
