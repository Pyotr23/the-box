// package query

// import (
// 	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
// 	"github.com/Pyotr23/the-box/internal/model"
// 	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// type TemperatureQuery struct {
// 	BaseQuery
// }

// func NewTemperatureQuery(socket rfcomm.Socket, botApi *tgapi.BotAPI, chatID int64) Query {
// 	return TemperatureQuery{
// 		BaseQuery: BaseQuery{
// 			Socket: socket,
// 			BotApi: botApi,
// 			ChatID: chatID,
// 		},
// 	}
// }

// func (tq TemperatureQuery) ReadFromHardware(code model.Code) (string, error) {
// 	return tq.Socket.Write(string(model.Temperature))
// }

// func (tq TemperatureQuery) SendToBot(answer string) error {
// 	message := tgapi.NewMessage(tq.ChatID, answer)
// 	_, err := tq.BotApi.Send(message)
// 	return err
// }
