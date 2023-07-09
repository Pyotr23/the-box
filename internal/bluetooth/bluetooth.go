package bluetooth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Pyotr23/the-box/internal/helper"
	bl "tinygo.org/x/bluetooth"
)

const (
	deviceName      = "HC-06"
	maxScanDuration = time.Second * 5
)

func GetMACAddress() (string, error) {
	adapter := bl.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return "", fmt.Errorf("adapter enable: %w", err)
	}

	mac, err := getScanResult(adapter)
	if err != nil {
		return "", fmt.Errorf("scan: %w", err)
	}
	if mac == "" {
		return "", errors.New("hc-06 not found")
	}

	log.Printf("found hc-06 with mac %s", mac)

	return mac, nil
}

func getScanResult(adapter *bl.Adapter) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxScanDuration)
	defer cancel()

	helper.Logln(fmt.Sprintf("scanning for '%s'...", deviceName))

	go func() {
		select {
		case <-ctx.Done():
			adapter.StopScan()
		}
	}()

	var macs []string
	err := adapter.Scan(func(adapter *bl.Adapter, device bl.ScanResult) {
		// if !strings.Contains(device.LocalName(), deviceName) {
		// 	return
		// }

		macs = append(macs, device.LocalName()) // device.Address.String())
	})

	helper.Logln(fmt.Sprintf("founded %s - %v", deviceName, macs))

	return "", err
}
