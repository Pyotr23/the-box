package enum

type Code byte

const (
	UnknownCode Code = iota
	TemperatureCode
	RelayOnCode
	RelayOffCode
	SetIDCode
	GetIDCode
)

type BotCommand string

const (
	Temperature BotCommand = "/temperature"
	RelayOn     BotCommand = "/relay-on"
	RelayOff    BotCommand = "/relay-off"
	SetID       BotCommand = "/set-id"
	GetID       BotCommand = "/get-id"
)

var codeByBotCommand = map[BotCommand]Code{
	Temperature: TemperatureCode,
	RelayOn:     RelayOnCode,
	RelayOff:    RelayOffCode,
	SetID:       SetIDCode,
	GetID:       GetIDCode,
}

func GetCode(c string) Code {
	return codeByBotCommand[BotCommand(c)]
}
