package module

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const botName = "bot manager"

type updateChSetter interface {
	SetUpdateChannel(ch chan io.ReadCloser)
}

type botManager struct {
	api    *tgbotapi.BotAPI
	bodyCh chan io.ReadCloser
	// inputTextCh     chan model.TextChatID
	// outputTextCh    chan model.TextChatID
	// waitInputTextCh chan struct{}
	// userMessageCh   chan string

	// inputMessageCh  chan model.Message
	// outputMessageCh chan model.Message
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

	b.bodyCh = make(chan io.ReadCloser)

	us.SetUpdateChannel(b.bodyCh)

	// 	mediator.updateCh = b.bodyCh

	// 	b.inputTextCh = make(chan model.TextChatID)
	// 	b.outputTextCh = make(chan model.TextChatID)
	// 	b.waitInputTextCh = make(chan struct{})
	// 	b.userMessageCh = make(chan string)

	// 	// a.inputMessageCh = b.inputMessageCh
	// 	// a.outputMessageCh = b.outputMessageCh

	// 	go func() {
	// 		for readCloser := range b.bodyCh {
	// 			update, err := decode(readCloser)
	// 			if err != nil {
	// 				helper.Logln(err.Error())
	// 				return
	// 			}

	// 			if update == nil {
	// 				helper.Logln("nil update")
	// 				return
	// 			}

	// 			if update.Message != nil {
	// 				// select {
	// 				// case <-b.waitInputTextCh:
	// 				// 	b.userMessageCh <- update.Message.Text
	// 				// default:
	// 				// 	break
	// 				// }

	// 			}
	// 		}
	// 	}()

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
	return nil
}

func (*botManager) CloseLog() {
	log.Printf("graceful shutdown of module '%s'", botName)
}

// func decode(readCloser io.ReadCloser) (*tgbotapi.Update, error) {
// 	defer func() {
// 		if err := readCloser.Close(); err != nil {
// 			helper.Logln(fmt.Sprintf("close body: %s", err.Error()))
// 		}
// 	}()

// 	var update *tgbotapi.Update
// 	if err := json.NewDecoder(readCloser).Decode(update); err != nil {
// 		return nil, fmt.Errorf("decode: %w", err)
// 	}

// 	return update, nil
// }

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
