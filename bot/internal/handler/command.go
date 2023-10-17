package handler

import (
	"fmt"

	"github.com/Pyotr23/the-box/bot/internal/enum"
	"github.com/Pyotr23/the-box/bot/internal/model"
	"github.com/Pyotr23/the-box/bot/internal/rfcomm"
)

type Command struct {
	base   baseHandler
	code   enum.Code
	socket rfcomm.Socket
}

func NewCommand(c model.Info) Command {
	return Command{
		// base:   newBaseHandler(c.ChatID, c.OutputTextCh),
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
