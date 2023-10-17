package rfcomm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/Pyotr23/the-box/bot/internal/enum"
	"golang.org/x/sys/unix"
)

const (
	defaultSize = 64
	errorByte   = byte(0)
)

type Socket struct {
	fd int
}

func NewSocket() (s Socket, err error) {
	s.fd, err = unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	return s, err
}

func (s Socket) Connect(mac string) error {
	leMAC, err := littleEndian(mac)
	if err != nil {
		return fmt.Errorf("little endian: %w", err)
	}

	err = unix.Connect(s.fd, &unix.SockaddrRFCOMM{
		Addr:    leMAC,
		Channel: 1,
	})
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return nil
}

func (s Socket) Close() error {
	return unix.Close(s.fd)

}

func (s Socket) Query(b enum.Code) (string, error) {
	if err := s.Command(b); err != nil {
		return "", fmt.Errorf("write only: %w", err)
	}

	answer, err := s.read(defaultSize)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	return string(answer), nil
}

func (s Socket) Command(b enum.Code) error {
	if err := s.write([]byte{byte(b)}); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	if err := s.readError(); err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	return nil
}

func (s Socket) SendText(b enum.Code, bs []byte) error {
	if err := s.write([]byte{byte(b)}); err != nil {
		return fmt.Errorf("write code: %w", err)
	}

	if err := s.write(bs); err != nil {
		return fmt.Errorf("write text: %w", err)
	}

	if err := s.readError(); err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	return nil
}

func (s Socket) write(data []byte) error {
	_, err := unix.Write(s.fd, data)
	return err
}

func (s Socket) readError() error {
	successBytes, err := s.read(1)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if len(successBytes) == 0 {
		return errors.New("no success bytes")
	}

	if successBytes[0] == errorByte {
		msg, err := s.read(defaultSize)
		if err != nil {
			return fmt.Errorf("read error message: %w", err)
		}
		return errors.New(string(msg))
	}

	return nil
}

func (s Socket) read(size int) ([]byte, error) {
	answer := make([]byte, size)
	_, err := unix.Read(s.fd, answer)
	return answer, err
}

// littleEndian converts MAC address string representation to little-endian byte array.
func littleEndian(mac string) ([6]byte, error) {
	var res [6]byte
	for i, cur := range strings.Split(mac, ":") {
		u, err := strconv.ParseUint(cur, 16, 8)
		if err != nil {
			return res, err
		}

		res[len(res)-1-i] = byte(u)
	}
	return res, nil
}
