package model

import (
	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	"github.com/Pyotr23/the-box/internal/model/enum"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct {
	Code   enum.Code
	Socket rfcomm.Socket
	BotAPI *tgbotapi.BotAPI
	ChatID int64
}
