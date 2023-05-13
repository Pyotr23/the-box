package handler

import (
	"fmt"

	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/handler/model"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

type QueryHandler struct {
	base   baseHandler
	code   enum.Code
	socket rfcomm.Socket
}

func NewQueryHandler(c model.Info) QueryHandler {
	return QueryHandler{
		base:   newBaseHandler(c.ChatID, c.OutputTextCh),
		code:   c.Code,
		socket: c.Socket,
	}
}

func (h QueryHandler) Handle() {
	var err error
	defer func() {
		if err != nil {
			h.base.ProcessError(fmt.Errorf("query handle: %w", err))
		}
	}()

	answer, err := h.socket.Query(h.code)
	if err != nil {
		err = fmt.Errorf("query: %w", err)
		return
	}

	h.base.SendText(answer)
}
