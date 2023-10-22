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
)

type Implementation struct {
	MacAddressService macAddressService
	pb.UnimplementedBluetoothServer
}

func NewBluetoothServer(service macAddressService) (*grpc.Server, error) {
	s := grpc.NewServer()
	impl := &Implementation{
		MacAddressService: service,
	}

	pb.RegisterBluetoothServer(s, impl)

	return s, nil
}
