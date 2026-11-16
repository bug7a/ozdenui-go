# 1. The first window

## Install

In the folder of your program (a Go module: `go mod init myapp` makes one):

```
go get github.com/bug7a/ozdenui-go
```

It needs Go 1.22 or newer.

- **macOS:** the window is the web view of the system (WKWebView), inside your program. It needs cgo, so the command
  line tools of Xcode must be there (`xcode-select --install`). The first build takes about ten seconds and may print
  warnings of the web view's C code: they are not errors, and the next builds are fast.
- **Windows, Linux:** the window is an app window of Edge or Chrome: no address bar, no tabs. No C compiler is
  needed. On Windows, Edge is always there.

The built program (`go build`) is one file: the page of the window and the components are inside it.

## The smallest window

```go
package main

import ui "github.com/bug7a/ozdenui-go"

func main() {
	ui.Run(func(page *ui.Object) {})
}
```

An empty window of 800 x 600 px, titled "OzdenUI". `go run .` opens it and returns when you close it.

## A window with a counter

```go
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
	ui.Run(start, ui.Window{Title: "Counter", Width: 400, Height: 320})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 16}, func() {
		label = ui.Label(ui.Props{Text: "0", FontSize: 64, Bold: true})
		ui.HGroup(ui.Props{Hug: true, Gap: 12}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "−", OnClick: func() { change(-1) }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "+", Tone: "cyan", OnClick: func() { change(1) }})
		})
	})
}

func change(step int) {
	count += step
	label.SetText(strconv.Itoa(count))
}
```

Read it from the top:

- **`ui.Run(start, ui.Window{...})`** is the window. It opens it, calls `start` when the page is ready, and returns
  when the window is closed. Call it once, at the end of `main()`.
- **`start(page)`** makes the objects. `page` is the ground of the window: it has a color, a width and a height, and
  every object without a box of its own is on it. Objects are made in `start()` or later (in an event, a timer),
  never before `ui.Run()`.
- **`ui.VGroup(ui.Props{Gap: 16}, func() { ... })`** is a group: the objects made inside its function go into it, one
  under the other, 16 px apart. The `HGroup` inside it puts the two buttons side by side; `Hug: true` makes it only as
  large as its buttons.
- **`ui.Label(ui.Props{Text: "0", FontSize: 64, Bold: true})`** is a text. Its values are the fields of `ui.Props`.
- **`OnClick: func() { change(-1) }`** is the event: a click runs the function.
- **`label.SetText(...)`** changes the label. `label.Text()` reads it back, without asking the window.

Every object is a `*ui.Object`. Keep the ones you change later in variables (here `label`); the others need none.

`examples/counter` of the module is a program like this one.

## How it works

The window shows a page that is inside your program. The page holds basic.js and the components of the kit; the
module carries them, so no internet is needed. It is served only to this computer.

```
your Go code ──commands──▶ the page (basic.js objects) ──▶ what you see
        ▲                                   │
        └─────────── events (a click) ◀─────┘
```

- An object in Go is a small proxy: it keeps an id and the values it was given, and sends commands to the page, where
  the real object lives. All the commands of one task (`start`, an event, a timer) go to the page together, when the
  task ends.
- A click on the page calls your Go function on the goroutine of the window (chapter 4).
- Your code is plain Go, with every package you like: files, databases, `net/http`, an AI client.

## The Window

```go
ui.Window{
	Title:    "Notes", // the title of the window ("OzdenUI" when empty)
	Width:    800,     // the size of the page in px (800 x 600 when 0)
	Height:   600,
	Theme:    "light", // "light", "dark" or "auto" (what the system asks): chapter 7
	Style:    "",      // "" or "minimal": chapter 7
	Language: "en",    // the language of its texts: "tr" writes the capital letters of Turkish ("İ")
	Debug:    false,   // true: the developer tools of the web view (right click > Inspect)
}
```

Other functions of the window work in `start()`, an event or a timer (`ui.Close()` also from any goroutine):

```go
ui.SetTitle("Notes — saved")
ui.Close()          // closes the window; Run() returns
ui.Theme()          // "light" or "dark"
ui.SetTheme("dark") // the window again in another theme: chapter 7
```

## A program in more than one function

Most programs keep the objects they change later. Variables at the top of the file are the simplest place:

```go
package main

import (
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	name   *ui.Object
	answer *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Greeting", Width: 420, Height: 320})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
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
```

`name.Text()` is always the text in the field: the page sends it with every key. A `Label` shows a typed text as it
is (a `<b>` stays `<b>`).

A struct works just as well (`app.name = ui.TextField(...)`, then `ui.Run(app.start)`): chapter 9 makes a component
of yours that way. `examples/greeting` is the program above.

## Looking at a window without opening it

```
OZDENUI_LOOK=app.png go run .
```

runs the program as it is, saves a picture of its window as `app.png` and ends; no window opens. It needs Chrome,
Edge or Chromium on the computer. `OZDENUI_THEME=dark` shows the dark window, and `OZDENUI_LOOK_WAIT=3` waits three
seconds before the picture. A program that runs without an error can still have a part covered or cut: a look at the
picture tells it.

## When the page is loaded again

With `Debug: true`, a reload of the page (right click > Reload) makes the objects again: `start()` runs again, the
timers stop, and the terminal says `ozdenui: the page was loaded again: start() runs again`. A program that keeps its
data in Go (a slice, a file) shows it again; nothing else is needed.

## Next

[Chapter 2: Objects](02-objects.md): the five plain objects and their values.
