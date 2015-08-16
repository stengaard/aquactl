// +build !linux

package main

import (
	"fmt"
	"log"
)

func newDigitalPin(pinid uint) (DigitalPin, error) {
	return &mockPin{
		name:  fmt.Sprintf("pin%d", pinid),
		value: false,
		dir:   PinOut,
	}, nil
}

type mockPin struct {
	name  string
	value bool
	dir   PinDirection
}

func (m *mockPin) Close() error {
	return nil
}

func (m *mockPin) ReadBool() (bool, error) {
	if m.dir == PinIn {
		return m.value, nil
	}

	return false, ErrIllegalPinOp
}

func (m *mockPin) WriteBool(v bool) error {
	if m.dir == PinOut {
		m.value = v
		out := "off"
		if m.value {
			out = "on"
		}

		log.Println(m.name, out)
		return nil
	}
	return ErrIllegalPinOp
}

func (m *mockPin) SetDirection(n PinDirection) (PinDirection, error) {
	old := m.dir
	m.dir = n
	return old, nil
}
