/* Bismillah */

// A counter: a label and two buttons.
package main

import (
	"strconv"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	count = 0
	label *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Counter", Width: 360, Height: 260})
}

func start(page *ui.Object) {
	if !ui.IsDark() { // a dark window has its own ground
		page.Set("color", "#F4F5F7")
	}
	ui.VGroup(ui.Props{Gap: 16}, func() {
		label = ui.Label(ui.Props{Text: "0", FontSize: 48})
		ui.HGroup(ui.Props{Hug: true, Gap: 10}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Minus", OnClick: func() { add(-1) }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Plus", Tone: "cyan", OnClick: func() { add(1) }})
		})
	})
}

func add(step int) {
	count += step
	label.SetText(strconv.Itoa(count))
}
