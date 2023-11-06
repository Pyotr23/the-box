package server

import (
	"context"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (impl *Implementation) GetDevicesFullInfo(
	ctx context.Context,
	in *pb.GetDevicesFullInfoRequest,
) (*pb.GetDevicesFullInfoResponse, error) {
	if len(in.GetIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no ids")
	}

	var ids = make([]int, 0, len(in.Ids))
	for _, id := range in.Ids {
		if id == 0 {
			return nil, status.Error(codes.InvalidArgument, "empty id in ids")
		}
		ids = append(ids, int(id))
	}

	devices, err := impl.databaseService.GetDeviceByIDs(ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "blink: %s", err)
	}

	var res = make([]*pb.GetDevicesFullInfoResponse_Device, 0, len(devices))
	for _, d := range devices {
		res = append(res, &pb.GetDevicesFullInfoResponse_Device{
			ID:         int32(d.ID),
			MacAddress: d.MacAddress,
			Name:       d.Name,
			CreatedAt:  timestamppb.New(d.CreatedAt),
			UpdatedAt:  timestamppb.New(d.UpdatedAt),
		})
	}
	return &pb.GetDevicesFullInfoResponse{
		Devices: res,
	}, nil
}
