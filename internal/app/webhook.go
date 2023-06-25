package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	webhookName      = "webhook"
	setWebhookFormat = "https://api.telegram.org/bot%s/setWebhook?url=%s/api/v1/update"
)

type (
	webhook struct {
		description string
	}

	webhookResp struct {
		Ok          bool   `json:"ok"`
		Result      bool   `json:"result"`
		Description string `json:"description"`
	}
)

func newWebhook() *webhook {
	return &webhook{}
}

func (*webhook) Name() string {
	return webhookName
}

func (w *webhook) Init(ctx context.Context, mediator *mediator) error {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	url := fmt.Sprintf(setWebhookFormat, token, mediator.tunnel.URL())
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	var wr webhookResp
	err = json.NewDecoder(resp.Body).Decode(&wr)
	if err != nil {
		return fmt.Errorf("decode webhook response: %w", err)
	}

	failed := !(wr.Ok && wr.Result)
	if failed {
		return fmt.Errorf("webhook response description: %s", wr.Description)
	}

	w.description = strings.ToLower(wr.Description)

	return nil
}

func (w *webhook) SuccessLog() {
	log.Println(w.description)
}

func (*webhook) Close(ctx context.Context) error {
	return nil
}

func (*webhook) CloseLog() {
	closeLog(webhookName)
}
