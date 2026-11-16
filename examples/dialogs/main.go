/* Bismillah */

// The dialogs: each one returns its answer (the code waits; the window keeps working meanwhile).
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor = "#F4F5F7"
	textSoft  = "#5A6270"
)

var answer *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "Dialogs", Width: 640, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 16, Padding: 36, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "Dialogs", FontSize: 24, Bold: true})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Alert", OnClick: showAlert})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Confirm", OnClick: showConfirm})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Prompt", OnClick: showPrompt})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Ask", OnClick: showAsk})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Toast", OnClick: func() { ui.ToastInfo("A short message.", "Info") }})
		})
		ui.Label(ui.Props{Text: "The answer:", FontSize: 13, TextColor: textSoft})
		answer = ui.Label(ui.Props{Text: "-", FontSize: 20})
	})
}

// tell shows a value as Go writes it: a text in quotes, true / false as they are.
func tell(value any) {
	answer.SetText(fmt.Sprintf("%#v", value))
}

func showAlert() {
	ui.Alert("The file was saved.", ui.DialogOptions{Title: "Saved"})
	tell("closed")
}

func showConfirm() {
	tell(ui.Confirm("The old copy will be replaced.", ui.DialogOptions{Title: "Replace it?", ConfirmText: "Replace"}))
}

func showPrompt() {
	name, ok := ui.Prompt("A name for the new folder:", ui.DialogOptions{Title: "New folder", Value: "Photos", Required: true})
	if !ok {
		tell("cancelled")
		return
	}
	tell(name)
}

func showAsk() {
	// Ask takes no title: its message says it all.
	tell(ui.Ask("You have changes that are not saved. Save them?", "Save", "Don't save", "Cancel"))
}
