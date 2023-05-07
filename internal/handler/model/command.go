package model

import (
	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

type Command struct {
	Code         enum.Code
	Socket       rfcomm.Socket
	ChatID       int64
	OutputTextCh chan TextChatID
}
