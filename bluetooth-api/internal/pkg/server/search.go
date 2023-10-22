package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	addresses, err := impl.MacAddressService.Search(ctx, in.GetDeviceNames())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SearchResponse{
		MacAddresses: addresses,
	}, nil
}
