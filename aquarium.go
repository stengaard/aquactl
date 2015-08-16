package main

import (
	"flag"
	"time"
)

func main() {
	flag.Parse()

	pinid := uint(18)
	pin, err := NewDigitalPin(pinid)
	if err != nil {
		panic(err)
	}
	sched := Schedule{
		Times: []time.Duration{
			1 * time.Second,
			2 * time.Second,
		},
		Pin:   pin,
		Name:  pinid,
		State: true,
	}
	sched.run()

}

func (s *Schedule) run() {
	_, err := s.Pin.SetDirection(PinOut)
	if err != nil {
		panic(err)
	}
	for {
		for i := 0; i < len(s.Times); i++ {
			s.State = !s.State
			err := s.Pin.WriteBool(s.State)
			if err != nil {
				panic(err)
			}
			time.Sleep(s.Times[i])
		}
	}
}

type Schedule struct {
	Times []time.Duration
	Pin   DigitalPin
	Name  uint
	State bool
}
