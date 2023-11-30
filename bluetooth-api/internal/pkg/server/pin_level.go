package server

import (
	"context"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/model"
	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) PinLevel(ctx context.Context, req *pb.PinLevelRequest) (*empty.Empty, error) {
	if req.GetDeviceId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty device id")
	}

	if req.GetPinNumber() < 0 {
		return nil, status.Error(codes.InvalidArgument, "negative pin number")
	}

	addr, err := impl.databaseService.GetMacAddressByID(ctx, int(req.DeviceId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get addr by id: %s", err)
	}
	if addr == "" {
		return nil, status.Errorf(codes.NotFound, "addr not found: %s", err)
	}

	data := model.SetPinData{
		MacAddress:   addr,
		PinNumber:    int(req.PinNumber),
		SetHighLevel: req.High,
	}
	if err = impl.socketService.SetPinLevel(ctx, data); err != nil {
		return nil, status.Errorf(codes.Internal, "set pin level: %s", err)
	}

	return &empty.Empty{}, nil
}
