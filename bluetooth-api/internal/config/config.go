package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pyotr23/the-box/bluetooth-api/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

func GetPort() (int, error) {
	path, _ := filepath.Abs("./config/config.yaml")
	bs, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read file: %w", err)
	}

	var cfg model.Config
	if err = yaml.Unmarshal(bs, &cfg); err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}
	return cfg.Port, nil
}
