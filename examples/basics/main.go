/* Bismillah */

// The basic objects: a Box, an Icon, a Label, an Input and a Button, each in a place of its own.
// Write a text and press the button (or Enter): the card shows the text and takes the next color.
package main

import (
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const pageColor = "#F4F5F7"

var cardColors = []string{"#5BBAC1", "#E0A038", "#E8796D", "#65A293"}

// A white circle with a check: a drawing as a text, so no picture file is needed.
const checkSvg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><circle cx="12" cy="12" r="11" fill="white"/>` +
	`<path d="M7.5 12.5l3 3 6-7" fill="none" stroke="#2B2B2B" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"/></svg>`

// The address of the drawing for the page.
var picture = ui.Picture([]byte(checkSvg))

var (
	colorNumber = 0
	card        *ui.Object
	label       *ui.Object
	input       *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "The basic objects", Width: 360, Height: 260})
}

func start(page *ui.Object) {
	if !ui.IsDark() { // a dark window has its own ground
		page.Set("color", pageColor)
	}

	// An object has a place: Left, Top, Width, Height. The objects made in the function of a Box are inside it.
	card = ui.Box(ui.Props{Left: 20, Top: 20, Width: 320, Height: 110, Color: cardColors[0], Round: 12}, func() {
		ui.Icon(ui.Props{Left: 20, Top: 25, Width: 60, Height: 60, Picture: picture})
		label = ui.Label(ui.Props{Left: 100, Top: 36, Width: 200, Height: 38, Text: "Hello", FontSize: 26, Bold: true, TextColor: "white"})
	})

	// An event is a function: OnClick runs at a click, OnEnter at Enter.
	input = ui.Input(ui.Props{Left: 20, Top: 150, Width: 200, Height: 40, Placeholder: "Write a text", OnEnter: change})
	ui.Button(ui.Props{Left: 230, Top: 150, Width: 110, Height: 40, Text: "Change", OnClick: change})
}

func change() {
	typed := strings.TrimSpace(input.Text())
	if typed == "" {
		typed = "Hello"
	}
	label.SetText(typed)
	colorNumber = (colorNumber + 1) % len(cardColors)
	card.Set("color", cardColors[colorNumber])
}
