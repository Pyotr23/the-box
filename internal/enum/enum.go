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
	SetModeCode
	GetModeCode
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
	SetMode                       BotCommand = "/set_mode"
	GetMode                       BotCommand = "/get_mode"
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
	SetMode:                       SetModeCode,
	GetMode:                       GetModeCode,
}

func GetCode(c string) Code {
	return codeByBotCommand[BotCommand(c)]
}
