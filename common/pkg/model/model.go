package model

type Config struct {
	BluetoothApiCOnfig BluetoothApiConfig `yaml:"bluetooth-api"`
}

type BluetoothApiConfig struct {
	Port int `yaml:"port"`
}
