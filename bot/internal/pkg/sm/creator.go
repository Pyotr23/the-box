package sm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Pyotr23/the-box/bot/internal/pkg/helper"
	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	"github.com/looplab/fsm"
)

const (
	blinkCommand      botCommand = "blink"
	searchCommand     botCommand = "search"
	registerCommand   botCommand = "register"
	unregisterCommand botCommand = "unregister"

	leavePrefix = "leave_"
	enterPrefix = "enter_"

	eventKey = "event"

	finishState = "finish"
)

type botCommand string

type bluetoothService interface {
	Search(ctx context.Context) ([]string, error)
	Blink(ctx context.Context, addr string) error
	RegisterDevice(ctx context.Context, name, macAddress string) error
	UnregisterDevice(ctx context.Context, id int) error
	DevicesMap(ctx context.Context) (map[string]model.Device, error)
	RegisteredDevicesMap(ctx context.Context) (map[string]model.Device, error)
	GetDeviceAliases(ctx context.Context) ([]string, error)
}

type fsmCreator struct {
	service        bluetoothService
	textChatIdCh   chan<- model.TextChatID
	keyboardCh     chan<- model.Keyboard
	smByBotCommand map[botCommand]func() *fsm.FSM
}

func newFsmCreator(
	service bluetoothService,
	textChatIdCh chan<- model.TextChatID,
	keyboarCh chan<- model.Keyboard,
) *fsmCreator {
	c := &fsmCreator{
		service:      service,
		textChatIdCh: textChatIdCh,
		keyboardCh:   keyboarCh,
	}
	c.smByBotCommand = map[botCommand]func() *fsm.FSM{
		searchCommand:     c.newSearchFSM,
		blinkCommand:      c.newBlinkFSM,
		registerCommand:   c.newRegisterDeviceFSM,
		unregisterCommand: c.newUnregisterDeviceFSM,
	}
	return c
}

func (c *fsmCreator) create(command string) (*fsm.FSM, error) {
	if len(command) < 2 {
		return nil, fmt.Errorf("too short command '%s'", command)
	}

	if create, ok := c.smByBotCommand[botCommand(command[1:])]; ok {
		return create(), nil
	}

	return nil, fmt.Errorf("unknown command '%s'", command)
}

func (c *fsmCreator) newSearchFSM() *fsm.FSM {
	const (
		startState  = "start"
		searchEvent = "search"
	)
	var sm = fsm.NewFSM(startState,
		fsm.Events{
			{
				Name: searchEvent,
				Src:  []string{startState},
				Dst:  finishState,
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

				aliases, err := c.service.GetDeviceAliases(ctx)
				if err != nil {
					err = fmt.Errorf("get devices aliases: %w", err)
					return
				}

				c.textChatIdCh <- model.TextChatID{
					Text:   strings.Join(aliases, ","),
					ChatID: chatID,
				}

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)

	return sm
}

func (c *fsmCreator) newRegisterDeviceFSM() *fsm.FSM {
	const (
		startState         = "start"
		choiceWaitingState = "choice_waiting"
		nameWaitingState   = "name_wating"
		searchEvent        = "search"
		getNameEvent       = "get_name"
		registerEvent      = "choice"
		macKey             = "address"
	)
	var sm = fsm.NewFSM("start",
		fsm.Events{
			{
				Name: searchEvent,
				Src:  []string{startState},
				Dst:  choiceWaitingState,
			},
			{
				Name: getNameEvent,
				Src:  []string{choiceWaitingState},
				Dst:  nameWaitingState,
			},
			{
				Name: registerEvent,
				Src:  []string{nameWaitingState},
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

				deviceByAddress, err := c.service.DevicesMap(ctx)
				if err != nil {
					err = fmt.Errorf("devices map: %w", err)
					return
				}

				if len(deviceByAddress) == 0 {
					return
				}

				var buttons = make([]model.Button, 0, len(deviceByAddress))
				for addr, d := range deviceByAddress {
					var key string
					if d.ID > 0 {
						key = d.Name
					} else {
						key = d.MacAddress
					}
					buttons = append(buttons,
						model.Button{
							Key:   key,
							Value: addr,
						},
					)
				}

				c.keyboardCh <- model.Keyboard{
					ChatID:  chatID,
					Message: "choose the device:",
					Buttons: buttons,
				}

				e.FSM.SetMetadata(eventKey, getNameEvent)
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

				e.FSM.SetMetadata(macKey, userChoice)

				c.textChatIdCh <- model.TextChatID{
					ChatID: chatID,
					Text:   "enter device name",
				}

				e.FSM.SetMetadata(eventKey, registerEvent)

				return
			},
			withLeavePrefix(nameWaitingState): func(ctx context.Context, e *fsm.Event) {
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

				deviceName, ok := e.Args[0].(string)
				if !ok {
					err = errors.New("first arg not string")
					return
				}

				macValue, ok := e.FSM.Metadata(macKey)
				if !ok {
					err = fmt.Errorf("metadata '%s' not found", macKey)
					return
				}
				address, ok := macValue.(string)
				if !ok {
					err = fmt.Errorf("metadat '%s' not string", macKey)
					return
				}

				if err := c.service.RegisterDevice(ctx, deviceName, address); err != nil {
					err = fmt.Errorf("register device: %w", err)
					return
				}

				c.textChatIdCh <- model.TextChatID{
					ChatID: chatID,
					Text:   "device registered",
				}

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)

	return sm
}

func (c *fsmCreator) newUnregisterDeviceFSM() *fsm.FSM {
	const (
		startState         = "start"
		choiceWaitingState = "choice_waiting"
		nameWaitingState   = "name_wating"

		searchEvent     = "search"
		unregisterEvent = "unregister"
		notFoundEvent   = "not_found"

		macKey = "address"
	)
	var sm = fsm.NewFSM(startState,
		fsm.Events{
			{
				Name: searchEvent,
				Src:  []string{startState},
				Dst:  choiceWaitingState,
			},
			{
				Name: unregisterEvent,
				Src:  []string{choiceWaitingState},
				Dst:  finishState,
			},
			{
				Name: notFoundEvent,
				Src:  []string{startState},
				Dst:  finishState,
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

				deviceByAddress, err := c.service.RegisteredDevicesMap(ctx)
				if err != nil {
					err = fmt.Errorf("devices map: %w", err)
					return
				}

				if len(deviceByAddress) == 0 {
					c.textChatIdCh <- model.TextChatID{
						ChatID: chatID,
						Text:   "registered devices not found",
					}

					e.FSM.SetMetadata(eventKey, notFoundEvent)

					return
				}

				var buttons = make([]model.Button, 0, len(deviceByAddress))
				for _, d := range deviceByAddress {
					buttons = append(buttons,
						model.Button{
							Key:   d.Name,
							Value: strconv.Itoa(d.ID),
						},
					)
				}

				c.keyboardCh <- model.Keyboard{
					ChatID:  chatID,
					Message: "choose the device:",
					Buttons: buttons,
				}

				e.FSM.SetMetadata(eventKey, unregisterEvent)
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

				id, err := strconv.Atoi(userChoice)
				if err != nil {
					err = errors.New("user choice not integer")
					return
				}

				if err = c.service.UnregisterDevice(ctx, id); err != nil {
					err = fmt.Errorf("unregister device: %w", err)
					return
				}

				c.textChatIdCh <- model.TextChatID{
					ChatID: chatID,
					Text:   "device unregistered",
				}

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)

	return sm
}

func (c *fsmCreator) newBlinkFSM() *fsm.FSM {
	const (
		startState         = "start"
		choiceWaitingState = "choice_waiting"
		searchEvent        = "search"
		choiceEvent        = "choice"
	)
	var sm = fsm.NewFSM("start",
		fsm.Events{
			{
				Name: searchEvent,
				Src:  []string{startState},
				Dst:  choiceWaitingState,
			},
			{
				Name: choiceEvent,
				Src:  []string{choiceWaitingState},
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

				deviceByAddress, err := c.service.DevicesMap(ctx)
				if err != nil {
					err = fmt.Errorf("devices map: %w", err)
					return
				}

				if len(deviceByAddress) == 0 {
					return
				}

				var buttons = make([]model.Button, 0, len(deviceByAddress))
				for addr, d := range deviceByAddress {
					var key string
					if d.ID > 0 {
						key = d.Name
					} else {
						key = d.MacAddress
					}
					buttons = append(buttons,
						model.Button{
							Key:   key,
							Value: addr,
						},
					)
				}

				c.keyboardCh <- model.Keyboard{
					ChatID:  chatID,
					Message: "choose the device:",
					Buttons: buttons,
				}

				e.FSM.SetMetadata(eventKey, choiceEvent)
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

				if err = c.service.Blink(ctx, userChoice); err != nil {
					err = fmt.Errorf("blink: %w", err)
					return
				}

				c.textChatIdCh <- model.TextChatID{
					ChatID: chatID,
					Text:   "finish blinking",
				}

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)

	return sm
}

func withLeavePrefix(state string) string {
	return leavePrefix + state
}

func withEnterPrefix(state string) string {
	return enterPrefix + state
}
