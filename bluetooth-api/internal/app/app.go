package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	macl "github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/client/mac_address"
	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/server"
	masrv "github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/service/mac_address"
	socketrv "github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/service/socket"
	common "github.com/Pyotr23/the-box/common/pkg/config"
	"google.golang.org/grpc"
)

type closer interface {
	Close(ctx context.Context)
}

type App struct {
	Listener net.Listener
	Server   *grpc.Server
	closers  []closer
}

func NewApp() (*App, error) {
	var (
		err error
		app = &App{}
	)

	if app.Listener, err = getListener(); err != nil {
		return nil, fmt.Errorf("get listener: %w", err)
	}

	maClient, err := macl.NewMacAddressClient()
	if err != nil {
		return nil, fmt.Errorf("create mac address client: %w", err)

	}

	maService := masrv.NewMacAddressService(maClient)
	socketService := socketrv.NewSocketService()

	if app.Server, err = server.NewBluetoothServer(maService, socketService); err != nil {
		return nil, fmt.Errorf("get server: %w", err)
	}

	app.closers = append(app.closers, socketService)

	return app, nil
}

func getListener() (net.Listener, error) {
	port, err := common.GetBluetoothApiPort()
	if err != nil {
		return nil, fmt.Errorf("get port: %w", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return listener, nil
}

func (a *App) Run(ctx context.Context) (chan os.Signal, chan error) {
	errCh := make(chan error)
	go func() {
		log.Printf("server listening at %v", a.Listener.Addr())
		if err := a.Server.Serve(a.Listener); err != nil {
			errCh <- fmt.Errorf("serve listener: %w", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return signalCh, errCh
}

func (a *App) Exit(ctx context.Context) {
	a.Server.GracefulStop()
	for _, closer := range a.closers {
		closer.Close(ctx)
	}
}
