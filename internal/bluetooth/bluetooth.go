package bluetooth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	bl "tinygo.org/x/bluetooth"
)

const (
	deviceName      = "HC-06"
	maxScanDuration = time.Second * 5
)

func GetMACAddress() (string, error) {
	ba := bl.DefaultAdapter
	err := ba.Enable()
	if err != nil {
		return "", fmt.Errorf("adapter enable: %w", err)
	}

	mac, err := getScanResult(ba)
	if err != nil {
		return "", fmt.Errorf("scan: %w", err)
	}
	if mac == "" {
		return "", errors.New("hc-06 not found")
	}

	log.Printf("found hc-06 with mac %s", mac)

	return mac, nil
}

func getScanResult(ba *bl.Adapter) (mac string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxScanDuration)
	defer cancel()

	err = ba.Scan(func(adapter *bl.Adapter, device bl.ScanResult) {
		select {
		case <-ctx.Done():
			adapter.StopScan()
		default:
			if !strings.Contains(device.LocalName(), deviceName) {
				return
			}

			mac = device.Address.String()
			adapter.StopScan()

			return
		}
	})
	return
}
