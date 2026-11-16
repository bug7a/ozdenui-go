/* Bismillah */

// A calculator: the keys stand in a GridGroup. The keys of the keyboard work too: 12 + 7, then Enter.
package main

import (
	"strconv"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor   = "#F4F5F7"
	screenColor = "#FFFFFF"
	errorText   = "Error"
	mostDigits  = 12
)

var keys = []string{"C", "÷", "7", "8", "9", "×", "4", "5", "6", "−", "1", "2", "3", "+", "0", ".", "="}

var wideKeys = map[string]int{"C": 3, "0": 2} // how many cells a key takes

var keyboardNames = map[string]string{"*": "×", "/": "÷", "-": "−", "Enter": "=", "Escape": "C", "c": "C"}

var operations = map[string]func(a, b float64) float64{
	"+": func(a, b float64) float64 { return a + b },
	"−": func(a, b float64) float64 { return a - b },
	"×": func(a, b float64) float64 { return a * b },
	"÷": func(a, b float64) float64 { return a / b },
}

var (
	typed     = "0"
	left      float64
	operation = ""
	isNew     = true
	screen    *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Calculator", Width: 320, Height: 410})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 10, Padding: 16, Align: "left top"}, func() {
		screen = ui.Label(ui.Props{Text: "0", Width: "100%", FontSize: 34, Bold: true, TextAlign: "right", Color: screenColor,
			Round: 10, Padding: []int{14, 8}, Shrink: 0})
		wide := map[*ui.Object]int{}
		grid := ui.GridGroup(ui.GridGroupProps{Columns: 4, Gap: 8}, func() {
			for _, key := range keys {
				tone := "white"
				if key == "=" {
					tone = "cyan"
				}
				button := ui.SmartButton(ui.SmartButtonProps{LabelText: key, Width: "100%", Tone: tone, OnClick: func() { press(key) }})
				if span, ok := wideKeys[key]; ok {
					wide[button] = span
				}
			}
		})
		// WHY after the group: the group is not in its variable yet while its function makes the keys
		for button, span := range wide {
			grid.Call("place", button, ui.P{"columnSpan": span})
		}
	})
	page.On("keydown", func(obj *ui.Object, event ui.Event) { pressOnKeyboard(event.Key) })
}

func pressOnKeyboard(name string) {
	key, ok := keyboardNames[name]
	if !ok {
		key = name
	}
	for _, known := range keys {
		if key == known {
			press(key)
			return
		}
	}
}

func press(key string) {
	if key == "C" {
		clear()
	} else if key == "=" {
		calculate()
	} else if _, ok := operations[key]; ok {
		chooseOperation(key)
	} else {
		write(key)
	}
}

func clear() {
	typed, left, operation, isNew = "0", 0, "", true
	screen.SetText("0")
}

func write(key string) {
	text := typed
	if isNew {
		text = ""
	}
	if key == "." && strings.Contains(text, ".") {
		return
	}
	if key == "." && text == "" {
		text = "0"
	}
	text = strings.TrimLeft(text+key, "0")
	if text == "" {
		text = "0"
	}
	if strings.HasPrefix(text, ".") {
		text = "0" + text
	}
	if len(text) > mostDigits {
		text = text[:mostDigits]
	}
	typed, isNew = text, false
	screen.SetText(typed)
}

func chooseOperation(key string) {
	if operation != "" && !isNew {
		calculate()
	}
	left, _ = strconv.ParseFloat(typed, 64)
	operation, isNew = key, true
}

func calculate() {
	if operation == "" {
		return
	}
	right, _ := strconv.ParseFloat(typed, 64)
	if operation == "÷" && right == 0 {
		clear()
		screen.SetText(errorText)
		return
	}
	typed = strconv.FormatFloat(operations[operation](left, right), 'g', 10, 64)
	operation, isNew = "", true
	screen.SetText(typed)
}
