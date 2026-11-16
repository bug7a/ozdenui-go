# OzdenUI for Go — The Handbook

Build a desktop window in plain Go: objects, groups and ready components. Only Go: no JavaScript, no HTML, no CSS,
no build step of a front end. This handbook teaches it step by step, for people who know Go and have never written a
web page.

*Türkçe: [OzdenUI for Go — El Kitabı](turkce/README.md)*

```go
package main

import ui "github.com/bug7a/ozdenui-go"

func main() {
	ui.Run(start, ui.Window{Title: "Hello", Width: 360, Height: 200})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Gap: 12}, func() {
		ui.Label(ui.Props{Text: "Hello from Go", FontSize: 24, Bold: true})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Close", OnClick: ui.Close})
	})
}
```

## The chapters

| | Chapter | You learn |
|---|---|---|
| 1 | [The first window](01-first-window.md) | install, `ui.Run`, `start(page)`, `ui.Window`, how the window and your code talk |
| 2 | [Objects](02-objects.md) | `Box`, `Label`, `Button`, `Input`, `Icon`; `ui.Props`; values you set and read; removing |
| 3 | [Layout](03-layout.md) | `VGroup`, `HGroup`, `GridGroup`, `Fill: true`, `ScrollArea`, a place of its own |
| 4 | [Events and time](04-events-and-time.md) | clicks and keys, `ui.Event`, `ui.Every` and `ui.After`, a goroutine and `ui.Do` |
| 5 | [Components](05-components.md) | the ready parts: buttons, fields, choices, lists, tables, tabs, bars |
| 6 | [Dialogs, messages and files](06-dialogs-and-files.md) | `ui.Confirm()`, toasts, bubbles, the file dialogs of the system |
| 7 | [Looks](07-looks.md) | dark windows, the minimal look in your color, the colors of one component |
| 8 | [Data and pictures](08-data-and-pictures.md) | charts, big tables, logs, Markdown, JSON, code, pictures |
| 9 | [Your own component](09-your-own-component.md) | a type of yours made of objects, shared as a package |
| 10 | [Going further](10-going-further.md) | basic.js from Go, speed, goroutines, limits, finding a mistake |

Read 1 to 4 in order: they are the whole idea. Then take the chapters you need.

## How to read it

- Every program that starts with `package main` is complete: copy it into the `main.go` of a module that has
  `go get github.com/bug7a/ozdenui-go`, and run it with `go run .`. A check of the module builds and draws every one
  of them, so they stay true.
- The examples of the module (the `examples/` folder of its source) are longer programs: a to-do list, a calculator,
  a CSV viewer, an editor, a small game. The chapters name them where they fit. From a clone of the repository:
  `go run ./examples/todo`.
- A name in `code` is a name in Go. The page of the window is made with [basic.js](https://bug7a.github.io/basic.js/),
  a small JavaScript library; you never need to write JavaScript, but chapter 10 shows how to reach it.
- Every function and every field is in `go doc github.com/bug7a/ozdenui-go` and on
  [pkg.go.dev](https://pkg.go.dev/github.com/bug7a/ozdenui-go). The site: https://bug7a.github.io/ozdenui-go/

## In one minute

- `ui.Run(start, ui.Window{...})` opens a window. `start(page)` makes its objects; `page` is the window's ground.
- An object is made where it is written: inside the function given to a box or a group, it goes into that box or group.
- A group places its objects: `ui.VGroup` one under the other, `ui.HGroup` side by side.
- A value is a field when the object is made (`ui.Props{FontSize: 20}`) and a name later (`label.Set("fontSize", 24)`).
- An event is a function: `OnClick: save`. It may take nothing, the object, or the object and the event.
- Components keep their state in Go: `box.Checked()`, `selectBox.Value()`, `field.Text()`. Reading it asks nothing of the page.
- Your code runs on one goroutine: an event, a timer and `start()` never run at the same time.
