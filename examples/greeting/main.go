/* Bismillah */

// A TextField read in Go: name.Text() is always the text on the page. Enter or the button greets.
package main

import (
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const pageColor = "#F4F5F7"

var (
	name   *ui.Object
	answer *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Greeting", Width: 420, Height: 320})
}

func start(page *ui.Object) {
	if !ui.IsDark() { // a dark window has its own ground
		page.Set("color", pageColor)
	}
	ui.VGroup(ui.Props{Gap: 14, Padding: 40}, func() {
		name = ui.TextField(ui.TextFieldProps{TitleText: "What is your name?", Placeholder: "Ada", Width: 300, OnEnter: greet})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Greet", Tone: "cyan", OnClick: greet})
		answer = ui.Label(ui.Props{FontSize: 18})
	})
}

func greet() {
	typed := strings.TrimSpace(name.Text())
	if typed == "" {
		answer.SetText("Write your name first.")
		return
	}
	answer.SetText("Hello, " + typed + "!")
}
