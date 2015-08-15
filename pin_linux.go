package main

import (
	"os"
	"path"
)

type sysFsPin struct {
	dir readWriteSeeker
	val readWriteSeeker
}

func newDigitalPin(pin string) (DigitalPin, error) {
	var err error
	p := &sysFsPin{}

	p.dir, err = os.Open(path.Join(pinPath, "value"))
	if err != nil {
		return nil, err
	}

	p.val, err = os.Open(path.Join(pinPath, "direction"))
	if err != nil {
		p.dir.Close()
		return nil, err
	}

	return p, nil
}

func (p *sysFsPin) SetDirection(PinDirection) error {
	return nil
}

func (p *sysFsPin) Close() error {
	p.dir.Close() // ignoring err
	return p.val.Close()
}

func (p *sysFsPin) WriteBool(v bool) error {
	b := byte('0')
	if v {
		b = '1'
	}
	_, err := p.val.Write([]byte{b})
	return err
}

func (p *sysFsPin) ReadBool() (bool, error) {
	t := make([]byte, 2)
	n, err := p.val.Read(t)
	r := n > 1 && t[0] == '1'
	return r, err
}
