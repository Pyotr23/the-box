package bluetooth

import (
	"context"
	"errors"

	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	helper "github.com/Pyotr23/the-box/common/pkg/context"
	"google.golang.org/grpc/metadata"
)

const (
	hc05 = "HC-05"
	hc06 = "HC-06"
)

var devicesTypes = []string{hc05, hc06}

type client interface {
	Search(ctx context.Context, deviceNames []string) ([]string, error)
	Blink(ctx context.Context) error
	RegisterDevice(ctx context.Context, name, address string) error
	DevicesList(ctx context.Context, deviceNames []string) ([]model.Device, error)
}

type Service struct {
	c client
}

func NewService(c client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) DevicesMap(ctx context.Context) (map[string]model.Device, error) {
	devices, err := s.c.DevicesList(ctx, devicesTypes)
	if err != nil {
		return nil, err
	}

	var m = make(map[string]model.Device, len(devices))
	for _, d := range devices {
		m[d.MacAddress] = d
	}

	return m, nil
}

func (s *Service) Search(ctx context.Context) ([]string, error) {
	return s.c.Search(ctx, devicesTypes)
}

func (s *Service) Blink(ctx context.Context, addr string) error {
	if addr == "" {
		return errors.New("empty mac address")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, helper.MacAddressKey, addr)

	return s.c.Blink(ctx)
}

func (s *Service) RegisterDevice(ctx context.Context, name, address string) error {
	return s.c.RegisterDevice(ctx, name, address)
}
