package message_processor

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Pyotr23/the-box/bot/internal/helper"
	"github.com/Pyotr23/the-box/bot/internal/model"
	"github.com/looplab/fsm"
)

type command string

const (
	search command = "search"
)

type bluetoothService interface {
	Search(ctx context.Context) ([]string, error)
}

type Service struct {
	currentFsmByChatID map[int64]*fsm.FSM
	bs                 bluetoothService
	textChatIdCh       chan model.TextChatID
}

func NewService(bs bluetoothService, textChatIdCh chan model.TextChatID) *Service {
	return &Service{
		bs:                 bs,
		currentFsmByChatID: make(map[int64]*fsm.FSM),
		textChatIdCh:       textChatIdCh,
	}
}

func (s *Service) ProcessCommand(ctx context.Context, text string) error {
	chatID := helper.ChatIdFromCtx(ctx)
	if chatID == 0 {
		return errors.New("chat id not found in context")
	}

	sm := s.createFsm(command(text[1:]))
	if sm == nil {
		return fmt.Errorf("unknown command '%s'", text)
	}

	s.currentFsmByChatID[chatID] = sm

	transitions := sm.AvailableTransitions()
	switch len(transitions) {
	case 0:
		return fmt.Errorf("no transitions for command '%s", text)
	case 1:
		return sm.Event(ctx, transitions[0])
	default:
		return fmt.Errorf("two much transitions for command '%s", text)
	}
}

func (s *Service) createFsm(command command) *fsm.FSM {
	switch command {
	case search:
		return fsm.NewFSM("start",
			fsm.Events{
				{
					Name: "search",
					Src:  []string{"start"},
					Dst:  "finish",
				},
			},
			fsm.Callbacks{
				"enter_finish": func(ctx context.Context, e *fsm.Event) {
					chatID := helper.ChatIdFromCtx(ctx)
					if chatID == 0 {
						log.Print("chat id not found in context")
					}

					macAddresses, err := s.bs.Search(ctx)

					var text string
					if err == nil {
						text = fmt.Sprintf("%v", macAddresses)
					} else {
						text = fmt.Sprintf("search: %s", err)
					}

					s.textChatIdCh <- model.TextChatID{
						Text:   text,
						ChatID: chatID,
					}

					return
				},
			},
		)
	}
	return nil
}
