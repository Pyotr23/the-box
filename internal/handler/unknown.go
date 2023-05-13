package handler

import (
	"errors"

	"github.com/Pyotr23/the-box/internal/handler/model"
)

type UnknownHandler struct {
	base baseHandler
}

func NewUnknownHandler(c model.Info) UnknownHandler {
	return UnknownHandler{
		base: newBaseHandler(c.ChatID, c.OutputTextCh),
	}
}

func (h UnknownHandler) Handle() {
	h.base.ProcessError(errors.New("unknown command"))
}
