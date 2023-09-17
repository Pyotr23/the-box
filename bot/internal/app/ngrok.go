package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

const ngrokTunnelName = "ngrok tunnel"

type tunnelSetter interface {
	setTunnel(tunnel net.Listener)
}

type ngrokTunnel struct {
	tunnel ngrok.Tunnel
}

func newNgrokTunnel() *ngrokTunnel {
	return &ngrokTunnel{}
}

func (*ngrokTunnel) Name() string {
	return ngrokTunnelName
}

func (nt *ngrokTunnel) Init(ctx context.Context, app interface{}) error {
	ts, ok := app.(tunnelSetter)
	if !ok {
		return errors.New("app not implements tunnel setter")
	}

	var err error
	nt.tunnel, err = ngrok.Listen(context.Background(),
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return fmt.Errorf("ngrok listen: %w", err)
	}

	ts.setTunnel(nt.tunnel)

	return nil
}

func (nt *ngrokTunnel) SuccessLog() {
	log.Printf("tunnel URL %s", nt.tunnel.URL())
}

func (nt *ngrokTunnel) Close(ctx context.Context) error {
	return nt.tunnel.CloseWithContext(ctx)
}

func (*ngrokTunnel) CloseLog() {
	closeLog(ngrokTunnelName)
}
