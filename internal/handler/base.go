package base

import (
	"log"

	"github.com/Pyotr23/the-box/internal/handler/model"
)

type BaseHandler struct {
	chatID       int64
	outputTextCh chan<- model.TextChatID
}

func NewBaseHandler(chatID int64, outputTextCh chan<- model.TextChatID) BaseHandler {
	return BaseHandler{
		chatID:       chatID,
		outputTextCh: outputTextCh,
	}
}

func (h BaseHandler) ProcessError(err error) {
	log.Println(err.Error())

	h.outputTextCh <- model.TextChatID{
		Text:   err.Error(),
		ChatID: h.chatID,
	}
}
