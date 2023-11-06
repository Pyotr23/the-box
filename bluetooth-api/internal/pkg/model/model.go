package model

import "time"

type Device struct {
	ID         int
	MacAddress string
	Name       string
}

type DeviceInfo struct {
	ID         int
	MacAddress string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
