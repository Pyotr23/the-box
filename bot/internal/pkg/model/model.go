package model

import (
	"github.com/Pyotr23/the-box/bot/internal/pkg/enum"
)

const ErrMessageNoChatID = "chat id not found in context"

type TextChatID struct {
	ChatID int64
	Text   string
}

type Keyboard struct {
	ChatID  int64
	Message string
	Buttons []Button
}

type Button struct {
	Key   string
	Value string
}

type Info struct {
	ChatID       int64
	OutputTextCh chan TextChatID
	Code         enum.Code
}
