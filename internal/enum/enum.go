package enum

type Code byte

const (
	UnknownCode Code = iota
	TemperatureCode
	RelayOnCode
	RelayOffCode
	SetIDCode
	GetIDCode
	GetLowerTemperatureThresholdCode
	GetHigherTemperatureThresholdCode
	SetLowerTemperatureThresholdCode
	SetHigherTemperatureThresholdCode
)

type BotCommand string

const (
	Temperature                   BotCommand = "/temperature"
	RelayOn                       BotCommand = "/relay_on"
	RelayOff                      BotCommand = "/relay_off"
	SetID                         BotCommand = "/set_id"
	GetID                         BotCommand = "/get_id"
	GetLowerTemperatureThreshold  BotCommand = "/get_low_temp_thrld"
	GetHigherTemperatureThreshold BotCommand = "/get_high_temp_thrld"
	SetLowerTemperatureThreshold  BotCommand = "/set_low_temp_thrld"
	SetHigherTemperatureThreshold BotCommand = "/set_high_temp_thrld"
)

var codeByBotCommand = map[BotCommand]Code{
	Temperature:                   TemperatureCode,
	RelayOn:                       RelayOnCode,
	RelayOff:                      RelayOffCode,
	SetID:                         SetIDCode,
	GetID:                         GetIDCode,
	GetLowerTemperatureThreshold:  GetLowerTemperatureThresholdCode,
	GetHigherTemperatureThreshold: GetHigherTemperatureThresholdCode,
	SetLowerTemperatureThreshold:  SetLowerTemperatureThresholdCode,
	SetHigherTemperatureThreshold: SetHigherTemperatureThresholdCode,
}

func GetCode(c string) Code {
	return codeByBotCommand[BotCommand(c)]
}
