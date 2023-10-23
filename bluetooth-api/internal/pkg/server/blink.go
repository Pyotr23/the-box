package server

import (
	"context"

	helper "github.com/Pyotr23/the-box/common/pkg/context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) Blink(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	addrs := metadata.ValueFromIncomingContext(ctx, helper.MacAddressKey)
	if len(addrs) == 0 || addrs[0] == "" {
		return nil, status.Error(codes.InvalidArgument, "empty mac address")
	}

	if err := impl.SocketService.Blink(ctx, addrs[0]); err != nil {
		return nil, status.Errorf(codes.Internal, "blink: %s", err)
	}

	return &empty.Empty{}, nil
}
