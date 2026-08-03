package main

import (
	"testing"
	"time"
)

// newTestManager returns a Manager that uses the given command instead of
// caffeinate, so tests run without side effects.
func newTestManager(name string, args ...string) *Manager {
	return &Manager{
		Expired: make(chan struct{}, 1),
		newCmd: func(_ ...string) (string, []string) {
			return name, args
		},
	}
}

func TestManagerInitialState(t *testing.T) {
	m := newTestManager("sleep", "10")
	if m.IsRunning() {
		t.Error("expected not running on init")
	}
}

func TestManagerStartStop(t *testing.T) {
	m := newTestManager("sleep", "10")
	t.Cleanup(m.Stop)

	if err := m.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !m.IsRunning() {
		t.Error("expected running after Start")
	}
	m.Stop()
	if m.IsRunning() {
		t.Error("expected not running after Stop")
	}
}

func TestManagerRestartReplacesProcess(t *testing.T) {
	m := newTestManager("sleep", "10")
	t.Cleanup(m.Stop)

	if err := m.Start(); err != nil {
		t.Fatalf("first Start: %v", err)
	}
	if err := m.Start(); err != nil {
		t.Fatalf("second Start: %v", err)
	}
	if !m.IsRunning() {
		t.Error("expected running after restart")
	}
}

func TestManagerNaturalExpiry(t *testing.T) {
	// "true" exits immediately — simulates a timed caffeinate finishing
	m := newTestManager("true")

	if err := m.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	select {
	case <-m.Expired:
		// good
	case <-time.After(2 * time.Second):
		t.Fatal("Expired channel not signalled within 2s")
	}
	if m.IsRunning() {
		t.Error("expected not running after natural expiry")
	}
}

func TestManagerStopDoesNotTriggerExpired(t *testing.T) {
	m := newTestManager("sleep", "10")

	if err := m.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	m.Stop()

	// Give the goroutine a moment to (wrongly) send on Expired
	time.Sleep(100 * time.Millisecond)

	select {
	case <-m.Expired:
		t.Error("Stop should not signal Expired")
	default:
		// correct — channel is empty
	}
}
