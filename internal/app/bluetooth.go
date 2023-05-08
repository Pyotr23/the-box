package app

import (
	"context"
	"fmt"
	"log"

	hardware "github.com/Pyotr23/the-box/internal/bluetooth"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

const bluetoothName = "bluetooth"

type bluetooth struct{}

func newBluetooth() *bluetooth {
	return &bluetooth{}
}

func (*bluetooth) Name() string {
	return bluetoothName
}

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

func (*bluetooth) Close(ctx context.Context, a *App) (err error) {
	for _, socket := range a.sockets {
		err = socket.Close()
	}
	return
}

func (*bluetooth) CloseLog() {
	closeLog(bluetoothName)
}
