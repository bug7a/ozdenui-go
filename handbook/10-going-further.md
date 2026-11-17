# 10. Going further

## Finding a mistake

- **A mistake of Go** (a wrong field, a wrong type) does not build: a field of `ui.Props` or of a `...Props` type that
  does not exist is a compile error, and your editor marks it before the program runs.
- **A name the compiler can not check** (in `Set()`, `Extra`, a `ui.P` of a component) is told once, with the line of
  your program:

  ```
  ozdenui: main.go:17: Label: "fontSzie" is not a value of Label (did you mean "fontSize"?)
  ozdenui: main.go:24: TextField: "warningTxt" is not a value of TextField (did you mean "warningText"?)
  ```

- **A method that does not exist** (`obj.Call("setValu", 3)`) and **a value that was never given**
  (`label.Get("width")` of an `"auto"` label) are told the same way; the second one names `Measure()`.
- **A panic in an event or a timer** does not close the window with your work in it: the terminal writes
  `ozdenui: a task failed: ...` and the lines of your program that led to it, and the window keeps working.
- **An object changed from a goroutine of your own** panics with the fix in its message: `ui.Do(func() { ... })`.
- **A message of the page** (a component that could not do what it was asked) comes as `ozdenui (page ...): ...`.
- `ui.Window{Debug: true}` gives the window the developer tools of the web view: right click > Inspect.

So: run the program, and **read the terminal**.

## basic.js from Go

Everything of basic.js and of the kit is reachable, also what has no Go name yet:

```go
gauge.Call("setValue", 80)            // a method of the object on the page; not waited for
text := editor.Ask("getSelectedText") // a method, and its result
width := label.Measure("width")       // a value as the page has it now (a size of "auto")
field.Part("input").Call("select")    // a box inside a component

// Any JavaScript; label.JS() is the object.
value := ui.Evaluate(label.JS() + ".elem.offsetWidth")

// Any component of the kit by its name, and a param that has no field, by its name in the kit.
gauge = ui.Component("Gauge", ui.P{"value": 72, "width": 200})
ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Extra: ui.P{"ariaLabel": "Save the file"}})
```

- `Call()` goes with the other commands of the task. `Ask()`, `Measure()` and `Evaluate()` wait for the page, so keep
  them out of loops.
- What comes back from the page is plain: a number is a `float64`, a flag of the kit is `1` / `0`, a list is `[]any`,
  an object `map[string]any`. `ui.Number(v)`, `ui.Text(v)` and `ui.IsOn(v)` read them.
- `go doc` of a component lists its methods; what each takes is in its file of the kit:
  [web/comp](https://github.com/bug7a/ozdenui-go/tree/main/web/comp). The params keep their names in camelCase.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	title *ui.Object
	line  *ui.Object
	name  *ui.Object
	info  *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Measured", Width: 480, Height: 340})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14, Padding: 24, Align: "left top"}, func() {
		title = ui.Label(ui.Props{Text: "Measured by the page", FontSize: 26, Bold: true})
		line = ui.Box(ui.Props{Width: 20, Height: 4, Color: "#2A78D6", Round: 2})
		name = ui.TextField(ui.TextFieldProps{TitleText: "A new title", Value: "Hello", Width: 300, OnEnter: retitle})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Use it", Tone: "cyan", OnClick: retitle})
		info = ui.Label(ui.Props{Text: "-", FontSize: 14, TextColor: "#5A6270"})
	})
	line.Call("setMotion", "width 0.3s") // the page moves the width; Go gives only the end
	underline()
}

func retitle() {
	title.SetText(name.Text())
	underline()
	name.Part("input").Call("select") // the text of the field is selected, ready for the next one
}

// underline makes the line as wide as the title: a size of "auto" is known only by the page.
func underline() {
	width := ui.Number(title.Measure("width"))
	line.Set("width", width)
	pixels := ui.Number(ui.Evaluate("window.devicePixelRatio"))
	info.SetText(fmt.Sprintf("The title is %.0f px wide; the pixel ratio of this screen is %.0f.", width, pixels))
}
```

## Goroutines, again

| | may change objects | may ask the page (`Ask`, `Measure`, `Evaluate`) and show a dialog |
|---|---|---|
| `start()`, events, timers (the goroutine of the window) | yes | yes |
| a goroutine of yours | inside `ui.Do(func() { ... })` | inside `ui.Do(...)`, but the goroutine does not wait for it |

Any goroutine may call `ui.Do()`, `ui.Close()`, `ui.SetTheme()` and write to `logView.Writer(level)`.

## Speed

What makes a window slow is almost always one of two things:

- **A question to the page in a loop**: every `Measure()`, `Ask()` or `Evaluate()` waits for an answer. Reading what
  your program gave or what a component told (`obj.Get()`, `field.Text()`, `box.Checked()`) asks nothing.
- **A slow job on the goroutine of the window**: the window does not answer until it ends. Put it in a goroutine
  (chapter 4).

All the commands of one task go to the page together when it ends, so making 500 objects in a loop is one message,
not 500. A long task that must stay on the goroutine of the window can show its progress with `ui.Flush()`: the
commands made until then go to the page at once.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	limit = 10_000_000
	parts = 20
)

var (
	bar    *ui.Object
	result *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Primes", Width: 420, Height: 260})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 16, Padding: 24}, func() {
		bar = ui.ProgressBar(ui.ProgressBarProps{TitleText: "Numbers checked", MaxValue: parts, Width: 320})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Count the primes", Tone: "cyan", OnClick: countPrimes})
		result = ui.Label(ui.Props{Text: fmt.Sprintf("Primes under %d: ?", limit), FontSize: 15})
	})
}

// countPrimes runs on the goroutine of the window: Flush shows each part as soon as it is done.
func countPrimes() {
	count := 0
	step := limit / parts
	for part := 0; part < parts; part++ {
		for number := part * step; number < (part+1)*step; number++ {
			if isPrime(number) {
				count++
			}
		}
		bar.SetValue(part + 1)
		ui.Flush()
	}
	result.SetText(fmt.Sprintf("Primes under %d: %d", limit, count))
}

func isPrime(number int) bool {
	if number < 2 {
		return false
	}
	for divisor := 2; divisor*divisor <= number; divisor++ {
		if number%divisor == 0 {
			return false
		}
	}
	return true
}
```

The bar moves, but a click waits until the count ends. A job of seconds belongs in a goroutine with `ui.Do()`.

## A picture of the window, a check without a person

```
OZDENUI_LOOK=app.png go run .                               # runs the program, saves app.png, ends; no window
OZDENUI_THEME=dark OZDENUI_LOOK=dark.png go run .           # the dark window, without a change of the code
OZDENUI_LOOK_WAIT=3 OZDENUI_LOOK=app.png go run .           # 3 s after the first drawing (default 1.5)
OZDENUI_WINDOW=headless go run .                            # a hidden window, until the program calls ui.Close()
OZDENUI_WINDOW=browser go run .                             # the app window of Edge or Chrome, also on macOS
OZDENUI_WINDOW=tab go run .                                 # a tab of the default browser
```

- A picture needs Chrome, Edge or Chromium on the computer. A program that runs without an error can still have a
  part covered, cut or squeezed: the picture shows it.
- A check of a program closes it from the code: `ui.After(3000, ui.Close)` as the first line of `start()`. It must
  not reach a dialog: a dialog waits for a person.

## With an AI assistant

OzdenUI is new, so an AI assistant (Claude Code, Codex, Cursor, Copilot...) does not know it yet and guesses from the
toolkits it knows. One command writes the rules of the module into the folder of your project, where the assistant
reads them by itself:

```
go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents              # into AGENTS.md
go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents CLAUDE.md    # into the file Claude Code reads
go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents --print      # only prints them
```

- The rules are short: how a program is made, the layout, the events, the components, the mistakes to avoid, and
  `OZDENUI_LOOK`, so the assistant looks at the window it made.
- A file that is there keeps your lines: the rules go to its end, between two marks. Run the command again after an
  upgrade of the module: only that part changes.

## One file to give away

```
go build -o notes .                       # one file: the page of the window and the kit are inside it
GOOS=windows GOARCH=amd64 go build .      # a program for Windows, built on another computer
```

The page, basic.js and the components are embedded (`go:embed`), so the built program needs no other file. A small
program is about 15 MB. On macOS the build needs cgo (the web view of the system); Windows and Linux need no C compiler.

## Limits of v0.2

- macOS and Windows are tested by hand. Linux opens the app window of Edge or Chrome (with none of them, a tab of the
  default browser); it is checked by the tests of the repository, not yet by hand.
- On Windows the window is sharp on the scale of the main screen; on a second screen with another scale the system
  stretches it.
- One window for a program.
- A callback that must answer at once can not be a Go function, because the answer comes later: a validator, a
  function that builds the objects of a row (`createItemView`). Format the values of a table in Go; the numbers of a
  `GraphBox` or a `Stepper` take a format text (chapter 8).
- A file dropped on a `FilePicker` gives its name, not its path; a file chosen with a click has its path (chapter 6).
- A bubble of a `Tooltip` shows a text, not objects of your own.
- Fifteen components change their theme while they are on the page; the others take it when they are made (chapter 7).
- No hot reload: `go run .` again.

## Where to go next

- **The examples** of the module, from the smallest: `counter`, `greeting`, `basics`, `notes`, `todo`, `form`,
  `settings`, `calculator`, `dialogs`, `clock`, `game`, `charts`, `csv-viewer`, `json-viewer`, `image-viewer`,
  `editor`, `markdown-editor`, `job-log`, `chat`, `gallery` (the main components), `minimal-theme`.
- **The README** of the module: every component in one table.
- **`go doc github.com/bug7a/ozdenui-go SelectBoxProps`**: every param of a component with its default;
  `go doc github.com/bug7a/ozdenui-go SelectBox`: its looks and its methods.
- **[basic.js](https://bug7a.github.io/basic.js/)** and **[the components](https://bug7a.github.io/js-components/)**:
  the library under it, with a live example of every component.

[Back to the start](README.md)
