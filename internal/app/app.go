package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Pyotr23/the-box/configs"
	"github.com/Pyotr23/the-box/internal/models"
	"github.com/go-telegram/bot"
	bm "github.com/go-telegram/bot/models"
)

const workersCount = 1

var botInitTimeout = time.Second * 10

func Run() error {
	appConf, err := configs.NewAppConf()
	if err != nil {
		return fmt.Errorf("new app conf: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithCheckInitTimeout(botInitTimeout),
		bot.WithWorkers(workersCount),
		bot.WithNotAsyncHandlers(),
		bot.WithErrorsHandler(models.ProcessBotError),
	}

	b, err := bot.New(appConf.GetBotToken(), opts...)
	if err != nil {
		return fmt.Errorf("bot new: %w", err)
	}

	slog.Info("start the bot")

	b.Start(ctx)

	slog.Info("stop the bot")

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *bm.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
