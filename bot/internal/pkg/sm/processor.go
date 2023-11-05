package sm

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pyotr23/the-box/bot/internal/pkg/helper"
	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	"github.com/looplab/fsm"
)

type creator interface {
	create(chatID int64, command string) (*fsm.FSM, error)
}

type fsmProcessor struct {
	fsmByChatID map[int64]*fsm.FSM
	creator     creator
}

func NewFsmProcessor(
	bs bluetoothService,
	sts settingsWriter,
	textChatIdCh chan<- model.TextChatID,
	keyboardCh chan<- model.Keyboard,
) *fsmProcessor {
	return &fsmProcessor{
		fsmByChatID: make(map[int64]*fsm.FSM),
		creator:     newFsmCreator(bs, sts, textChatIdCh, keyboardCh),
	}
}

func (p *fsmProcessor) GetCommandProcessor() func(ctx context.Context, command string) error {
	return func(ctx context.Context, command string) error {
		chatID := helper.ChatIdFromCtx(ctx)
		if chatID == 0 {
			return errors.New("chat id not found in context")
		}

		sm, err := p.creator.create(chatID, command)
		if err != nil {
			return fmt.Errorf("create fsm: %w", err)
		}

		p.fsmByChatID[chatID] = sm

		if err = makeEvent(ctx, sm); err != nil {
			return fmt.Errorf("make event: %w", err)
		}

		defer func() {
			if len(sm.AvailableTransitions()) == 0 {
				delete(p.fsmByChatID, chatID)
			}
		}()

		return nil
	}
}

func (p *fsmProcessor) GetTextProcessor() func(ctx context.Context, text string) error {
	return func(ctx context.Context, text string) error {
		chatID := helper.ChatIdFromCtx(ctx)
		if chatID == 0 {
			return errors.New("chat id not found in context")
		}

		sm := p.fsmByChatID[chatID]
		if sm == nil {
			return errors.New("text but no command before")
		}

		defer func() {
			if len(sm.AvailableTransitions()) == 0 {
				p.fsmByChatID[chatID] = nil
			}
		}()

		if err := makeEvent(ctx, sm, text); err != nil {
			return fmt.Errorf("make event: %w", err)
		}

		return nil
	}
}

func makeEvent(ctx context.Context, sm *fsm.FSM, args ...interface{}) error {
	startState := sm.Current()

	genericEvent, ok := sm.Metadata(eventKey)
	if !ok {
		return fmt.Errorf("event not found in metadata by key '%s'", eventKey)
	}

	event, ok := genericEvent.(string)
	if !ok {
		return errors.New("event not string in metadata")
	}
	if err := sm.Event(ctx, event, args...); err != nil {
		sm.SetState(startState)
		return fmt.Errorf("make event '%s': %w", event, err)
	}

	return nil
}
