package server

import (
	"context"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/model"
	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc"
)

type (
	macAddressService interface {
		Search(ctx context.Context, devicesNames []string) ([]string, error)
	}

	dbService interface {
		RegisterDevice(ctx context.Context, name, address string) error
		UnregisterDevice(ctx context.Context, id int) error
		GetDeviceByAddressMap(ctx context.Context, addresses []string) (map[string]model.Device, error)
		GetDeviceByIDs(ctx context.Context, ids []int) ([]model.DeviceInfo, error)
	}

	socketService interface {
		Blink(ctx context.Context, macAddress string) error
	}
)

type Implementation struct {
	macAddressService macAddressService
	databaseService   dbService
	socketService     socketService
	pb.UnimplementedBluetoothServer
}

func NewBluetoothServer(
	macAddressService macAddressService,
	socketService socketService,
	dbService dbService,
) (*grpc.Server, error) {
	s := grpc.NewServer()
	impl := &Implementation{
		macAddressService: macAddressService,
		socketService:     socketService,
		databaseService:   dbService,
	}

	pb.RegisterBluetoothServer(s, impl)

	return s, nil
}
