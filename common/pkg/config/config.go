package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pyotr23/the-box/common/pkg/model"
	yaml "gopkg.in/yaml.v3"
)

const configPath = "./common/config/config.yaml"

func GetBluetoothApiPort() (int, error) {
	path, err := filepath.Abs(configPath)
	if err != nil {
		return 0, fmt.Errorf("abs: %w", err)
	}

	bs, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read file: %w", err)
	}

	var cfg model.Config
	if err = yaml.Unmarshal(bs, &cfg); err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	return cfg.BluetoothApiConfig.Port, nil
}
