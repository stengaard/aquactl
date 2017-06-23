package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type LEDController struct {
	LightConf
	led    *LED
	state  bool
	schedC chan Schedule
	done   chan struct{}
}

func NewLEDController(lc LightConf) (*LEDController, error) {
	led, err := NewLED(lc.Pin)
	if err != nil {
		return nil, err
	}

	c := &LEDController{
		LightConf: lc,
		led:       led,
		schedC:    make(chan Schedule),
		done:      make(chan struct{}),
	}
	go c.run()
	c.schedC <- lc.Schedule

	return c, nil
}

func (s *LEDController) UpdateSchedule(sched Schedule) {
	s.schedC <- sched
}

func (s *LEDController) run() {
	err := s.led.Off()
	if err != nil {
		panic(err)
	}
	var sched Schedule
	t := time.NewTimer(sched.NextTick())
	for {
		select {
		case <-t.C:
			err := s.led.Set(sched.CurrentState())
			if err != nil {
				panic(err)
			}
			t.Reset(sched.NextTick())

		case sched = <-s.schedC:
			t.Reset(10 * time.Millisecond)
		case <-s.done:
			t.Stop()
			s.led.Close()
			return
		case s.schedC <- sched:

		}
	}
}

func (s *LEDController) Close() error {
	close(s.done)
	return nil
}

const (
	ScheduleTicks = 48
	day           = 24 * time.Hour
	durPrTick     = day / ScheduleTicks
)

// 30 minute interval
type Schedule [ScheduleTicks]LEDState

func (s Schedule) State(t time.Time) LEDState {
	return s[stateIdx(t)]
}

func (s Schedule) CurrentState() LEDState {
	return s.State(time.Now())
}

// NextTick returns the time until the next tick is to happen.
func (s Schedule) NextTick() time.Duration {
	now := stateIdx(time.Now())
	return timeTill(now, now+1)
}

func (s Schedule) String() string {

	out := []string{}
	streak := s[0]
	n := 0
	for i := 0; i < len(s); i++ {
		n++
		if s[i] == streak && i < len(s)-1 {
			continue
		}

		out = append(out, fmt.Sprintf("%t:%d", streak, n))
		streak = s[i]
		n = 1
	}
	return strings.Join(out, ",")

}

func (s *Schedule) FromString(in string) error {
	phases := strings.Split(in, ",")
	i := 0
	for _, phase := range phases {
		p := strings.Split(phase, ":")
		if len(p) < 2 {
			return fmt.Errorf("bad schedule: %s", s)
		}

		n, err := strconv.Atoi(p[1])
		if err != nil {
			return err
		}

		val := false
		switch p[0] {
		case "true":
			val = true
		case "false":
			val = false
		default:
			return fmt.Errorf("bad value: %s", p[0])
		}

		fmt.Println(i, n, val)
		// n-1 since it's now an index
		for n > 0 && i < len(s) {
			s[i] = LEDState(val)
			n--
			i++
		}
	}
	return nil
}

func stateIdx(t time.Time) int {
	a := time.Duration(t.Hour()) * time.Hour
	a += time.Duration(t.Minute()) * time.Minute
	a += time.Duration(t.Second()) * time.Second

	return int(a / durPrTick)
}

func timeTill(from, to int) time.Duration {
	if from >= ScheduleTicks {
		from %= ScheduleTicks
	}
	if to >= ScheduleTicks {
		to %= ScheduleTicks
	}

	// switch it around
	if to < from {
		from, to = to, from
	}

	return time.Duration((to - from)) * durPrTick
}

func (s Schedule) MarshalYAML() (interface{}, error) {
	return s.String(), nil
}

func (s *Schedule) UnmarshalYAML(ufn func(interface{}) error) error {
	in := ""
	err := ufn(&in)
	if err != nil {
		return err
	}
	return s.FromString(in)
}
