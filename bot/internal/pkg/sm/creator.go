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
	blinkCommand botCommand = "blink"

	searchCommand botCommand = "search"

	registerCommand   botCommand = "register"
	unregisterCommand botCommand = "unregister"

	setDeviceCommand    botCommand = "device_set"
	activeDeviceCommand botCommand = "active_device"

	leavePrefix = "leave_"
	enterPrefix = "enter_"

	eventKey  = "event"
	chatIdKey = "chatID"

	startState  = "start"
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
	GetDeviceFullInfo(ctx context.Context, id int) (model.DeviceInfo, error)
}

type settingsService interface {
	WriteDeviceID(id int) error
	ReadDeviceID() (int, error)
}

type fsmCreator struct {
	service         bluetoothService
	settingsService settingsService
	textChatIdCh    chan<- model.TextChatID
	keyboardCh      chan<- model.Keyboard
	smByBotCommand  map[botCommand]func(chatID int64) *fsm.FSM
	chatID          int64
}

func newFsmCreator(
	service bluetoothService,
	settingsWriter settingsService,
	textChatIdCh chan<- model.TextChatID,
	keyboarCh chan<- model.Keyboard,
) *fsmCreator {
	c := &fsmCreator{
		service:         service,
		settingsService: settingsWriter,
		textChatIdCh:    textChatIdCh,
		keyboardCh:      keyboarCh,
	}
	c.smByBotCommand = map[botCommand]func(chatID int64) *fsm.FSM{
		searchCommand:       c.newSearchFSM,
		blinkCommand:        c.newBlinkFSM,
		registerCommand:     c.newRegisterDeviceFSM,
		unregisterCommand:   c.newUnregisterDeviceFSM,
		setDeviceCommand:    c.newSetDeviceFSM,
		activeDeviceCommand: c.newActiveDeviceFSM,
	}
	return c
}

func (c *fsmCreator) create(chatID int64, command string) (*fsm.FSM, error) {
	if len(command) < 2 {
		return nil, fmt.Errorf("too short command '%s'", command)
	}

	if create, ok := c.smByBotCommand[botCommand(command[1:])]; ok {
		return create(chatID), nil
	}

	return nil, fmt.Errorf("unknown command '%s'", command)
}

func (c *fsmCreator) newActiveDeviceFSM(chatID int64) *fsm.FSM {
	const (
		gotActiveDeviceEvent = "active_device_got"
	)
	var sm = fsm.NewFSM(startState,
		fsm.Events{
			{
				Name: gotActiveDeviceEvent,
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

				id, err := c.settingsService.ReadDeviceID()
				if err != nil {
					err = fmt.Errorf("read device id: %w", err)
					return
				}
				if id == 0 {
					err = errors.New("no active device id in settings")
					return
				}

				device, err := c.service.GetDeviceFullInfo(ctx, id)
				if err != nil {
					err = fmt.Errorf("get device full info: %w", err)
					return
				}

				c.sendText(e.FSM, fmt.Sprintf("%v", device))

				return
			},
		},
	)

	sm.SetMetadata(eventKey, gotActiveDeviceEvent)
	sm.SetMetadata(chatIdKey, chatID)

	return sm
}

func (c *fsmCreator) newSetDeviceFSM(chatID int64) *fsm.FSM {
	const (
		idWaitingState  = "id_waiting"
		gotAliasesEvent = "aliasese_got"
		idWaitedEvent   = "id_waited"
	)

	var sm = fsm.NewFSM(startState,
		fsm.Events{
			{
				Name: gotAliasesEvent,
				Src:  []string{startState},
				Dst:  idWaitingState,
			},
			{
				Name: idWaitedEvent,
				Src:  []string{idWaitingState},
				Dst:  finishState,
			},
		},
		fsm.Callbacks{
			withLeavePrefix(startState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

				deviceByAddress, err := c.service.DevicesMap(ctx)
				if err != nil {
					err = fmt.Errorf("devices map: %w", err)
					return
				}

				if len(deviceByAddress) == 0 {
					return
				}

				var buttons = make([]model.Button, 0, len(deviceByAddress))
				for _, d := range deviceByAddress {
					var (
						key   string
						value int
					)
					if d.ID > 0 {
						key = d.Name
						value = d.ID
					} else {
						key = d.MacAddress
					}
					buttons = append(buttons,
						model.Button{
							Key:   key,
							Value: strconv.Itoa(value),
						},
					)
				}

				c.sendButtons(e.FSM, "choose current device:", buttons)

				e.FSM.SetMetadata(eventKey, idWaitedEvent)

				return
			},
			withLeavePrefix(idWaitingState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

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

				if id == 0 {
					c.textChatIdCh <- model.TextChatID{
						Text:   fmt.Sprintf("device unregistered, use '/%s' command", registerCommand),
						ChatID: chatID,
					}
					return
				}

				if err := c.settingsService.WriteDeviceID(id); err != nil {
					err = fmt.Errorf("write device id: %w", err)
					return
				}

				c.sendText(e.FSM, "current device set")

				return
			},
		},
	)

	sm.SetMetadata(eventKey, gotAliasesEvent)
	sm.SetMetadata(chatIdKey, chatID)

	return sm
}

func (c *fsmCreator) newSearchFSM(chatID int64) *fsm.FSM {
	const (
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

				aliases, err := c.service.GetDeviceAliases(ctx)
				if err != nil {
					err = fmt.Errorf("get devices aliases: %w", err)
					return
				}

				c.sendText(e.FSM, strings.Join(aliases, ", "))

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)
	sm.SetMetadata(chatIdKey, chatID)

	return sm
}

func (c *fsmCreator) newRegisterDeviceFSM(chatID int64) *fsm.FSM {
	const (
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

				c.sendText(e.FSM, "enter device name:")

				e.FSM.SetMetadata(eventKey, registerEvent)

				return
			},
			withLeavePrefix(nameWaitingState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

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

				c.sendText(e.FSM, "device registered")

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)
	sm.SetMetadata(chatIdKey, chatID)

	return sm
}

func (c *fsmCreator) newUnregisterDeviceFSM(chatID int64) *fsm.FSM {
	const (
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

				deviceByAddress, err := c.service.RegisteredDevicesMap(ctx)
				if err != nil {
					err = fmt.Errorf("devices map: %w", err)
					return
				}

				if len(deviceByAddress) == 0 {
					c.sendText(e.FSM, "registered devices not found")

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

				c.sendButtons(e.FSM, "choose the device:", buttons)

				e.FSM.SetMetadata(eventKey, unregisterEvent)
			},
			withLeavePrefix(choiceWaitingState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

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

				c.sendText(e.FSM, "device unregistered")

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)
	sm.SetMetadata(chatIdKey, chatID)

	return sm
}

func (c *fsmCreator) newBlinkFSM(chatID int64) *fsm.FSM {
	const (
		choiceWaitingState = "choice_waiting"
		searchEvent        = "search"
		choiceEvent        = "choice"
	)
	var sm = fsm.NewFSM(startState,
		fsm.Events{
			{
				Name: searchEvent,
				Src:  []string{startState},
				Dst:  choiceWaitingState,
			},
			{
				Name: choiceEvent,
				Src:  []string{choiceWaitingState},
				Dst:  finishState,
			},
		},
		fsm.Callbacks{
			withLeavePrefix(startState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

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

				c.sendButtons(e.FSM, "choose the device:", buttons)

				e.FSM.SetMetadata(eventKey, choiceEvent)
			},
			withLeavePrefix(choiceWaitingState): func(ctx context.Context, e *fsm.Event) {
				var err error
				defer func() {
					e.Err = err
				}()

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

				c.sendText(e.FSM, "finish blinking")

				return
			},
		},
	)

	sm.SetMetadata(eventKey, searchEvent)
	sm.SetMetadata(chatIdKey, chatID)

	return sm
}

func getChatID(sm *fsm.FSM) (int64, error) {
	genericChatID, ok := sm.Metadata(chatIdKey)
	if !ok {
		return 0, errors.New("no chat id in metadata")
	}
	chatID, ok := genericChatID.(int64)
	if !ok {
		return 0, errors.New("chat id in metadat not int")
	}
	return chatID, nil
}

func (c *fsmCreator) sendButtons(sm *fsm.FSM, msg string, buttons []model.Button) {
	chatID, err := getChatID(sm)
	if err != nil {
		log.Print("no chat id in state machine metadata")
		return
	}
	c.keyboardCh <- model.Keyboard{
		ChatID:  chatID,
		Message: msg,
		Buttons: buttons,
	}
}

func (c *fsmCreator) sendText(sm *fsm.FSM, msg string) {
	chatID, err := getChatID(sm)
	if err != nil {
		log.Print("no chat id in state machine metadata")
		return
	}
	c.textChatIdCh <- model.TextChatID{
		ChatID: chatID,
		Text:   msg,
	}
}

func withLeavePrefix(state string) string {
	return leavePrefix + state
}

func withEnterPrefix(state string) string {
	return enterPrefix + state
}
