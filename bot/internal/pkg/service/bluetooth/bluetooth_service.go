package bluetooth

import (
	"context"
)

const (
	hc05 = "HC-05"
	hc06 = "HC-06"
)

type client interface {
	Search(ctx context.Context, deviceNames []string) ([]string, error)
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
