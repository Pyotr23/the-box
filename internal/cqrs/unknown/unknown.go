package unknown

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const errText = "unknown command"

type UnknownCommandHandler struct {
	botAPI *tgapi.BotAPI
	chatID int64
}

func NewUnknownCommandHandler(botAPI *tgapi.BotAPI, chatID int64) UnknownCommandHandler {
	return UnknownCommandHandler{
		botAPI: botAPI,
		chatID: chatID,
	}
}

func (h UnknownCommandHandler) Handle() {
	message := tgapi.NewMessage(h.chatID, errText)
	_, err := h.botAPI.Send(message)
	if err != nil {
		log.Printf("send: %s\n", err.Error())
	}
}
