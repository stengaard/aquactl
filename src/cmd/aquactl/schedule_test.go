package main

import (
	"testing"
	"time"
)

func TestStateIdx(t *testing.T) {

	format := "15:04:05"
	tcs := []struct {
		time string
		idx  int
	}{
		{"00:00:00", 0},
		{"00:29:00", 0},
		{"00:30:01", 1},
		{"00:59:59", 1},
		// ....
		{"23:59:59", 47},
	}

	for _, tc := range tcs {
		tim, err := time.Parse(format, tc.time)
		if err != nil {
			t.Errorf("bad time  %s : %s", tc.time, err)
			continue
		}
		got := stateIdx(tim)
		if got != tc.idx {
			t.Errorf("Got %d - expected %d for time %s",
				got, tc.idx, tc.time)
		}
	}

}

func TestTimeTill(t *testing.T) {

}
