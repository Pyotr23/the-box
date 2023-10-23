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

type Service struct {
	socketByAddress map[string]socket.Socket
}

func NewSocketService() *Service {
	return &Service{
		socketByAddress: make(map[string]socket.Socket),
	}
}

func (s *Service) Blink(ctx context.Context, macAddress string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	skt, err := s.getConnectedSocket(macAddress)
	if err != nil {
		return fmt.Errorf("get connected socket: %w", err)
	}

	if err = skt.Command(enum.Blink); err != nil {
		return fmt.Errorf("command: %w", err)
	}

	return nil
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
