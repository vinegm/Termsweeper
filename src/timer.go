package src

import (
	"fmt"
	"time"
)

// encapsulates timing logic.
type Timer struct {
	running bool

	start   time.Time
	elapsed time.Duration
}

// initializes the timer. If already running does nothing.
func (timer *Timer) StartTimer() {
	if !timer.start.IsZero() && timer.elapsed == 0 {
		return
	}

	timer.running = true
	timer.start = time.Now()
	timer.elapsed = 0
}

// stops the timer and records the elapsed duration.
func (timer *Timer) StopTimer() {
	if timer.start.IsZero() {
		return
	}

	timer.elapsed = time.Since(timer.start)
	timer.start = time.Time{}
	timer.running = false
}

// Elapsed returns the current elapsed duration.
func (timer *Timer) ElapsedTime() time.Duration {
	if timer.running {
		return time.Since(timer.start)
	}

	return timer.elapsed
}

// Formatted returns a human-readable elapsed time string.
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
