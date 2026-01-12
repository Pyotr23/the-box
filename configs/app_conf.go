package configs

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

const secretPath = "./configs/secret.yaml"

type AppConf struct {
	secret Secret
}

func (a *AppConf) GetBotToken() string {
	return a.secret.BotToken
}

func (a *AppConf) GetWebhookToken() string {
	return a.secret.WebhookToken
}

func NewAppConf() (*AppConf, error) {
	secretContent, err := os.ReadFile(secretPath)
	if err != nil {
		return nil, fmt.Errorf("read file by path '%s': %w", secretPath, err)
	}

	var secret Secret
	if err := yaml.Unmarshal(secretContent, &secret); err != nil {
		return nil, fmt.Errorf("unmarshal '%s': %w", string(secretContent), err)
	}

	return &AppConf{
		secret: secret,
	}, nil
}
