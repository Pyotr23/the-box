package model

type (
	DeviceSection struct {
		ID int `ini:"id"`
	}
	SettingsInfo struct {
		DeviceSection `ini:"device"`
	}
)
