package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	webhookName      = "webhook"
	botTokenEnv      = "THEBOX_BOTTOKEN"
	setWebhookFormat = "https://api.telegram.org/bot%s/setWebhook?url=%s/api/v1/update"
)

type tunnelGetter interface {
	getTunnel() net.Listener
}

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

func (w *webhook) Init(ctx context.Context, app interface{}) error {
	tg, ok := app.(tunnelGetter)
	if !ok {
		return errors.New("app not implements tunnel getter")
	}

	token, err := getToken()
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	url := fmt.Sprintf(setWebhookFormat, token, tg.getTunnel().Addr().String())
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
	log.Print(w.description)
}

func (*webhook) Close(ctx context.Context) error {
	return nil
}

func (*webhook) CloseLog() {
	closeLog(webhookName)
}

func getToken() (string, error) {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return "", fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}
	return token, nil
}
