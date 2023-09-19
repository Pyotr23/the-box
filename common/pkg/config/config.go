package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Pyotr23/the-box/common/pkg/model"
	yaml "gopkg.in/yaml.v3"
)

var (
	cfg  model.Config
	path string
)

func init() {
	var err error
	path, err = filepath.Abs("./config/config.yaml")
	if err != nil {
		log.Printf("abs: %s", err)
		return
	}

	bs, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read file: %s", err)
		return
	}

	if err = yaml.Unmarshal(bs, &cfg); err != nil {
		log.Printf("unmarshal: %s", err)
		return
	}
}

func GetBluetoothApiPort() (int, error) {
	// path, err := filepath.Abs("./config/config.yaml")
	fmt.Println(path)
	// if err != nil {
	// 	return 0, fmt.Errorf("abs: %w", err)
	// }

	// bs, err := os.ReadFile(path)
	// if err != nil {
	// 	return 0, fmt.Errorf("read file: %w", err)
	// }

	// if err = yaml.Unmarshal(bs, &cfg); err != nil {
	// 	return 0, fmt.Errorf("unmarshal: %w", err)
	// }

	return 5001, nil
}
