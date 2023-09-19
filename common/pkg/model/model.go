package model

type Config struct {
	BluetoothApiConfig BluetoothApiConfig `yaml:"bluetooth-api"`
}

type BluetoothApiConfig struct {
	Port int `yaml:"port"`
}
