package message_processor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Pyotr23/the-box/bot/internal/helper"
	"github.com/Pyotr23/the-box/bot/internal/model"
	"github.com/looplab/fsm"
)

const textKey = "text"

type command string

const (
	search command = "search"
	blink  command = "blink"
)

type bluetoothService interface {
	Search(ctx context.Context) ([]string, error)
}

type Service struct {
	currentFsmByChatID map[int64]*fsm.FSM
	bs                 bluetoothService
	textChatIdCh       chan model.TextChatID
	keyboardCh         chan model.Keyboard
}

func NewService(bs bluetoothService, textChatIdCh chan model.TextChatID, keyboardCh chan model.Keyboard) *Service {
	return &Service{
		bs:                 bs,
		currentFsmByChatID: make(map[int64]*fsm.FSM),
		textChatIdCh:       textChatIdCh,
		keyboardCh:         keyboardCh,
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

func (s *Service) Process(ctx context.Context, text string) error {
	log.Print("start process text", text)

	chatID := helper.ChatIdFromCtx(ctx)
	if chatID == 0 {
		return errors.New("chat id not found in context")
	}

	sm := s.currentFsmByChatID[chatID]
	if sm == nil {
		return errors.New("text but no command before")
	}

	sm.SetMetadata(textKey, text)

	transitions := sm.AvailableTransitions()
	log.Print(transitions)
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
				"leave_start": func(ctx context.Context, e *fsm.Event) {
					chatID := helper.ChatIdFromCtx(ctx)
					if chatID == 0 {
						log.Print("chat id not found in context")
					}

					macAddresses, err := s.bs.Search(ctx)

					var text string
					if err == nil {
						text = strings.Join(macAddresses, ",")
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
	case blink:
		return fsm.NewFSM("start",
			fsm.Events{
				{
					Name: "search",
					Src:  []string{"start"},
					Dst:  "choice_waiting",
				},
				{
					Name: "choice",
					Src:  []string{"choice_waiting"},
					Dst:  "blinking",
				},
				// {
				// 	Name: "blink",
				// 	Src:  []string{"blinking"},
				// 	Dst:  "finish",
				// },
			},
			fsm.Callbacks{
				"leave_start": func(ctx context.Context, e *fsm.Event) {
					chatID := helper.ChatIdFromCtx(ctx)
					if chatID == 0 {
						log.Print("chat id not found in context")
					}

					macAddresses, err := s.bs.Search(ctx)
					if err != nil {
						s.textChatIdCh <- model.TextChatID{
							Text:   fmt.Sprintf("search: %s", err),
							ChatID: chatID,
						}
					}

					if len(macAddresses) == 0 {
						return
					}

					var buttons = make([]model.Button, 0, len(macAddresses))
					for _, ma := range macAddresses {
						buttons = append(buttons, model.Button{Key: ma})
					}

					s.keyboardCh <- model.Keyboard{
						ChatID:  chatID,
						Message: "choose the device",
						Buttons: buttons,
					}

					return
				},
				"leave_choice_waiting": func(ctx context.Context, e *fsm.Event) {
					chatID := helper.ChatIdFromCtx(ctx)
					if chatID == 0 {
						log.Print("chat id not found in context")
					}

					var text string
					defer func() {
						s.textChatIdCh <- model.TextChatID{
							Text:   text,
							ChatID: chatID,
						}
					}()

					textInterface, exists := e.FSM.Metadata(textKey)
					if !exists {
						text = fmt.Sprintf("metadata '%s' not found", textKey)
						return
					}

					t, ok := textInterface.(string)
					if !ok {
						text = "metadata not string"
					}

					text = t

					return
				},
			},
		)
	}
	return nil
}
