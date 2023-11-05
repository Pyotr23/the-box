package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	bc "github.com/Pyotr23/the-box/bot/internal/pkg/client/bluetooth"
	"github.com/Pyotr23/the-box/bot/internal/pkg/helper"
	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	bs "github.com/Pyotr23/the-box/bot/internal/pkg/service/bluetooth"
	sts "github.com/Pyotr23/the-box/bot/internal/pkg/service/settings"
	"github.com/Pyotr23/the-box/bot/internal/pkg/sm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	botName   = "bot manager"
	chatIdKey = "chatID"
)

type updateChSetter interface {
	SetUpdateChannel(ch chan *json.Decoder)
}

type process func(ctx context.Context, text string) error

type processorGetter interface {
	GetCommandProcessor() func(ctx context.Context, text string) error
	GetTextProcessor() func(ctx context.Context, text string) error
}

type botManager struct {
	api             *tgbotapi.BotAPI
	bodyCh          chan *json.Decoder
	textChatIdCh    chan model.TextChatID
	keyboardCh      chan model.Keyboard
	processorGetter processorGetter
}

func NewBotManager() *botManager {
	return &botManager{}
}

func (*botManager) Name() string {
	return botName
}

func (b *botManager) Init(ctx context.Context, app any) (err error) {
	us, ok := app.(updateChSetter)
	if !ok {
		return errors.New("app not implements update channel setter")
	}

	token := os.Getenv(botTokenEnv)
	if token == "" {
		return fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	b.api, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	b.bodyCh = make(chan *json.Decoder)
	b.textChatIdCh = make(chan model.TextChatID)
	b.keyboardCh = make(chan model.Keyboard)

	bluetoothClient, err := bc.NewClient()
	if err != nil {
		return fmt.Errorf("bluetooth client: %w", err)
	}

	bluetoothService := bs.NewService(bluetoothClient)
	settingsService := sts.NewService()

	b.processorGetter = sm.NewFsmProcessor(bluetoothService, settingsService, b.textChatIdCh, b.keyboardCh)

	us.SetUpdateChannel(b.bodyCh)

	go func() {
		for decoder := range b.bodyCh {
			var update tgbotapi.Update
			if err := decoder.Decode(&update); err != nil {
				log.Printf("decode: %s", err)
				continue
			}

			var (
				chatID  int64
				text    string
				process process
			)
			switch {
			case update.Message != nil:
				chatID = update.Message.Chat.ID
				text = update.Message.Text
				if update.Message.IsCommand() {
					process = b.processorGetter.GetCommandProcessor()
				} else {
					process = b.processorGetter.GetTextProcessor()
				}
			case update.CallbackQuery != nil:
				chatID = update.CallbackQuery.Message.Chat.ID
				text = update.CallbackQuery.Data
				process = b.processorGetter.GetTextProcessor()
			default:
				log.Print("unknown message")
				continue
			}

			ctx := helper.CtxWithChatIdValue(context.Background(), chatID)

			if err := process(ctx, text); err != nil {
				b.textChatIdCh <- model.TextChatID{
					Text:   err.Error(),
					ChatID: chatID,
				}
			}
		}
	}()

	go func() {
		for textChatID := range b.textChatIdCh {
			message := tgbotapi.NewMessage(textChatID.ChatID, textChatID.Text)
			if _, err := b.api.Send(message); err != nil {
				log.Printf("send message fail: %s", err)
			}
		}
	}()

	go func() {
		for keyboard := range b.keyboardCh {
			var message = tgbotapi.NewMessage(keyboard.ChatID, keyboard.Message)
			var keyboardButtons = make([]tgbotapi.InlineKeyboardButton, 0, len(keyboard.Buttons))
			for _, b := range keyboard.Buttons {
				keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(b.Key, b.Value))
			}
			message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardButtons)
			if _, err := b.api.Send(message); err != nil {
				b.textChatIdCh <- model.TextChatID{
					ChatID: keyboard.ChatID,
					Text:   fmt.Sprintf("send one time reply keyboard fail: %s", err),
				}
			}
		}
	}()

	return nil
}

func (*botManager) SuccessLog() {
	log.Println("ready bot manager")
}

func (b *botManager) Close(_ context.Context) error {
	// TODO write to closed channel
	close(b.bodyCh)
	close(b.textChatIdCh)
	close(b.keyboardCh)
	return nil
}

func (*botManager) CloseLog() {
	log.Printf("graceful shutdown of module '%s'", botName)
}
