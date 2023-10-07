package model

import (
	"github.com/Pyotr23/the-box/bot/internal/enum"
	"github.com/Pyotr23/the-box/bot/internal/rfcomm"
)

type TextChatID struct {
	ChatID int64
	Text   string
}

type Info struct {
	ChatID       int64
	OutputTextCh chan TextChatID
	Code         enum.Code
	Socket       rfcomm.Socket
}
