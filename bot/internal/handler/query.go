package handler

import (
	"fmt"
	"log"

	"github.com/Pyotr23/the-box/bot/internal/enum"
	"github.com/Pyotr23/the-box/bot/internal/model"
	"github.com/Pyotr23/the-box/bot/internal/rfcomm"
)

type QueryHandler struct {
	base   baseHandler
	code   enum.Code
	socket rfcomm.Socket
}

func NewQueryHandler(c model.Info) QueryHandler {
	return QueryHandler{
		// base:   newBaseHandler(c.ChatID, c.OutputTextCh),
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
	log.Printf("query code %v\n", h.code)
	answer, err := h.socket.Query(h.code)
	if err != nil {
		err = fmt.Errorf("query: %w", err)
		return
	}
	log.Printf("answer '%s'\n", answer)
	h.base.SendText(answer)
}
