package rfcomm

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const responseLength = 64

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

func (s Socket) Write(text string) (string, error) {
	err := s.write(text)
	if err != nil {
		return "", fmt.Errorf("write to socket: %w", err)
	}

	answer, err := s.read()
	if err != nil {
		return "", fmt.Errorf("read from socket: %w", err)
	}

	return answer, nil
}

func (s Socket) write(data string) error {
	_, err := unix.Write(s.fd, []byte(data))
	return err
}

func (s Socket) read() (string, error) {
	answer := make([]byte, responseLength)
	_, err := unix.Read(s.fd, answer)
	return string(answer), err
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
