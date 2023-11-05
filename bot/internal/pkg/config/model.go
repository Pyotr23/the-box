package config

type (
	Device struct {
		ID int `ini:"id"`
	}
	SettingsInfo struct {
		Device `ini:"device"`
	}
)
