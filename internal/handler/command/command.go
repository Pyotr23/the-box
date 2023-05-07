package command

import (
	"fmt"
	"strconv"

	"github.com/Pyotr23/the-box/internal/enum"
	base "github.com/Pyotr23/the-box/internal/handler"
	"github.com/Pyotr23/the-box/internal/handler/model"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

type Command struct {
	base   base.BaseHandler
	code   enum.Code
	socket rfcomm.Socket
}

func NewCommand(c model.Command) Command {
	return Command{
		base:   base.NewBaseHandler(c.ChatID, c.OutputTextCh),
		code:   c.Code,
		socket: c.Socket,
	}
}

func (c Command) Handle() {
	err := c.socket.Command(c.code)
	if err != nil {
		c.base.ProcessError(fmt.Errorf("command handle: %w", err))
	}
}

type CallbackCommand struct {
	base        base.BaseHandler
	code        enum.Code
	socket      rfcomm.Socket
	inputCh     <-chan string
	waitInputCh chan struct{}
}

func NewCallbackCommand(c model.Command, inputCh <-chan string, waitInputCh chan struct{}) CallbackCommand {
	return CallbackCommand{
		base:        base.NewBaseHandler(c.ChatID, c.OutputTextCh),
		code:        c.Code,
		socket:      c.Socket,
		inputCh:     inputCh,
		waitInputCh: waitInputCh,
	}
}

func (c CallbackCommand) Handle() {
	go func() {
		var err error
		defer func() {
			if err != nil {
				c.base.ProcessError(fmt.Errorf("callback command handle: %w", err))
			}
		}()

		c.waitInputCh <- struct{}{}

		n, err := strconv.ParseUint(<-c.inputCh, 10, 8)
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
