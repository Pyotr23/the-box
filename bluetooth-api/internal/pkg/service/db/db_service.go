package db

import "context"

type (
	bluetoothRepository interface {
		UpsertDevice(ctx context.Context, name, macAddress string) error
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
