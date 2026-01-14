package models

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-telegram/bot"
)

var (
	ErrBot error = new(BotError)

	ProcessBotError bot.ErrorsHandler = func(err error) {
		slog.Error(err.Error())

		var be *BotError
		if errors.As(err, &be) {
			be.Send()
		}
	}
)

type BotError struct {
	bot    *bot.Bot
	chatID int64
	ctx    context.Context
	err    error
}

func (be *BotError) Error() string {
	return be.err.Error()
}

func (be *BotError) Send() {
	be.bot.SendMessage(be.ctx, &bot.SendMessageParams{
		ChatID: be.chatID,
		Text:   be.err.Error(),
	})
}

func NewBotError(ctx context.Context, bot *bot.Bot, chatID int64, err error) *BotError {
	return &BotError{
		bot:    bot,
		chatID: chatID,
		ctx:    ctx,
		err:    err,
	}
}
