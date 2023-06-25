package model

import (
	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/rfcomm"
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
