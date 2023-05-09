package handler

import (
	"log"
)

type baseHandler struct {
	chatID       int64
	outputTextCh chan<- TextChatID
}

func newBaseHandler(chatID int64, outputTextCh chan<- TextChatID) baseHandler {
	return baseHandler{
		chatID:       chatID,
		outputTextCh: outputTextCh,
	}
}

func (h baseHandler) ProcessError(err error) {
	log.Println(err.Error())

	h.outputTextCh <- TextChatID{
		Text:   err.Error(),
		ChatID: h.chatID,
	}
}

func (h baseHandler) SendText(text string) {
	h.outputTextCh <- TextChatID{
		Text:   text,
		ChatID: h.chatID,
	}
}
