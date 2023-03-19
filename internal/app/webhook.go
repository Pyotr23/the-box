package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Pyotr23/the-box/internal/model"
)

const setWebhookFormat = "https://api.telegram.org/bot%s/setWebhook?url=%s/api/v1/update"

type webhook struct {
	description string
}

func (w *webhook) Init(ctx context.Context, a *App) error {
	token := os.Getenv(botTokenEnv)
	if token == "" {
		return fmt.Errorf("empty bot token environment %s", botTokenEnv)
	}

	url := fmt.Sprintf(setWebhookFormat, token, a.tunnel.URL())
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	var wr model.WebhookResp
	err = json.NewDecoder(resp.Body).Decode(&wr)
	if err != nil {
		return fmt.Errorf("decode webhook response: %w", err)
	}

	if !(wr.Ok && wr.Result) {
		return fmt.Errorf("webhook response description: %s", wr.Description)
	}

	w.description = strings.ToLower(wr.Description)

	return nil
}

func (w *webhook) SuccessLog() {
	log.Println(w.description)
}
