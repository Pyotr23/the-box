package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) CheckPin(
	ctx context.Context,
	req *pb.CheckPinRequest,
) (*pb.CheckPinResponse, error) {
	if req.GetDeviceId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty device id")
	}
	if req.GetPinNumber() < 0 {
		return nil, status.Error(codes.InvalidArgument, "pin number must be positive")
	}

	ma, err := impl.databaseService.GetMacAddressByID(ctx, int(req.DeviceId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get mac address by id %d: %s", req.DeviceId, err)
	}

	checkResult, err := impl.socketService.CheckPin(ctx, ma, int(req.PinNumber))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check pin by mac address: %s", err)
	}

	return &pb.CheckPinResponse{
		IsAvailable: checkResult,
	}, nil
}
