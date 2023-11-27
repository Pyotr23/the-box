package socket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/enum"
	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/model"
	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/socket"
)

var (
	defaultTimeout    = time.Second * 5
	pinIsBusyErrorMsg = "pin is busy"
)

type Service struct {
	socketByAddress map[string]socket.Socket
}

func NewSocketService() *Service {
	return &Service{
		socketByAddress: make(map[string]socket.Socket),
	}
}

func (s *Service) CheckPin(_ context.Context, macAddress string, pin int) (bool, error) {
	skt, err := s.getConnectedSocket(macAddress)
	if err != nil {
		return false, fmt.Errorf("get connected socket: %w", err)
	}

	err = skt.SendInt(enum.CheckPinCode, pin)
	if err != nil {
		if errors.Is(err, model.PinIsBusyError) {
			return false, nil
		}
		return false, fmt.Errorf("send int: %w", err)
	}

	return true, nil
}

func (s *Service) Blink(_ context.Context, macAddress string) error {
	skt, err := s.getConnectedSocket(macAddress)
	if err != nil {
		return fmt.Errorf("get connected socket: %w", err)
	}

	if err = skt.Command(enum.Blink); err != nil {
		return fmt.Errorf("command: %w", err)
	}

	return nil
}

func (s *Service) GetTemperature(_ context.Context, macAddress string) (string, error) {
	skt, err := s.getConnectedSocket(macAddress)
	if err != nil {
		return "", fmt.Errorf("get connected socket: %w", err)
	}

	t, err := skt.Query(enum.TemperatureCode)
	if err != nil {
		return "", fmt.Errorf("query: %w", err)
	}

	return t, nil
}

func (s *Service) Close(ctx context.Context) {
	log.Print("close sockets...")
	defer log.Print("sockets closed")

	for _, skt := range s.socketByAddress {
		if err := skt.Close(); err != nil {
			log.Printf("close socket error: %s", err)
		}
	}
}

func (s *Service) getConnectedSocket(addr string) (socket.Socket, error) {
	if skt, ok := s.socketByAddress[addr]; ok {
		return skt, nil
	}

	skt, err := socket.NewSocket()
	if err != nil {
		return socket.Socket{}, fmt.Errorf("new socket: %w", err)
	}

	if err = skt.Connect(addr); err != nil {
		return socket.Socket{}, fmt.Errorf("connect: %w", err)
	}

	s.socketByAddress[addr] = skt

	return skt, nil
}
