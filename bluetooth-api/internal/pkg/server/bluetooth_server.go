package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc"
)

type (
	macAddressService interface {
		Search(ctx context.Context, devicesNames []string) ([]string, error)
	}

	dbService interface {
		RegisterDevice(ctx context.Context, name, address string) error
	}

	socketService interface {
		Blink(ctx context.Context, macAddress string) error
	}
)

type Implementation struct {
	MacAddressService macAddressService
	DatabaseService   dbService
	SocketService     socketService
	pb.UnimplementedBluetoothServer
}

func NewBluetoothServer(
	macAddressService macAddressService,
	socketService socketService,
	dbService dbService,
) (*grpc.Server, error) {
	s := grpc.NewServer()
	impl := &Implementation{
		MacAddressService: macAddressService,
		SocketService:     socketService,
		DatabaseService:   dbService,
	}

	pb.RegisterBluetoothServer(s, impl)

	return s, nil
}
