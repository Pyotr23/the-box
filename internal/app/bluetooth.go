package app

import (
	"context"
	"fmt"
	"log"

	hardware "github.com/Pyotr23/the-box/internal/bluetooth"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

type bluetooth struct{}

func (*bluetooth) Init(ctx context.Context, a *App) error {
	mac, err := hardware.GetMACAddress()
	if err != nil {
		return fmt.Errorf("get mac address: %w", err)
	}

	socket, err := rfcomm.NewSocket()
	if err != nil {
		return fmt.Errorf("new socket: %w", err)
	}

	err = socket.Connect(mac)
	if err != nil {
		return fmt.Errorf("socket connect: %w", err)
	}

	a.sockets = append(a.sockets, socket)

	return nil
}

func (*bluetooth) SuccessLog() {
	log.Println("init hc-06")
}
