# 4. Events and time

## An event is a function

```go
ui.Button(ui.Props{Text: "Save", OnClick: save}) // when it is made
button.Set("onClick", save)                      // later: it takes the place of the one before
remove := button.On("click", logClick)           // one more handler; remove() takes it away
```

The fields of `ui.Props` for events: `OnClick`, `OnInput`, `OnEnter`, `OnKeyDown`, `OnMouseEnter`, `OnMouseLeave`,
`OnFocus`, `OnBlur`. Any other event of the page works with `On()`, written as the page writes it: `"dblclick"`,
`"mouseenter"`, `"keyup"`, `"contextmenu"`.

### The handler takes what it needs

```go
func save()                                  // nothing
func save(button *ui.Object)                 // the object
func save(button *ui.Object, event ui.Event) // the object and the event
```

The field of an event is `any`, so all three fit. A value that is not a function panics with the line of your
program. A method value works too: `OnClick: app.save`.

**In a loop, every closure has its own variable** (Go 1.22 and newer), so the obvious code is right:

```go
for number := 1; number <= 5; number++ {
	ui.SmartButton(ui.SmartButtonProps{LabelText: strconv.Itoa(number), OnClick: func() { pick(number) }})
}
```

`ui.Event` has what the event knows:

| Field | When |
|---|---|
| `Key`, `Code` | a key: `"Enter"`, `"a"`, `"ArrowLeft"`; `Code` is the place on the keyboard (`"KeyA"`) |
| `Shift`, `Ctrl`, `Alt`, `Meta` | a key: `true` while held (`Meta` is Cmd on a Mac) |
| `X`, `Y` | the mouse, in px from the window's corner |
| `Text` | an `Input`: its text |
| `Name` | the name of the event: `"click"`, `"keydown"` |

```go
package main

import (
	"fmt"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	spotColor  = "#5BBAC1"
	hoverColor = "#E0A038"
)

var (
	pressed *ui.Object
	spot    *ui.Object
	where   *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Events", Width: 420, Height: 300})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14}, func() {
		pressed = ui.Label(ui.Props{Text: "Press a key", FontSize: 28, Bold: true})
		spot = ui.Box(ui.Props{Width: 160, Height: 60, Color: spotColor, Round: 10})
		where = ui.Label(ui.Props{Text: "Click the box", FontSize: 14, TextColor: "#5A6270"})
	})
	page.On("keydown", showKey)
	spot.On("mouseenter", func() { spot.Set("color", hoverColor) })
	spot.On("mouseleave", func() { spot.Set("color", spotColor) })
	spot.On("click", func(box *ui.Object, event ui.Event) {
		where.SetText(fmt.Sprintf("Clicked at %.0f, %.0f", event.X, event.Y))
	})
}

func showKey(page *ui.Object, event ui.Event) {
	held := []string{}
	if event.Ctrl {
		held = append(held, "Ctrl")
	}
	if event.Alt {
		held = append(held, "Alt")
	}
	if event.Shift {
		held = append(held, "Shift")
	}
	if event.Meta {
		held = append(held, "Meta")
	}
	pressed.SetText(strings.Join(append(held, event.Key), " + "))
}
```

`page.On("keydown", ...)` hears every key of the window; a field that has the focus hears its own keys too.

## One goroutine: no locks

`start()`, every event and every timer run **one at a time, on one goroutine** (the app goroutine). Two clicks never
run at the same time, and a value you read is the value of now. You need no locks, also for the variables of your
program that only these functions use.

The other side of it: while one function runs, the window waits. A function that takes 3 seconds freezes the window
for 3 seconds. So **never `time.Sleep()` in an event**: wait with a timer, and give a slow job a goroutine (below).

## Time: Every() and After()

```go
package main

import (
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	clock     *ui.Object
	stopClock func()
)

func main() {
	ui.Run(start, ui.Window{Title: "Clock", Width: 360, Height: 220})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Gap: 12}, func() {
		clock = ui.Label(ui.Props{Text: time.Now().Format("15:04:05"), FontSize: 48, Bold: true})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Stop the clock", OnClick: func() { stopClock() }})
	})
	stopClock = ui.Every(1000, tick)                      // again and again, every 1000 ms
	ui.After(5000, func() { ui.SetTitle("Clock — 5 s") }) // once
}

func tick() {
	clock.SetText(time.Now().Format("15:04:05"))
}
```

- `ui.Every(ms, fn)` returns a function that stops it; `ui.After(ms, fn)` one that cancels it.
- A timer is a task too: it runs on the app goroutine, between the clicks, so it may change any object. A tick that
  comes while the one before still runs is not piled up.
- `OnClick: func() { stopClock() }` and not `OnClick: stopClock`: the variable is still `nil` when the button is made,
  and the closure reads it at the click.
- `examples/clock` adds a stopwatch.

## A slow job: a goroutine and ui.Do()

```go
package main

import (
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const steps = 100

var (
	bar    *ui.Object
	button *ui.Object
	result *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A slow job", Width: 420, Height: 260})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14}, func() {
		bar = ui.ProgressBar(ui.ProgressBarProps{TitleText: "Working", Value: 0, MaxValue: steps, Width: 320})
		button = ui.SmartButton(ui.SmartButtonProps{LabelText: "Start", Tone: "cyan", OnClick: begin})
		result = ui.Label(ui.Props{FontSize: 16})
	})
}

func begin() {
	button.Set("waiting", true) // a clock on the button: it takes no clicks now
	result.SetText("")
	go work()
}

func work() {
	for step := 1; step <= steps; step++ {
		time.Sleep(30 * time.Millisecond) // the slow part: a big file, a download, a model
		ui.Do(func() { bar.SetValue(step) })
	}
	ui.Do(func() {
		button.Set("waiting", false)
		result.SetText("Done.")
	})
}
```

- `go work()`: the job runs on a goroutine of its own; the window keeps working.
- **Objects are changed only on the app goroutine.** A goroutine of yours gives its changes to it with
  `ui.Do(func() { ... })`: the function runs there, after the task that runs now. Any goroutine may call `ui.Do`: a
  server, a socket, a watcher of files.
- An object changed from another goroutine panics, and the message gives the fix:
  `ozdenui: ... works only in start(), an event or a timer. From another goroutine: ozdenui.Do(func() { ... })`.
  Nothing breaks without a word.
- Give the job what it needs before it starts (a text of a field, a path): reading an object is also the business of
  the app goroutine.
- `ui.Close()` works from any goroutine.
- `examples/job-log` copies files this way and writes what it does into a `LogView`.

### A long task that shows its progress

The commands of a task go to the page when the task ends. A task that must stay on the app goroutine for a while
(many objects at once) can send what it made so far with `ui.Flush()`.

## The order of events

Events reach Go in the order they happened, also when many come at once (fast typing), and each one runs as its own
task. A fast and continuous one (`obj.On("mousemove", fn)`) sends every move to Go. A component that is dragged (a
`Slider`, a `Splitter`) does its moves on the page and tells you the result (`OnChange`, `OnChangeEnd`).

## Next

[Chapter 5: Components](05-components.md): the ready parts of a window.
