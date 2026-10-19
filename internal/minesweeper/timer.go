package minesweeper

import (
	"fmt"
	"time"
)

// Timer measures how long a game takes.
type Timer struct {
	running bool

	start   time.Time
	elapsed time.Duration
}

// StartTimer starts timing from zero, unless already running.
func (timer *Timer) StartTimer() {
	if timer.running {
		return
	}

	timer.running = true
	timer.start = time.Now()
	timer.elapsed = 0
}

// StopTimer stops the timer, keeping the elapsed duration.
func (timer *Timer) StopTimer() {
	if timer.start.IsZero() {
		return
	}

	timer.elapsed = time.Since(timer.start)
	timer.start = time.Time{}
	timer.running = false
}

// ElapsedTime is how long the timer has run.
func (timer *Timer) ElapsedTime() time.Duration {
	if timer.running {
		return time.Since(timer.start)
	}

	return timer.elapsed
}

// FormattedTime is ElapsedTime as text.
func (timer *Timer) FormattedTime() string {
	duration := timer.ElapsedTime()

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
