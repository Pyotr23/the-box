package server

import (
	"context"
	"fmt"

	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *Implementation) DevicesList(
	ctx context.Context,
	req *pb.DevicesListRequest,
) (*pb.DevicesListResponse, error) {
	if len(req.GetDeviceTypes()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty device types")
	}

	addresses, err := impl.macAddressService.Search(ctx, req.GetDeviceTypes())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("search: %s", err))
	}

	deviceByAddress, err := impl.databaseService.GetDeviceByAddressMap(ctx, addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("get devices by addresses: %s", err))
	}

	var res = make([]*pb.DevicesListResponse_Device, 0, len(addresses))
	for _, a := range addresses {
		if device, ok := deviceByAddress[a]; ok {
			res = append(res, &pb.DevicesListResponse_Device{
				ID:         int32(device.ID),
				MacAddress: device.MacAddress,
				Name:       device.Name,
			})
			continue
		}

		res = append(res, &pb.DevicesListResponse_Device{
			MacAddress: a,
		})
	}

	return &pb.DevicesListResponse{
		Devices: res,
	}, nil
}
