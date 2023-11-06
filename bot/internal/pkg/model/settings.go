package model

type (
	DeviceSection struct {
		ID int `ini:"id"`
	}
	SettingsInfo struct {
		Device DeviceSection `ini:"device"`
	}
)
