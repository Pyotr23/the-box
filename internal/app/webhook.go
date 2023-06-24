package app

import (
	"context"
	"encoding/json"
	"errors"
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

func (w *webhook) Init(ctx context.Context, tunnelUrl interface{}) (interface{}, error) {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return nil, fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	tunnelUrl, ok := tunnelUrl.(string)
	if !ok {
		return nil, errors.New("tunnel url not a string type")
	}

	url := fmt.Sprintf(setWebhookFormat, token, tunnelUrl)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("set webhook: %w", err)
	}

	var wr webhookResp
	err = json.NewDecoder(resp.Body).Decode(&wr)
	if err != nil {
		return nil, fmt.Errorf("decode webhook response: %w", err)
	}

	if !(wr.Ok && wr.Result) {
		return nil, fmt.Errorf("webhook response description: %s", wr.Description)
	}

	w.description = strings.ToLower(wr.Description)

	return nil, nil
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
