package bluetooth

import (
	"context"
	"errors"

	helper "github.com/Pyotr23/the-box/common/pkg/context"
	"google.golang.org/grpc/metadata"
)

const (
	hc05 = "HC-05"
	hc06 = "HC-06"
)

type client interface {
	Search(ctx context.Context, deviceNames []string) ([]string, error)
	Blink(ctx context.Context) error
	RegisterDevice(ctx context.Context, name, address string) error
}

type Service struct {
	c client
}

func NewService(c client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) Search(ctx context.Context) ([]string, error) {
	return s.c.Search(ctx, []string{hc05, hc06})
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
