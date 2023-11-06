package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) GetTemperature(
	ctx context.Context,
	in *pb.GetTemperatureRequest,
) (*pb.GetTemperatureResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty id")
	}

	ma, err := impl.databaseService.GetMacAddressByID(ctx, int(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get mac address by id: %s", err)
	}
	if ma == "" {
		return nil, status.Error(codes.NotFound, "mac address not found by id")
	}

	t, err := impl.socketService.GetTemperature(ctx, ma)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get temperature: %s", err)
	}

	return &pb.GetTemperatureResponse{
		Value: t,
	}, nil
}
