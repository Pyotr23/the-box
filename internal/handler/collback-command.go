package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/handler/model"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

const (
	timeForAnswerIsOut     = "time for answer is out"
	enterBluetoothDeviceID = "enter bluetooth device id"
	secondsBeforeTimeout   = 5
)

var timeIsOutError = errors.New(timeForAnswerIsOut)

type callbackCommand struct {
	base        baseHandler
	code        enum.Code
	socket      rfcomm.Socket
	inputCh     <-chan string
	waitInputCh chan struct{}
}

type SetIDCallbackCommand struct {
	callbackCommand
}

func newCallbackCommand(c model.Info, inputCh <-chan string, waitInputCh chan struct{}) callbackCommand {
	return callbackCommand{
		base:        newBaseHandler(c.ChatID, c.OutputTextCh),
		code:        c.Code,
		socket:      c.Socket,
		inputCh:     inputCh,
		waitInputCh: waitInputCh,
	}
}

func NewSetIDCallbackCommand(c model.Info, inputCh <-chan string, waitInputCh chan struct{}) SetIDCallbackCommand {
	return SetIDCallbackCommand{
		callbackCommand: newCallbackCommand(c, inputCh, waitInputCh),
	}
}

func (c SetIDCallbackCommand) Handle() {
	go func() {
		var err error
		defer func() {
			if err != nil {
				if !errors.Is(err, timeIsOutError) {
					err = fmt.Errorf("callback command handle: %w", err)
				}
				c.base.ProcessError(err)
			}
		}()

		c.base.SendText(enterBluetoothDeviceID)

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
