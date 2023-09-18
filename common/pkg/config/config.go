package config

import (
	"fmt"
	"os"

	"github.com/Pyotr23/the-box/common/pkg/model"
	yaml "gopkg.in/yaml.v3"
)

func GetBluetoothApiPort() (int, error) {
	bs, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return 0, fmt.Errorf("read file: %w", err)
	}

	var cfg model.Config
	if err = yaml.Unmarshal(bs, &cfg); err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	return cfg.BluetoothApiCOnfig.Port, nil
}
