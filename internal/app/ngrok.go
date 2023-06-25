package app

import (
	"context"
	"fmt"

	"github.com/Pyotr23/the-box/internal/helper"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

const ngrokTunnelName = "ngrok tunnel"

type ngrokTunnel struct {
	tunnel ngrok.Tunnel
}

func newNgrokTunnel() *ngrokTunnel {
	return &ngrokTunnel{}
}

func (*ngrokTunnel) Name() string {
	return ngrokTunnelName
}

func (nt *ngrokTunnel) Init(ctx context.Context, mediator *mediator) error {
	var err error
	nt.tunnel, err = ngrok.Listen(context.Background(),
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return fmt.Errorf("ngrok listen: %w", err)
	}

	mediator.tunnel = nt.tunnel

	return nil
}

func (nt *ngrokTunnel) SuccessLog() {
	helper.Logln(fmt.Sprintf("tunnel URL %s", nt.tunnel.URL()))
}

func (nt *ngrokTunnel) Close(ctx context.Context) error {
	return nt.tunnel.CloseWithContext(ctx)
}

func (*ngrokTunnel) CloseLog() {
	closeLog(ngrokTunnelName)
}
