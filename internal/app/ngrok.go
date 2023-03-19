package app

import (
	"context"
	"fmt"
	"log"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type ngrokTunnel struct {
	url string
}

func (n *ngrokTunnel) Init(ctx context.Context, a *App) error {
	tunnel, err := ngrok.Listen(context.Background(),
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return fmt.Errorf("ngrok listen: %w", err)
	}

	a.tunnel = tunnel
	n.url = tunnel.URL()

	return nil
}

func (n *ngrokTunnel) SuccessLog() {
	log.Printf("tunnel URL %s\n", n.url)
}
