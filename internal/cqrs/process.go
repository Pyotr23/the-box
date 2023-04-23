package cqrs

import (
	"log"

	"github.com/Pyotr23/the-box/internal/cqrs/command"
	"github.com/Pyotr23/the-box/internal/cqrs/query"
	"github.com/Pyotr23/the-box/internal/cqrs/unknown"
	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
	"github.com/Pyotr23/the-box/internal/model"
	"github.com/Pyotr23/the-box/internal/model/enum"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Handle()
}

func Process(update *tgapi.Update, socket rfcomm.Socket, botAPI *tgapi.BotAPI) {
	text := update.Message.Text

	log.Printf("text from bot '%s'\n", text)

	c := model.Command{
		Code:   enum.GetCode(text),
		Socket: socket,
		BotAPI: botAPI,
		ChatID: update.Message.Chat.ID,
	}

	var handler Handler
	switch c.Code {
	case enum.TemperatureCode:
		handler = query.NewQueryHandler(c)
	case enum.RelayOnCode:
		handler = command.NewCommandHandler(c)
	case enum.RelayOffCode:
		handler = command.NewCommandHandler(c)
	case enum.UnknownCode:
		log.Printf("unknown command '%s'\n", text)
		handler = unknown.NewUnknownCommandHandler(botAPI, update.Message.Chat.ID)
	}

	handler.Handle()
}
