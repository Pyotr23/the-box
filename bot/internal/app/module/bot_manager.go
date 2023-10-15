package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	bc "github.com/Pyotr23/the-box/bot/internal/client/bluetooth"
	"github.com/Pyotr23/the-box/bot/internal/helper"
	"github.com/Pyotr23/the-box/bot/internal/model"
	bs "github.com/Pyotr23/the-box/bot/internal/service/bluetooth"
	mp "github.com/Pyotr23/the-box/bot/internal/service/message_processor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	botName   = "bot manager"
	chatIdKey = "chatID"
)

type updateChSetter interface {
	SetUpdateChannel(ch chan *json.Decoder)
}

type processor interface {
	Process(ctx context.Context, text string) error
	ProcessCommand(ctx context.Context, text string) error
}

type botManager struct {
	api          *tgbotapi.BotAPI
	bodyCh       chan *json.Decoder
	textChatIdCh chan model.TextChatID
	keyboardCh   chan model.Keyboard
	processor    processor
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

	b.processor = mp.NewService(bluetoothService, b.textChatIdCh, b.keyboardCh)

	us.SetUpdateChannel(b.bodyCh)

	// 	b.inputTextCh = make(chan model.TextChatID)
	// 	b.outputTextCh = make(chan model.TextChatID)
	// 	b.waitInputTextCh = make(chan struct{})
	// 	b.userMessageCh = make(chan string)

	// 	// a.inputMessageCh = b.inputMessageCh
	// 	// a.outputMessageCh = b.outputMessageCh

	go func() {
		for decoder := range b.bodyCh {
			var update tgbotapi.Update
			if err := decoder.Decode(&update); err != nil {
				log.Printf("decode: %s", err)
				continue
			}

			var (
				isCommand bool
				chatID    int64
				text      string
			)
			switch {
			case update.Message != nil:
				isCommand = update.Message.IsCommand()
				chatID = update.Message.Chat.ID
				text = update.Message.Text
			case update.CallbackQuery != nil:
				chatID = update.CallbackQuery.Message.Chat.ID
				text = update.CallbackQuery.Data
			default:
				log.Print("unknown message")
				continue
			}

			ctx := helper.CtxWithChatIdValue(context.Background(), chatID)

			var errText string
			if isCommand {
				if err := b.processor.ProcessCommand(ctx, text); err != nil {
					errText = fmt.Sprintf("process command: %s", err)
				}
			} else {
				if err := b.processor.Process(ctx, text); err != nil {
					errText = fmt.Sprintf("process message: %s", err)
				}
			}

			if errText != "" {
				b.textChatIdCh <- model.TextChatID{
					Text:   errText,
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
				keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(b.Key, b.Key))
			}
			message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardButtons)
			if _, err := b.api.Send(message); err != nil {
				log.Printf("send one time reply keyboard fail: %s", err)
			}
		}
	}()

	// 	go func() {
	// 		for data := range b.outputTextCh {
	// 			helper.Logln(fmt.Sprintf("data for user - %v", data))

	// 			message := tgbotapi.NewMessage(data.ChatID, data.Text)
	// 			if _, err := b.api.Send(message); err != nil {
	// 				helper.Logln(fmt.Sprintf("send fail: %s", err.Error()))
	// 			}
	// 		}
	// 	}()

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

// func (b *botManager) processBody(body io.ReadCloser) {
// 	defer func() {
// 		if err := body.Close(); err != nil {
// 			helper.Logln(fmt.Sprintf("close body: %s", err.Error()))
// 		}
// 	}()

// 	var update *tgbotapi.Update
// 	if err := json.NewDecoder(body).Decode(update); err != nil {
// 		helper.Logln(fmt.Sprintf("decode: %s", err.Error()))
// 		return
// 	}

// 	if update == nil {
// 		helper.Logln("empty update")
// 		return
// 	}

// 	if update.Message != nil {
// 		if update.Message.Text != "" {
// 			b.inputTextCh <- model.TextChatID{
// 				Text:   update.Message.Text,
// 				ChatID: update.Message.Chat.ID,
// 			}
// 		}
// 	}

// 	helper.Logln("empty/unknown update")
// }
