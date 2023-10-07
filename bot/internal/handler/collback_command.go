package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Pyotr23/the-box/bot/internal/enum"
	"github.com/Pyotr23/the-box/bot/internal/model"
	"github.com/Pyotr23/the-box/bot/internal/rfcomm"
)

const (
	timeForAnswerIsOut              = "time for answer is out"
	enterBluetoothDeviceID          = "enter bluetooth device id"
	enterLowerTemperatureThreshold  = "enter lower temperature threshold"
	enterHigherTemperatureThreshold = "enter higher temperature threshold"
	enterMode                       = "enter mode"
	secondsBeforeTimeout            = 5
)

var timeIsOutError = errors.New(timeForAnswerIsOut)

type baseCallbackCommand struct {
	baseHandler baseHandler
	code        enum.Code
	socket      rfcomm.Socket
	inputCh     <-chan string
	waitInputCh chan struct{}
	textForUser string
}

func NewBaseCallbackCommand(c model.Info, inputCh <-chan string, waitInputCh chan struct{}, textForUser string) baseCallbackCommand {
	return baseCallbackCommand{
		// baseHandler: newBaseHandler(c.ChatID, c.OutputTextCh),
		code:        c.Code,
		socket:      c.Socket,
		inputCh:     inputCh,
		waitInputCh: waitInputCh,
		textForUser: textForUser,
	}
}

func (c baseCallbackCommand) Handle() {
	go func() {
		var err error
		defer func() {
			if err != nil {
				if !errors.Is(err, timeIsOutError) {
					err = fmt.Errorf("callback command handle: %w", err)
				}
				c.baseHandler.ProcessError(err)
			}
		}()

		c.baseHandler.SendText(c.textForUser)

		go func() {
			c.waitInputCh <- struct{}{}
		}()

		input, err := getTextBeforeTimeout(c.inputCh, c.waitInputCh)
		if err != nil {
			return
		}

		n, err := strconv.ParseUint(input, 10, 8)
		if err != nil {
			err = fmt.Errorf("parse uint: %w", err)
			return
		}
		log.Printf("parse uint %d\n", n)
		err = c.socket.SendText(c.code, []byte{byte(n)})
		if err != nil {
			err = fmt.Errorf("send text %d: %w", byte(n), err)
			return
		}
	}()
}

func getTextBeforeTimeout(inputCh <-chan string, waitInputCh <-chan struct{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), secondsBeforeTimeout*time.Second)
	defer cancel()

	var input string
	select {
	case input = <-inputCh:
		return input, nil
	case <-ctx.Done():
		<-waitInputCh
		return "", timeIsOutError
	}
}

type SetIDCallbackCommand struct {
	baseCallbackCommand
}

func NewSetIDCallbackCommand(c model.Info, inputCh <-chan string, waitInputCh chan struct{}) SetIDCallbackCommand {
	return SetIDCallbackCommand{
		baseCallbackCommand: NewBaseCallbackCommand(c, inputCh, waitInputCh, enterBluetoothDeviceID),
	}
}

type SetModeCallbackCommand struct {
	baseCallbackCommand
}

func NewSetModeCallbackCommand(c model.Info, inputCh <-chan string, waitInputCh chan struct{}) SetModeCallbackCommand {
	return SetModeCallbackCommand{
		baseCallbackCommand: NewBaseCallbackCommand(c, inputCh, waitInputCh, enterMode),
	}
}

type SetLowerTemperatureThresholdCallbackCommand struct {
	baseCallbackCommand
}

func NewSetLowerTemperatureThresholdCallbackCommand(c model.Info,
	inputCh <-chan string,
	waitInputCh chan struct{},
) SetLowerTemperatureThresholdCallbackCommand {
	return SetLowerTemperatureThresholdCallbackCommand{
		baseCallbackCommand: NewBaseCallbackCommand(c, inputCh, waitInputCh, enterLowerTemperatureThreshold),
	}
}

type SetHigherTemperatureThresholdCallbackCommand struct {
	baseCallbackCommand
}

func NewSetHigherTemperatureThresholdCallbackCommand(c model.Info,
	inputCh <-chan string,
	waitInputCh chan struct{},
) SetHigherTemperatureThresholdCallbackCommand {
	return SetHigherTemperatureThresholdCallbackCommand{
		baseCallbackCommand: NewBaseCallbackCommand(c, inputCh, waitInputCh, enterHigherTemperatureThreshold),
	}
}
