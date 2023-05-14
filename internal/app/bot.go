package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Pyotr23/the-box/internal/handler/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const botName = "bot"

type bot struct {
	api             *tgbotapi.BotAPI
	updateCh        chan *tgbotapi.Update
	inputMessageCh  chan model.Message
	outputMessageCh chan model.Message
}

func newBot() *bot {
	return &bot{
		updateCh:        make(chan *tgbotapi.Update),
		inputMessageCh:  make(chan model.Message),
		outputMessageCh: make(chan model.Message),
	}
}

func (*bot) Name() string {
	return botName
}

func (b *bot) Init(ctx context.Context, a *App) (err error) {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	b.api, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	a.updateCh = b.updateCh
	a.inputMessageCh = b.inputMessageCh
	a.outputMessageCh = b.outputMessageCh

	go func() {
		for update := range b.updateCh {
			b.inputMessageCh <- model.Message{
				ChatID: update.Message.Chat.ID,
				Text:   update.Message.Text,
			}
		}
	}()

	go func() {
		for m := range b.outputMessageCh {
			log.Printf("message %v\n", m)
			message := tgbotapi.NewMessage(m.ChatID, m.Text)
			_, err := b.api.Send(message)
			if err != nil {
				log.Printf("send fail: %s\n", err.Error())
			}
		}
	}()

	return nil
}

func (b *bot) SuccessLog() {
	log.Println("ready bot service")
}

func (b *bot) Close(ctx context.Context, _ *App) error {
	close(b.inputMessageCh)
	close(b.outputMessageCh)
	return nil
}

func (*bot) CloseLog() {
	closeLog(botName)
}
