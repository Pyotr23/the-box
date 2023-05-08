package app

import (
	"context"
	"fmt"
	"log"
	"os"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const botName = "bot"

type bot struct {
	username string
}

func newBot() *bot {
	return &bot{}
}

func (*bot) Name() string {
	return botName
}

func (b *bot) Init(ctx context.Context, a *App) error {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	api, err := tgapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	a.botAPI = api
	b.username = api.Self.UserName

	return nil
}

func (b *bot) SuccessLog() {
	log.Printf("authorized on account %s\n", b.username)
}

func (*bot) Close(ctx context.Context, a *App) error {
	a.botAPI.StopReceivingUpdates()
	return nil
}

func (*bot) CloseLog() {
	closeLog(botName)
}
