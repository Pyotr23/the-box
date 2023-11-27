package model

import "errors"

const (
	PinIsBusyErrText = "pin is busy"
)

var (
	PinIsBusyError = errors.New(PinIsBusyErrText)
)
