package app

import (
	"context"
	"log"

	"github.com/Pyotr23/the-box/internal/enum"
	"github.com/Pyotr23/the-box/internal/handler"
	"github.com/Pyotr23/the-box/internal/handler/model"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

const messageName = "message"

type message struct {
	inlineTextCh    chan string
	inputMessageCh  chan model.Message
	outputMessageCh chan model.Message
	waitInputCh     chan struct{}
	socket          rfcomm.Socket
}

type botCommandHandler interface {
	Handle()
}

func newMessage() *message {
	return &message{
		waitInputCh:  make(chan struct{}),
		inlineTextCh: make(chan string),
	}
}

func (*message) Name() string {
	return messageName
}

func (m *message) Init(ctx context.Context, a *App) error {
	m.inputMessageCh = a.inputMessageCh
	m.outputMessageCh = a.outputMessageCh
	m.socket = a.sockets[0]

	go func() {
		for inputMessage := range m.inputMessageCh {
			m.handle(inputMessage)
		}
	}()

	return nil
}

func (*message) SuccessLog() {
	log.Println("setup message processor")
}

func (m *message) Close(ctx context.Context, a *App) error {
	close(m.waitInputCh)
	return nil
}

func (*message) CloseLog() {
	closeLog(messageName)
}

func (m *message) handle(message model.Message) {
	select {
	case <-m.waitInputCh:
		m.inlineTextCh <- message.Text
		return
	default:
		break
	}

	log.Printf("message text from bot '%s'\n", message.Text)

	h := m.getHandler(message)

	h.Handle()
}

func (m *message) createCommand(msg model.Message) model.Info {
	return model.Info{
		Code:         enum.GetCode(msg.Text),
		Socket:       m.socket,
		ChatID:       msg.ChatID,
		OutputTextCh: m.outputMessageCh,
	}
}

func (m *message) getHandler(msg model.Message) (h botCommandHandler) {
	info := m.createCommand(msg)

	switch info.Code {
	case enum.TemperatureCode:
		return handler.NewQueryHandler(info)
	case enum.RelayOnCode:
		return handler.NewCommand(info)
	case enum.RelayOffCode:
		return handler.NewCommand(info)
	case enum.SetIDCode:
		return handler.NewSetIDCallbackCommand(info, m.inlineTextCh, m.waitInputCh)
	case enum.GetIDCode:
		return handler.NewQueryHandler(info)
	case enum.UnknownCode:
		return handler.NewUnknownHandler(info)
	}

	log.Println("unexpected info code")

	return
}
