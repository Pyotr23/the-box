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
	UnregisterDevice(ctx context.Context, id int) error
	DevicesList(ctx context.Context, deviceNames []string) ([]model.Device, error)
	GetDevicesFullInfo(ctx context.Context, ids []int) ([]model.DeviceInfo, error)
	GetTemperature(ctx context.Context, id int) (string, error)
	CheckPin(ctx context.Context, deviceID, pin int) (bool, error)
	SetPinLevel(ctx context.Context, deviceID, pinNumber int, high bool) error
}

type Service struct {
	c client
}

func NewService(c client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) SetPinLevel(ctx context.Context, deviceID, pinNumber int, high bool) error {
	return s.c.SetPinLevel(ctx, deviceID, pinNumber, high)
}

func (s *Service) CheckPin(ctx context.Context, deviceID, pin int) (bool, error) {
	return s.c.CheckPin(ctx, deviceID, pin)
}

func (s *Service) GetTemperature(ctx context.Context, id int) (string, error) {
	return s.c.GetTemperature(ctx, id)
}

func (s *Service) GetDeviceFullInfo(ctx context.Context, id int) (model.DeviceInfo, error) {
	devices, err := s.c.GetDevicesFullInfo(ctx, []int{id})
	if err != nil {
		return model.DeviceInfo{}, err
	}
	if len(devices) == 0 {
		return model.DeviceInfo{}, errors.New("device not found")
	}
	return devices[0], nil
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

func (s *Service) RegisteredDevicesMap(ctx context.Context) (map[string]model.Device, error) {
	devices, err := s.c.DevicesList(ctx, devicesTypes)
	if err != nil {
		return nil, err
	}

	var m = make(map[string]model.Device, len(devices))
	for _, d := range devices {
		if d.ID > 0 {
			m[d.MacAddress] = d
		}
	}

	return m, nil
}

func (s *Service) GetDeviceAliases(ctx context.Context) ([]string, error) {
	devices, err := s.DevicesMap(ctx)
	if err != nil {
		return nil, err
	}

	var aliases = make([]string, 0, len(devices))
	for _, d := range devices {
		var alias string
		if d.Name == "" {
			alias = d.MacAddress
		} else {
			alias = d.Name
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
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

func (s *Service) UnregisterDevice(ctx context.Context, id int) error {
	return s.c.UnregisterDevice(ctx, id)
}
