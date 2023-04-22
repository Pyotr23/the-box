package enum

type Code byte

const (
	UnknownCode Code = iota
	TemperatureCode
	RelayOnCode
	RelayOffCode
)

type BotCommand string

const (
	Temperature BotCommand = "/temperature"
	RelayOn     BotCommand = "/relay-on"
	RelayOff    BotCommand = "/relay-off"
)

var codeByBotCommand = map[BotCommand]Code{
	Temperature: TemperatureCode,
	RelayOn:     RelayOnCode,
	RelayOff:    RelayOffCode,
}

func GetCode(c string) Code {
	return codeByBotCommand[BotCommand(c)]
}
