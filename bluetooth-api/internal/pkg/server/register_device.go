package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) RegisterDevice(
	ctx context.Context,
	req *pb.RegisterDeviceRequest,
) (*empty.Empty, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty device name")
	}
	if req.GetMacAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty mac address")
	}
	if err := impl.databaseService.RegisterDevice(ctx, req.Name, req.MacAddress); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}
