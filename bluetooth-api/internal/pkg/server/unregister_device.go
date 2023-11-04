package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) UnregisterDevice(
	ctx context.Context,
	req *pb.UnregisterDeviceRequest,
) (*empty.Empty, error) {
	if req.GetID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty device id")
	}
	if err := impl.databaseService.UnregisterDevice(ctx, int(req.ID)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}
