package handler

import (
	"log"

	"github.com/Pyotr23/the-box/bot/internal/model"
)

type baseHandler struct {
	chatID       int64
	outputTextCh chan<- model.Info
}

func newBaseHandler(chatID int64, outputTextCh chan<- model.Info) baseHandler {
	return baseHandler{
		chatID:       chatID,
		outputTextCh: outputTextCh,
	}
}

func (h baseHandler) ProcessError(err error) {
	log.Println(err.Error())

	// h.outputTextCh <- model.Info{
	// 	Text:   err.Error(),
	// 	ChatID: h.chatID,
	// }
}

func (h baseHandler) SendText(text string) {
	// h.outputTextCh <- model.Message{
	// 	Text:   text,
	// 	ChatID: h.chatID,
	// }
}
