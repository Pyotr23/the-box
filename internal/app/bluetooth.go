package app

import (
	"context"
	"fmt"
	"log"

	hardware "github.com/Pyotr23/the-box/internal/bluetooth"
	"github.com/Pyotr23/the-box/internal/rfcomm"
)

const bluetoothName = "bluetooth"

type bluetooth struct {
	sockets []rfcomm.Socket
}

func newBluetooth() *bluetooth {
	return &bluetooth{}
}

func (*bluetooth) Name() string {
	return bluetoothName
}

func (b *bluetooth) Init(ctx context.Context, _ interface{}) (interface{}, error) {
	mac, err := hardware.GetMACAddress()
	if err != nil {
		return nil, fmt.Errorf("get mac address: %w", err)
	}

	socket, err := rfcomm.NewSocket()
	if err != nil {
		return nil, fmt.Errorf("new socket: %w", err)
	}

	err = socket.Connect(mac)
	if err != nil {
		return nil, fmt.Errorf("socket connect: %w", err)
	}

	b.sockets = append(b.sockets, socket)

	return b.sockets, nil
}

func (*bluetooth) SuccessLog() {
	log.Println("init hc-06")
}

func (b *bluetooth) Close(ctx context.Context) error {
	for _, socket := range b.sockets {
		if err := socket.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (*bluetooth) CloseLog() {
	closeLog(bluetoothName)
}
