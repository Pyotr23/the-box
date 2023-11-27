package enum

type Code byte

const (
	UnknownCode Code = iota
	TemperatureCode
	PinOnCode
	PinOffCode
	CheckPinCode
)
