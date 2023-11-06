package settings

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	"gopkg.in/ini.v1"
)

const (
	relativePath  = "./bot/configs/settings.ini"
	deviceSection = "device"
	idKey         = "id"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (_ Service) ReadConfig() (model.SettingsInfo, error) {
	absPath, err := getAbsIniPath()
	if err != nil {
		return model.SettingsInfo{}, fmt.Errorf("get abs ini path: %w", err)
	}

	var cfg = new(model.SettingsInfo)
	if err = ini.MapTo(cfg, absPath); err != nil {
		return model.SettingsInfo{}, fmt.Errorf("map to: %w", err)
	}
	if cfg == nil {
		return model.SettingsInfo{}, errors.New("nil config after mapping")
	}
	return *cfg, nil
}

func (_ Service) ReadDeviceID() (int, error) {
	absPath, err := getAbsIniPath()
	if err != nil {
		return 0, fmt.Errorf("get abs ini path: %w", err)
	}

	var cfg = new(model.SettingsInfo)
	if err = ini.MapTo(cfg, absPath); err != nil {
		return 0, fmt.Errorf("map to: %w", err)
	}
	if cfg == nil {
		return 0, errors.New("nil config after mapping")
	}
	return cfg.Device.ID, nil
}

func (_ Service) WriteDeviceID(id int) error {
	absPath, err := getAbsIniPath()
	if err != nil {
		return fmt.Errorf("get abs ini path: %w", err)
	}

	file, err := ini.Load(absPath)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	file.Section(deviceSection).Key(idKey).SetValue(strconv.Itoa(id))

	if err = file.SaveTo(absPath); err != nil {
		return fmt.Errorf("save to: %w", err)
	}

	return nil
}

func getAbsIniPath() (string, error) {
	return filepath.Abs(relativePath)
}
