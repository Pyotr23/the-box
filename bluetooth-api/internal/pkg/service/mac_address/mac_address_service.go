package mac_address

import (
	"context"
)

type (
	searcher interface {
		GetAddressesByName(ctx context.Context, deviceNames []string) (map[string][]string, error)
	}

	Service struct {
		client searcher
	}
)

func NewMacAddressService(searcher searcher) *Service {
	return &Service{
		client: searcher,
	}
}

func (s *Service) Search(ctx context.Context, devicesNames []string) ([]string, error) {
	addressesByDeviceName, err := s.client.GetAddressesByName(ctx, devicesNames)
	if err != nil {
		return nil, err
	}

	if len(addressesByDeviceName) == 0 {
		return nil, nil
	}

	var res = make([]string, 0, len(addressesByDeviceName))
	for _, addresses := range addressesByDeviceName {
		res = append(res, addresses...)
	}

	return res, nil
}
