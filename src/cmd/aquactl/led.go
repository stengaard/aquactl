package main

import "sync"

type LEDState bool

const (
	LEDOn  = true
	LEDOff = false
)

type LED struct {
	mux     sync.Mutex
	pin     DigitalPin
	current LEDState
}

func NewLED(pin uint) (*LED, error) {
	p, err := NewDigitalPin(pin)
	if err != nil {
		return nil, err
	}
	_, err = p.SetDirection(PinOut)
	p.WriteBool(false)
	return &LED{pin: p}, nil
}

func (l *LED) Set(s LEDState) error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.set(s)
}

func (l *LED) set(s LEDState) error {
	err := l.pin.WriteBool(bool(s))
	if err == nil {
		return err
	}

	l.current = s

	return nil
}

func (l *LED) Invert() error {
	l.mux.Lock()
	defer l.mux.Unlock()

	return l.set(!l.current)
}

func (l *LED) On() error    { return l.Set(LEDOn) }
func (l *LED) Off() error   { return l.Set(LEDOff) }
func (l *LED) Close() error { return l.pin.Close() }
