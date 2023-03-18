package app

import (
	"context"
	"fmt"
	"log"
	"os"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	username string
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
