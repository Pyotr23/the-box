// package bridge

// import (
// 	"log"

// 	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
// 	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// type Request string

// const (
// 	Temperature Request = "temperature"
// )

// type Reader interface {
// 	Read(chatID int64, request Request) error
// }

// type temperatureReader struct {
// 	socket rfcomm.Socket
// 	api    *tgapi.BotAPI
// }

// func NewTemperatureReader(socket rfcomm.Socket, botApi *tgapi.BotAPI) Reader {
// 	return temperatureReader{
// 		socket: socket,
// 		api:    botApi,
// 	}
// }

// func (r temperatureReader) Read(chatID int64, request Request) error {
// 	btAnswer, err := r.socket.Write("1")
// 	if err != nil {
// 		log.Printf("write: %s", err.Error())
// 		return err
// 	}

// 	message := tgapi.NewMessage(chatID, btAnswer)

// 	_, err = r.api.Send(message)
// 	if err != nil {
// 		log.Printf("send message: %s", err.Error())
// 	}

// 	return nil
// }

// type BotCommandHandler struct {
// 	socket rfcomm.Socket
// }

// func (h BotCommandHandler) Handle(bc string) {
// 	switch bc {
// 	case "temperature":

// 	}
// }
