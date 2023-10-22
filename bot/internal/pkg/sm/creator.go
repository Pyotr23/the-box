package sm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Pyotr23/the-box/bot/internal/pkg/helper"
	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	"github.com/looplab/fsm"
)

const (
	blinkCommand  string = "blink"
	searchCommand string = "search"

	leavePrefix = "leave_"
)

type bluetoothService interface {
	Search(ctx context.Context) ([]string, error)
}

type fsmCreator struct {
	service      bluetoothService
	textChatIdCh chan<- model.TextChatID
	keyboardCh   chan<- model.Keyboard
}

func newFsmCreator(
	service bluetoothService,
	textChatIdCh chan<- model.TextChatID,
	keyboarCh chan<- model.Keyboard,
) *fsmCreator {
	return &fsmCreator{
		service:      service,
		textChatIdCh: textChatIdCh,
		keyboardCh:   keyboarCh,
	}
}

func (c *fsmCreator) create(command string) (*fsm.FSM, error) {
	if len(command) < 2 {
		return nil, fmt.Errorf("too short command '%s'", command)
	}

	switch command[1:] {
	case searchCommand:
		return c.newSearchFSM(), nil
	case blinkCommand:
		return c.newBlinkFSM(), nil
	default:
		return nil, fmt.Errorf("unknown command '%s'", command)
	}
}

func (c *fsmCreator) newSearchFSM() *fsm.FSM {
	const startState = "start"
	return fsm.NewFSM(startState,
		fsm.Events{
			{
				Name: "search",
				Src:  []string{startState},
				Dst:  "finish",
			},
		},
		fsm.Callbacks{
			withLeavePrefix(startState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

				chatID := helper.ChatIdFromCtx(ctx)
				if chatID == 0 {
					log.Print(model.ErrMessageNoChatID)
					return
				}

				macAddresses, err := c.service.Search(ctx)
				if err != nil {
					err = fmt.Errorf("search: %w", err)
					return
				}

				c.textChatIdCh <- model.TextChatID{
					Text:   strings.Join(macAddresses, ","),
					ChatID: chatID,
				}

				return
			},
		},
	)
}

func (c *fsmCreator) newBlinkFSM() *fsm.FSM {
	const (
		startState         = "start"
		choiceWaitingState = "choice_waiting"
	)
	return fsm.NewFSM("start",
		fsm.Events{
			{
				Name: "search",
				Src:  []string{startState},
				Dst:  choiceWaitingState,
			},
			{
				Name: "choice",
				Src:  []string{choiceWaitingState},
				Dst:  "blinking",
			},
			// {
			// 	Name: "blink",
			// 	Src:  []string{"blinking"},
			// 	Dst:  "finish",
			// },
		},
		fsm.Callbacks{
			withLeavePrefix(startState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

				chatID := helper.ChatIdFromCtx(ctx)
				if chatID == 0 {
					log.Print(model.ErrMessageNoChatID)
					return
				}

				macAddresses, err := c.service.Search(ctx)
				if err != nil {
					err = fmt.Errorf("search: %w", err)
					return
				}

				if len(macAddresses) == 0 {
					return
				}

				var buttons = make([]model.Button, 0, len(macAddresses))
				for _, ma := range macAddresses {
					buttons = append(buttons,
						model.Button{
							Key:   ma,
							Value: ma,
						},
					)
				}

				c.keyboardCh <- model.Keyboard{
					ChatID:  chatID,
					Message: "choose the device:",
					Buttons: buttons,
				}
			},
			withLeavePrefix(choiceWaitingState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

				chatID := helper.ChatIdFromCtx(ctx)
				if chatID == 0 {
					log.Print(model.ErrMessageNoChatID)
					return
				}

				if len(e.Args) == 0 {
					err = errors.New("no args")
					return
				}

				userChoice, ok := e.Args[0].(string)
				if !ok {
					err = errors.New("first arg not string")
					return
				}

				c.textChatIdCh <- model.TextChatID{
					Text:   userChoice,
					ChatID: chatID,
				}

				return
			},
		},
	)
}

func withLeavePrefix(state string) string {
	return leavePrefix + state
}
