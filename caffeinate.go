package main

import (
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// Manager owns a single caffeinate child process. Expired is closed when
// a timed run finishes on its own (not when Stop is called).
type Manager struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	Expired chan struct{}
	// newCmd returns the executable and args to run. Defaults to caffeinate.
	// Overridden in tests.
	newCmd func(caffeArgs ...string) (string, []string)
}

func NewManager() *Manager {
	return &Manager{
		Expired: make(chan struct{}, 1),
		newCmd: func(caffeArgs ...string) (string, []string) {
			return "caffeinate", caffeArgs
		},
	}
}

func (m *Manager) Start() error {
	return m.launch("-dims")
}

func (m *Manager) StartTimed(hours int) error {
	return m.launch("-dims", "-t", strconv.Itoa(hours*3600))
}

func (m *Manager) launch(args ...string) error {
	m.mu.Lock()
	// Signal any running process to exit. We do NOT call Wait() here —
	// the goroutine started for that process owns the Wait() call.
	// Setting m.cmd = nil before releasing the lock means the goroutine
	// will see m.cmd != cmd and will not send to Expired.
	if m.cmd != nil {
		m.cmd.Process.Signal(os.Interrupt)
		m.cmd = nil
	}

	name, fullArgs := m.newCmd(args...)
	cmd := exec.Command(name, fullArgs...)
	if err := cmd.Start(); err != nil {
		m.mu.Unlock()
		return err
	}
	m.cmd = cmd
	// Drain any buffered expiry from a prior timed session — stale after a new launch.
	select {
	case <-m.Expired:
	default:
	}
	m.mu.Unlock()

	go func() {
		cmd.Wait()
		m.mu.Lock()
		natural := m.cmd == cmd
		if natural {
			m.cmd = nil
		}
		m.mu.Unlock()
		if natural {
			select {
			case m.Expired <- struct{}{}:
			default:
			}
		}
	}()
	return nil
}

func (m *Manager) Stop() {
	m.mu.Lock()
	if m.cmd == nil {
		m.mu.Unlock()
		return
	}
	cmd := m.cmd
	m.cmd = nil // nil before signalling so goroutine sees non-natural exit
	m.mu.Unlock()

	cmd.Process.Signal(os.Interrupt)
	// No Wait() here — the goroutine owns the Wait() call for this process.
}

func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cmd != nil
}
