package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	client "github.com/Pyotr23/the-box/bluetooth-api/internal/client/mac_address"
	pb "github.com/Pyotr23/the-box/bluetooth-api/pkg/pb/bluetooth"
	"github.com/Pyotr23/the-box/common/pkg/config"
	"google.golang.org/grpc"
)

type Implementation struct {
	MacAddressClient client.MacAddressClient
	pb.UnimplementedBluetoothServer
}

func (impl *Implementation) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	m, err := impl.MacAddressClient.GetAddressesByNameMap(in.GetDeviceNames())
	if err != nil {
		return nil, err
	}

	log.Println(m)

	var resp = &pb.SearchResponse{
		Items: make([]*pb.SearchResponse_AddressesByName, 0, len(m)),
	}
	for deviceName, addresses := range m {
		resp.Items = append(resp.Items, &pb.SearchResponse_AddressesByName{
			Name:         deviceName,
			MacAddresses: addresses,
		})
	}

	return resp, nil
}

type App struct {
	Listener net.Listener
	Server   *grpc.Server
}

func NewApp() (*App, error) {
	var (
		err error
		app = &App{}
	)

	if app.Listener, err = getListener(); err != nil {
		return nil, fmt.Errorf("get listener: %w", err)
	}

	if app.Server, err = getServer(); err != nil {
		return nil, fmt.Errorf("get server: %w", err)
	}

	return app, nil
}

func getListener() (net.Listener, error) {
	port, err := config.GetBluetoothApiPort()
	if err != nil {
		return nil, fmt.Errorf("get api port: %w", err)
	}
	if port == 0 {
		return nil, errors.New("empty api port")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return listener, nil
}

func getServer() (*grpc.Server, error) {
	client, err := maclient.NewClient()
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	impl := &Implementation{
		MacAddressClient: client,
	}

	_, _ = impl.Search(context.Background(), &pb.SearchRequest{DeviceNames: []string{"HC-05", "HC-06"}})

	pb.RegisterBluetoothServer(s, impl)

	return s, nil
}

func (a *App) Run(ctx context.Context) (chan os.Signal, chan error) {
	errCh := make(chan error)
	go func() {
		log.Printf("server listening at %v", a.Listener.Addr())
		if err := a.Server.Serve(a.Listener); err != nil {
			errCh <- fmt.Errorf("serve listener: %w", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return signalCh, errCh
}

func (a *App) Exit(_ context.Context) {
	a.Server.GracefulStop()
}
