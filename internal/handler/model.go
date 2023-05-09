package handler

import (
	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

type Info struct {
	Code         enum.Code
	Socket       rfcomm.Socket
	ChatID       int64
	OutputTextCh chan TextChatID
}

type TextChatID struct {
	Text   string
	ChatID int64
}
