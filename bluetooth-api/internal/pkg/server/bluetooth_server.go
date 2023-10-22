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

	socketService interface {
		Blink(ctx context.Context, macAddress string) error
	}
)

type Implementation struct {
	MacAddressService macAddressService
	SocketService     socketService
	pb.UnimplementedBluetoothServer
}

func NewBluetoothServer(service macAddressService, socketService socketService) (*grpc.Server, error) {
	s := grpc.NewServer()
	impl := &Implementation{
		MacAddressService: service,
		SocketService:     socketService,
	}

	pb.RegisterBluetoothServer(s, impl)

	return s, nil
}
