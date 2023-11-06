package db

import (
	"context"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/model"
)

type (
	bluetoothRepository interface {
		UpsertDevice(ctx context.Context, name, macAddress string) error
		GetByMacAddresses(ctx context.Context, macAddresses []string) ([]model.DbDevice, error)
		GetByIDs(ctx context.Context, ids []int) ([]model.DbDevice, error)
		DeleteDevice(ctx context.Context, id int) error
	}

	Service struct {
		repo bluetoothRepository
	}
)

func NewDbService(repo bluetoothRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterDevice(ctx context.Context, name, address string) error {
	return s.repo.UpsertDevice(ctx, name, address)
}

func (s *Service) UnregisterDevice(ctx context.Context, id int) error {
	return s.repo.DeleteDevice(ctx, id)
}

func (s *Service) GetDeviceByAddressMap(
	ctx context.Context,
	addresses []string,
) (map[string]model.Device, error) {
	dbDevices, err := s.repo.GetByMacAddresses(ctx, addresses)
	if err != nil {
		return nil, err
	}

	if len(dbDevices) == 0 {
		return nil, nil
	}

	var m = make(map[string]model.Device, len(dbDevices))
	for _, d := range dbDevices {
		m[d.MacAddress] = model.Device{
			ID:         d.ID,
			MacAddress: d.MacAddress,
			Name:       d.Name,
		}
	}

	return m, nil
}

func (s *Service) GetDeviceByIDs(
	ctx context.Context,
	ids []int,
) ([]model.DeviceInfo, error) {
	dbDevices, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(dbDevices) == 0 {
		return nil, nil
	}

	var res = make([]model.DeviceInfo, 0, len(dbDevices))
	for _, d := range dbDevices {
		res = append(res, model.DeviceInfo{
			ID:         d.ID,
			MacAddress: d.MacAddress,
			Name:       d.Name,
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
		})
	}

	return res, nil
}
