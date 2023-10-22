package bluetooth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	b "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	common "github.com/Pyotr23/the-box/common/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	host      = "localhost"
	separator = ":"
)

var (
	defaultDuration = time.Second * 5
)

type BluetoothClient interface{}

type Client struct {
	api b.BluetoothClient
}

func NewClient() (*Client, error) {
	port, err := common.GetBluetoothApiPort()
	if err != nil {
		return nil, errors.New("empty port")
	}

	serverAddr := host + separator + strconv.Itoa(port)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return &Client{
		api: b.NewBluetoothClient(cc),
	}, nil
}

func (c *Client) Search(ctx context.Context, deviceNames []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	req := &b.SearchRequest{
		DeviceNames: deviceNames,
	}
	resp, err := c.api.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("api call: %w", err)
	}

	return resp.GetMacAddresses(), nil
}
