package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"gopkg.in/ini.v1"
)

const (
	relativePath  = "./bot/configs/settings.ini"
	deviceSection = "device"
	idKey         = "id"
)

func ReadConfig() (SettingsInfo, error) {
	absPath, err := getAbsIniPath()
	if err != nil {
		return SettingsInfo{}, fmt.Errorf("get abs ini path: %w", err)
	}

	var cfg = new(SettingsInfo)
	if err = ini.MapTo(cfg, absPath); err != nil {
		return SettingsInfo{}, fmt.Errorf("map to: %w", err)
	}
	if cfg == nil {
		return SettingsInfo{}, errors.New("nil config after mapping")
	}
	return *cfg, nil
}

func WriteDeviceID(id int) error {
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
