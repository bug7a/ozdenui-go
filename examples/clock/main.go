/* Bismillah */

// A clock and a stopwatch: ui.Every() runs a function again and again (the loop() of basic.js).
package main

import (
	"fmt"
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor       = "#1E2230"
	textLight       = "#F4F5F7"
	textSoft        = "#9AA3B5"
	stopwatchStepMs = 100
)

var (
	startedAt     time.Time // zero: the stopwatch is not running
	total         time.Duration
	stopStopwatch func()

	clock       *ui.Object
	watch       *ui.Object
	startButton *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Clock", Width: 460, Height: 340})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 10, Padding: 30}, func() {
		clock = ui.Label(ui.Props{Text: time.Now().Format("15:04:05"), FontSize: 72, Bold: true, TextColor: textLight})
		ui.Label(ui.Props{Text: time.Now().Format("Monday, 02 January 2006"), FontSize: 16, TextColor: textSoft})
		watch = ui.Label(ui.Props{Text: "0.0 s", FontSize: 28, TextColor: textLight})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			startButton = ui.SmartButton(ui.SmartButtonProps{LabelText: "Start", StyleName: "ozdenDark", OnClick: toggle})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Reset", StyleName: "ozdenDark", OnClick: reset})
		})
	})
	ui.Every(1000, showTime)
}

func showTime() {
	clock.SetText(time.Now().Format("15:04:05"))
}

func elapsed() time.Duration {
	if startedAt.IsZero() {
		return total
	}
	return total + time.Since(startedAt)
}

func showWatch() {
	watch.SetText(fmt.Sprintf("%.1f s", elapsed().Seconds()))
}

func toggle() {
	if startedAt.IsZero() {
		startedAt = time.Now()
		stopStopwatch = ui.Every(stopwatchStepMs, showWatch)
		startButton.SetText("Stop")
		return
	}
	total = elapsed()
	startedAt = time.Time{}
	stopStopwatch()
	startButton.SetText("Start")
	showWatch()
}

func reset() {
	total = 0
	if !startedAt.IsZero() {
		startedAt = time.Now()
	}
	showWatch()
}
