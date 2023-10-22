package mac_address

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	bt "tinygo.org/x/bluetooth"
)

const maxScanDuration = time.Second * 3

type Client struct {
	adapter *bt.Adapter
}

func NewMacAddressClient() (Client, error) {
	adapter := bt.DefaultAdapter
	return Client{
		adapter: adapter,
	}, adapter.Enable()
}

func (c Client) GetAddressesByName(deviceNames []string) (map[string][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxScanDuration)
	defer cancel()

	log.Printf("scanning for %v...", deviceNames)

	go func() {
		select {
		case <-ctx.Done():
			c.adapter.StopScan()
		}
	}()

	var (
		macsByDeviceName = make(map[string][]string)
		macAdressesSet   = make(map[string]bool)
	)
	err := c.adapter.Scan(func(adapter *bt.Adapter, current bt.ScanResult) {
		macAddress := current.Address.String()

		for _, deviceName := range deviceNames {
			if strings.EqualFold(current.LocalName(), deviceName) && !macAdressesSet[macAddress] {
				macsByDeviceName[deviceName] = append(macsByDeviceName[deviceName], current.Address.String())
				macAdressesSet[macAddress] = true
				return
			}
		}

		return
	})
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return macsByDeviceName, nil
}
