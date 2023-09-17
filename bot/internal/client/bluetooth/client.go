package bluetooth

import (
	"errors"
	"fmt"
	"strconv"

	b "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	common "github.com/Pyotr23/the-box/common/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	host      = "localhost"
	separator = ":"
)

type BluetoothClient interface{}

type Client struct {
	api b.BluetoothClient
}

func NewClient() (Client, error) {
	cfg, err := common.GetBlutoothApiConfig()
	if err != nil {
		return Client{}, fmt.Errorf("get bluetooth api config: %w", err)
	}

	port := cfg.Port
	if port == 0 {
		return Client{}, errors.New("empty port")
	}

	serverAddr := host + separator + strconv.Itoa(port)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return Client{}, fmt.Errorf("dial: %w", err)
	}

	return Client{
		api: b.NewBluetoothClient(cc),
	}, nil
}
