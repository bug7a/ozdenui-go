/* Bismillah */

// The check of a real window: it writes the size the page sees into $PROBE_REPORT and closes itself when PROBE_CLOSE=1.
package main

import (
	"fmt"
	"os"

	ui "github.com/bug7a/ozdenui-go"
)

func main() {
	ui.Run(start, ui.Window{Title: "Probe", Width: 480, Height: 320})
}

func start(page *ui.Object) {
	if os.Getenv("PROBE_CLOSE") == "1" {
		ui.After(5000, ui.Close)
	}
	ui.VGroup(ui.Props{Gap: 16}, func() {
		ui.Label(ui.Props{Text: "OzdenUI in a real window", FontSize: 28})
		ui.TextField(ui.TextFieldProps{TitleText: "Name", Value: "Ada"})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan"})
	})
	ui.After(1000, func() {
		ui.SetTitle("Probe 2")
		size := ui.Evaluate("innerWidth + ' x ' + innerHeight + ', devicePixelRatio ' + devicePixelRatio + ', ' + navigator.userAgent")
		_ = os.WriteFile(os.Getenv("PROBE_REPORT"), []byte(fmt.Sprint(size)), 0o644)
	})
}
