package main

import (
	"fmt"

	"github.com/getlantern/systray"
)

func runTray(m *Manager) {
	systray.Run(func() { onReady(m) }, func() { m.Stop() })
}

func onReady(m *Manager) {
	systray.SetTemplateIcon(coldIconPNG, coldIconPNG)
	systray.SetTitle("")
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

	setOn := func() {
		systray.SetTemplateIcon(hotIconPNG, hotIconPNG)
		mStatus.SetTitle("● Caffeinated")
		mToggle.SetTitle("Turn Off")
	}

	setOff := func() {
		systray.SetTemplateIcon(coldIconPNG, coldIconPNG)
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
					setOff()
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
		setOff()
		return
	}
	setOff()
	setOn()
}
