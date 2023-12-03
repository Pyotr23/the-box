package bluetooth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	b "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	common "github.com/Pyotr23/the-box/common/pkg/config"
	"github.com/Pyotr23/the-box/common/pkg/convert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	host      = "localhost"
	separator = ":"
)

var (
	defaultDuration    = time.Second * 5
	longDefaultDuraion = time.Second * 10
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

func (c *Client) SetPinLevel(ctx context.Context, deviceID, pinNumber int, high bool) error {
	ctx, cancel := context.WithTimeout(ctx, longDefaultDuraion)
	defer cancel()

	req := &b.PinLevelRequest{
		DeviceId:  int32(deviceID),
		PinNumber: int32(pinNumber),
		High:      high,
	}
	_, err := c.api.PinLevel(ctx, req)
	if err != nil {
		return fmt.Errorf("api call: %w", err)
	}

	return nil
}

func (c *Client) DevicesList(ctx context.Context, deviceNames []string) ([]model.Device, error) {
	ctx, cancel := context.WithTimeout(ctx, longDefaultDuraion)
	defer cancel()

	req := &b.DevicesListRequest{
		DeviceTypes: deviceNames,
	}
	resp, err := c.api.DevicesList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("api call: %w", err)
	}

	var res = make([]model.Device, 0, len(resp.GetDevices()))
	for _, d := range resp.GetDevices() {
		res = append(res, model.Device{
			ID:         int(d.ID),
			MacAddress: d.MacAddress,
			Name:       d.Name,
		})
	}

	return res, nil
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

func (c *Client) Blink(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, longDefaultDuraion)
	defer cancel()

	if _, err := c.api.Blink(ctx, &emptypb.Empty{}); err != nil {
		return fmt.Errorf("api call: %w", err)
	}

	return nil
}

func (c *Client) RegisterDevice(ctx context.Context, name, address string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	req := &b.RegisterDeviceRequest{
		Name:       name,
		MacAddress: address,
	}
	if _, err := c.api.RegisterDevice(ctx, req); err != nil {
		return fmt.Errorf("api call: %w", err)
	}

	return nil
}

func (c *Client) UnregisterDevice(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	req := &b.UnregisterDeviceRequest{
		ID: int32(id),
	}
	if _, err := c.api.UnregisterDevice(ctx, req); err != nil {
		return fmt.Errorf("api call: %w", err)
	}

	return nil
}

func (c *Client) GetDevicesFullInfo(ctx context.Context, ids []int) ([]model.DeviceInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	req := &b.GetDevicesFullInfoRequest{
		Ids: convert.ConvertIDs[int, int32](ids),
	}
	resp, err := c.api.GetDevicesFullInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("api call: %w", err)
	}

	if len(resp.GetDevices()) == 0 {
		return nil, nil
	}

	var res = make([]model.DeviceInfo, 0, len(resp.GetDevices()))
	for _, d := range resp.GetDevices() {
		res = append(res, model.DeviceInfo{
			ID:         int(d.ID),
			MacAddress: d.MacAddress,
			Name:       d.Name,
			CreatedAt:  d.CreatedAt.AsTime(),
			UpdatedAt:  d.UpdatedAt.AsTime(),
		})
	}

	return res, nil
}

func (c *Client) GetTemperature(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	req := &b.GetTemperatureRequest{
		Id: int32(id),
	}
	resp, err := c.api.GetTemperature(ctx, req)
	if err != nil {
		return "", fmt.Errorf("api call: %w", err)
	}

	return resp.GetValue(), nil
}

func (c *Client) CheckPin(ctx context.Context, deviceID, pin int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	req := &b.CheckPinRequest{
		DeviceId:  int32(deviceID),
		PinNumber: int32(pin),
	}
	resp, err := c.api.CheckPin(ctx, req)
	if err != nil {
		return false, fmt.Errorf("api call: %w", err)
	}

	return resp.IsAvailable, nil
}
