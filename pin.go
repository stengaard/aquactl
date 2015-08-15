package main

import (
	"errors"
	"io"
)

// PinDirection determines if Reads or Writes can be done from a pin.
type PinDirection bool

const (
	PinOut = PinDirection(false)
	PinIn  = PinDirection(true)
)

var (
	ErrIllegalPinOp = errors.New("illegal operation on pin")
)

type readWriteSeeker interface {
	io.Reader
	io.Writer
	io.Closer
}

// DigitalPin can read or write from a single physical, external pin
type DigitalPin interface {
	// SetDirection sets the pin direction and returns the
	// previous direction
	SetDirection(PinDirection) (PinDirection, error)

	// Set the output pin to high (true) or low (false). Returns an error if
	// the pins is configured for reads.
	WriteBool(bool) error

	// Read the current value of pin. Returns a error if the pin is currently
	// configured for writes.
	ReadBool() (bool, error)
	Close() error
}

// NewDigitalPin returns an implementation of DigitalPin.
// The mapping from pinName to physical pins is platform dependendant.
// Returns an error of the pin cannot be found.
func NewDigitalPin(pinName string) (DigitalPin, error) {
	return newDigitalPin(pinName)
}
