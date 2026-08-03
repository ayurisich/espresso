package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

var (
	hotFrames = []string{"☕", "♨️☕", "☕♨️"}
	coldIcon  = "🥶☕"
	warnIcon  = "⚠️☕"
)

func runTray(m *Manager) {
	systray.Run(func() { onReady(m) }, func() { m.Stop() })
}

func onReady(m *Manager) {
	systray.SetTitle(coldIcon)
	systray.SetTooltip("Caffeinate Toggle")

	mStatus := systray.AddMenuItem("○ Decaffeinated", "")
	mStatus.Disable()
	systray.AddSeparator()

	mToggle := systray.AddMenuItem("Turn On", "")
	systray.AddSeparator()

	var mHours [5]*systray.MenuItem
	for i := range mHours {
		label := fmt.Sprintf("%d hour", i+1)
		if i > 0 {
			label = fmt.Sprintf("%d hours", i+1)
		}
		mHours[i] = systray.AddMenuItem(label, "")
	}

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	var animDone chan struct{}

	setOn := func() {
		animDone = make(chan struct{})
		done := animDone
		go func() {
			ticker := time.NewTicker(600 * time.Millisecond)
			defer ticker.Stop()
			frame := 0
			for {
				select {
				case <-ticker.C:
					systray.SetTitle(hotFrames[frame%len(hotFrames)])
					frame++
				case <-done:
					return
				}
			}
		}()
		mStatus.SetTitle("● Caffeinated")
		mToggle.SetTitle("Turn Off")
	}

	setOff := func() {
		if animDone != nil {
			close(animDone)
			animDone = nil
		}
		systray.SetTitle(coldIcon)
		mStatus.SetTitle("○ Decaffeinated")
		mToggle.SetTitle("Turn On")
	}

	for {
		select {
		case <-mToggle.ClickedCh:
			if m.IsRunning() {
				m.Stop()
				setOff()
			} else {
				if err := m.Start(); err != nil {
					systray.SetTitle(warnIcon)
				} else {
					setOn()
				}
			}

		case <-mHours[0].ClickedCh:
			startTimed(m, 1, setOn, setOff)
		case <-mHours[1].ClickedCh:
			startTimed(m, 2, setOn, setOff)
		case <-mHours[2].ClickedCh:
			startTimed(m, 3, setOn, setOff)
		case <-mHours[3].ClickedCh:
			startTimed(m, 4, setOn, setOff)
		case <-mHours[4].ClickedCh:
			startTimed(m, 5, setOn, setOff)

		case <-m.Expired:
			setOff()

		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func startTimed(m *Manager, hours int, setOn, setOff func()) {
	if err := m.StartTimed(hours); err != nil {
		systray.SetTitle(warnIcon)
		setOff()
		return
	}
	setOn()
}
