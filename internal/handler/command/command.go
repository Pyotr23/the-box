package command

import (
	"fmt"

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
