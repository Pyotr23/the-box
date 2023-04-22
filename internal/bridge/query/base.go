// package query

// import (
// 	"github.com/Pyotr23/the-box/internal/hardware/rfcomm"
// 	"github.com/Pyotr23/the-box/internal/model"
// 	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// type Query interface {
// 	ReadFromHardware(code model.Code) (string, error)
// 	SendToBot(answer string) error
// }

// type BaseQuery struct {
// 	Socket rfcomm.Socket
// 	BotApi *tgapi.BotAPI
// 	ChatID int64
// }
