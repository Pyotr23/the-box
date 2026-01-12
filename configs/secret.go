package configs

type Secret struct {
	BotToken     string `yaml:"bot_token"`
	WebhookToken string `yaml:"webhook_token"`
}

