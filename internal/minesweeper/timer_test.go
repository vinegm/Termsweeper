package minesweeper

import (
	"testing"
	"time"
)

func TestTimerStartStop(t *testing.T) {
	var timer Timer

	if got := timer.ElapsedTime(); got != 0 {
		t.Fatalf("ElapsedTime = %v on a fresh timer, want 0", got)
	}

	timer.StartTimer()
	start := timer.start

	timer.StartTimer() // starting twice must not restart the clock

	if !timer.start.Equal(start) {
		t.Error("StartTimer restarted an already running timer")
	}

	timer.StopTimer()
	stopped := timer.ElapsedTime()

	if stopped <= 0 {
		t.Fatalf("ElapsedTime = %v after stopping, want a positive duration", stopped)
	}

	if got := timer.ElapsedTime(); got != stopped {
		t.Errorf("ElapsedTime = %v, want the frozen %v", got, stopped)
	}
}

func TestTimerFormattedTime(t *testing.T) {
	tests := []struct {
		elapsed time.Duration
		want    string
	}{
		{0, "00:00"},
		{9 * time.Second, "00:09"},
		{75 * time.Second, "01:15"},
		{59*time.Minute + 59*time.Second, "59:59"},
		{time.Hour + 2*time.Minute + 3*time.Second, "01:02:03"},
	}

	for _, test := range tests {
		timer := Timer{elapsed: test.elapsed}

		if got := timer.FormattedTime(); got != test.want {
			t.Errorf("FormattedTime(%v) = %q, want %q", test.elapsed, got, test.want)
		}
	}
}
