package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Pyotr23/the-box/internal/helper"
	"github.com/Pyotr23/the-box/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const botName = "bot manager"

type botManager struct {
	api             *tgbotapi.BotAPI
	bodyCh          chan io.ReadCloser
	inputTextCh     chan model.TextChatID
	outputTextCh    chan model.TextChatID
	waitInputTextCh chan struct{}
	userMessageCh   chan string

	// inputMessageCh  chan model.Message
	// outputMessageCh chan model.Message
}

func newBotManager() *botManager {
	return &botManager{
		// inputMessageCh:  make(chan model.Message),
		// outputMessageCh: make(chan model.Message),
	}
}

func (*botManager) Name() string {
	return botName
}

func (b *botManager) Init(ctx context.Context, mediator *mediator) (err error) {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	b.api, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	b.bodyCh = make(chan io.ReadCloser)

	mediator.updateCh = b.bodyCh

	b.inputTextCh = make(chan model.TextChatID)
	b.outputTextCh = make(chan model.TextChatID)
	b.waitInputTextCh = make(chan struct{})
	b.userMessageCh = make(chan string)

	// a.inputMessageCh = b.inputMessageCh
	// a.outputMessageCh = b.outputMessageCh

	go func() {
		for readCloser := range b.bodyCh {
			update, err := decode(readCloser)
			if err != nil {
				helper.Logln(err.Error())
				return
			}

			if update == nil {
				helper.Logln("nil update")
				return
			}

			if update.Message != nil {
				// select {
				// case <-b.waitInputTextCh:
				// 	b.userMessageCh <- update.Message.Text
				// default:
				// 	break
				// }

			}
		}
	}()

	go func() {
		for data := range b.outputTextCh {
			helper.Logln(fmt.Sprintf("data for user - %v", data))

			message := tgbotapi.NewMessage(data.ChatID, data.Text)
			if _, err := b.api.Send(message); err != nil {
				helper.Logln(fmt.Sprintf("send fail: %s", err.Error()))
			}
		}
	}()

	return nil
}

func (b *botManager) SuccessLog() {
	helper.Logln("ready bot manager")
}

func (b *botManager) Close(ctx context.Context) error {
	close(b.inputTextCh)
	close(b.outputTextCh)
	return nil
}

func (*botManager) CloseLog() {
	closeLog(botName)
}

func decode(readCloser io.ReadCloser) (*tgbotapi.Update, error) {
	defer func() {
		if err := readCloser.Close(); err != nil {
			helper.Logln(fmt.Sprintf("close body: %s", err.Error()))
		}
	}()

	var update *tgbotapi.Update
	if err := json.NewDecoder(readCloser).Decode(update); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return update, nil
}

func (b *botManager) processBody(body io.ReadCloser) {
	defer func() {
		if err := body.Close(); err != nil {
			helper.Logln(fmt.Sprintf("close body: %s", err.Error()))
		}
	}()

	var update *tgbotapi.Update
	if err := json.NewDecoder(body).Decode(update); err != nil {
		helper.Logln(fmt.Sprintf("decode: %s", err.Error()))
		return
	}

	if update == nil {
		helper.Logln("empty update")
		return
	}

	if update.Message != nil {
		if update.Message.Text != "" {
			b.inputTextCh <- model.TextChatID{
				Text:   update.Message.Text,
				ChatID: update.Message.Chat.ID,
			}
		}
	}

	helper.Logln("empty/unknown update")
}
