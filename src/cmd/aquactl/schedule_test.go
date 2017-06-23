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

func TestScheduleSerialisation(t *testing.T) {
	tests := []struct {
		in, out string
		err     bool
	}{
		{"true:48", "true:48", false},
		{"false:48", "false:48", false},
		{"true:1,false:46,true:1", "true:1,false:46,true:1", false},
		//		{"false:500,false:45", "false:48", false},
		//		{"false:47,true:1", "false:47,true:1", false},
		//		{"true:5,false:45", "true:5,false:43", false},
		//		{"false:3,false:45", "false:48", false},
	}

	for i, test := range tests {
		s := &Schedule{}
		err := s.FromString(test.in)
		if err != nil {
			t.Error(err, "in test", i+1, "with text", test.in)
		}

		res := s.String()
		if res != test.out {
			t.Error("got", res, "but expected", test.out, "on", test.in)
		}

	}
}
