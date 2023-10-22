package socket

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/enum"
	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/socket"
)

var defaultTimeout = time.Second * 5

type Service struct{}

func NewSocketService() Service {
	return Service{}
}

func (s *Service) Blink(ctx context.Context, macAddress string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	skt, err := socket.NewSocket()
	if err != nil {
		return fmt.Errorf("new socket: %w", err)
	}

	defer func() {
		if dErr := skt.Close(); dErr != nil {
			log.Printf("close socket error: %s", dErr)
		}
	}()

	if err = skt.Command(enum.Blink); err != nil {
		return fmt.Errorf("command: %w", err)
	}

	return nil
}
