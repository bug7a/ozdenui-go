# 2. Objects

Five plain objects make every window: `Box`, `Label`, `Button`, `Input` and `Icon`. The components of chapter 5 are
made of them. All five take one value: `ui.Props{...}`.

| Object | What it is | Its own fields |
|---|---|---|
| `ui.Box` | a rectangle; the only object that holds other objects | `Color`, `Border`, `Round`, `ScrollY` |
| `ui.Label` | a text | `Text`, `Html`, `FontSize`, `Bold`, `TextColor`, `TextAlign` |
| `ui.Button` | a plain button of the system (the nicer one is `ui.SmartButton`) | `Text` |
| `ui.Input` | a one-line text field | `Text`, `Placeholder`, `OnEnter` |
| `ui.Icon` | a picture | `Picture` (a path, the bytes of a PNG or a JPEG, an `image.Image`, an address) |

## A place and a size

An object outside a group has a place of its own: `Left`, `Top`, `Width`, `Height`, in px:

```go
ui.Box(ui.Props{Left: 20, Top: 20, Width: 100, Height: 100, Color: "orangered"})
```

A size may also be a text: `Width: "100%"` (of the box it is in), `Height: "auto"` (as high as its content). An object
without a place starts at the top left corner of its box.

## Values

A value is a field of `ui.Props` when the object is made, and a name later:

```go
title := ui.Label(ui.Props{Text: "Report", FontSize: 24, Bold: true, TextColor: "#2B2B2B"})
title.SetText("Monthly report")                                   // changes it on the screen
title.Set("text", "Yearly report", "fontSize", 28, "bold", false) // more than one at once
fmt.Println(title.Text())                                         // "Yearly report": read in Go, no question to the page
```

- **Names:** a field is `FontSize`; the same value in `Set()` and `Get()` is `"fontSize"`, as basic.js writes it (the
  name of the field with a small first letter). A wrong name is not lost: the terminal says
  `ozdenui: main.go:12: Label: "fontSzie" is not a value of Label (did you mean "fontSize"?)`.
- **Numbers and flags:** a size in px is a number; a flag is `true` / `false`.
- **Colors:** any CSS color: `"red"`, `"#2A78D6"`, `"rgba(0, 0, 0, 0.5)"`.
- **A field that is not written keeps the default.** A zero that is written is a value (`Gap: 0`). The one exception
  is a text: `""` means "not given", so an empty text on purpose is `ui.Empty` (chapter 5).
- **Reading:** `obj.Get("fontSize")` is the value as `any`; `obj.Str()`, `obj.Num()`, `obj.Int()` and `obj.Bool()`
  give it in a type: `title.Num("fontSize")` is `28`.

The values every object has:

| Field | Example |
|---|---|
| `Color` | the background: `Color: "#FFFFFF"`, `Color: "transparent"` |
| `Border`, `BorderColor` | `Border: 1, BorderColor: "#DDD"` |
| `Round` | the corners: `Round: 8`; `Round: 100` makes a circle of a square |
| `Opacity` | `Opacity: 0.5` |
| `Visible` | `Visible: false` hides it; in a group the others close the gap |
| `Clickable` | an object with an event takes clicks by itself; `Clickable: false` lets them pass through |
| `Cursor` | `Cursor: "pointer"` |
| `Css` | anything else of CSS: `Css: ui.P{"letterSpacing": "1px"}` |
| `Extra` | any other value of basic.js by its name: `Extra: ui.P{"textSize": 12}` |

`ui.P` is `map[string]any`. You will meet it again: a style, an item of a list, a row of a table.

### A value you did not give

`title.Text()` answers from Go because your code gave it. A value the page decides, like the width of a label without
a width, is asked with `Measure()`:

```go
width := title.Measure("width") // the width on the screen now, in px (a float64 in an any)
```

`Get()` of a value that was never given is `nil`, and the terminal names `Measure()` for it. `Measure()` waits for the
page, so use it when you need it, not in a loop. The page itself is different: `page.Width()` and `page.Height()` are
always current.

### Text and HTML

The `Text` of a `Label` or a `Button` is shown exactly as it is written, so a text typed by a user or read from a file
goes in as it is (a `<b>` stays `<b>`). HTML goes in as `Html`:

```go
note.SetText(typedText)
title.Set("html", "<b>3</b> new")
```

In basic.js `text` is HTML; in Go it is the plain text of basic.js. The texts of the components are shown as they
are written too (`ui.SmartButtonProps{LabelText: typed}`); `AllowHtml: true` makes them HTML.

## Objects inside a Box

The objects made in the function given to a `Box` go inside it. Their place is counted from the box's corner:

```go
package main

import (
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

// A white circle with a check, as an address of a picture: no file is needed.
const check = `data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'>` +
	`<circle cx='12' cy='12' r='11' fill='white'/><path d='M7.5 12.5l3 3 6-7' fill='none' stroke='%232B2B2B' stroke-width='2.4'/></svg>`

var (
	label *ui.Object
	input *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A card", Width: 360, Height: 220})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.Box(ui.Props{Left: 20, Top: 20, Width: 320, Height: 110, Color: "#5BBAC1", Round: 12}, func() {
		ui.Icon(ui.Props{Left: 20, Top: 25, Width: 60, Height: 60, Picture: check})
		label = ui.Label(ui.Props{Left: 100, Top: 36, Width: 200, Height: 38, Text: "Hello", FontSize: 26, Bold: true, TextColor: "white"})
	})
	input = ui.Input(ui.Props{Left: 20, Top: 150, Width: 200, Height: 40, Placeholder: "Write a text", OnEnter: change})
	ui.Button(ui.Props{Left: 230, Top: 150, Width: 110, Height: 40, Text: "Change", OnClick: change})
}

func change() {
	typed := strings.TrimSpace(input.Text())
	if typed == "" {
		typed = "Hello"
	}
	label.SetText(typed)
}
```

- `ui.Box(props, func() { ... })`: the function runs at once, and what it makes goes into the box. The function is
  optional: `ui.Box(props)` is an empty box.
- **Only a `Box` (and the groups, which are boxes) holds objects.** A `Label` or a `Button` has no function to give.
  To put something over a button, make it next to the button in the same box, at the same place.
- `Picture` takes a path of this computer (`"photo.png"`), the bytes of a PNG or a JPEG, an `image.Image` made by
  your program, or an address. The `#` of a color in an address of an SVG is written `%23`.
- `OnEnter: change`: Enter in the field does what the button does.
- `examples/basics` is this card, with a color that changes at every press.

## Hiding, removing, adding later

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	note  *ui.Object
	info  *ui.Object
	shown = true
)

func main() {
	ui.Run(start, ui.Window{Title: "Hide and remove", Width: 420, Height: 300})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14}, func() {
		note = ui.Label(ui.Props{Text: "A note to hide", FontSize: 22, Color: "white", Round: 8, Padding: []int{14, 8}})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Hide / show", OnClick: toggle})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Measure", OnClick: measure})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Remove", Tone: "orange", OnClick: remove})
		})
		info = ui.Label(ui.Props{Text: "The width is asked of the page.", FontSize: 14, TextColor: "#5A6270"})
	})
}

func toggle() {
	if note == nil {
		return
	}
	shown = !shown
	if shown {
		note.Show()
	} else {
		note.Hide()
	}
}

func measure() {
	if note == nil {
		return
	}
	info.SetText(fmt.Sprintf("The note is %v px wide.", note.Measure("width")))
}

func remove() {
	if note == nil {
		return
	}
	note.Remove()
	note = nil // a removed object is not used again
	info.SetText("The note is gone.")
}
```

- `Hide()` and `Show()` change `Visible`; a hidden object keeps its values, and `Show()` brings it back.
- `Remove()` removes the object and every object inside it, for good: make a new one to show it again. A removed object
  is not used again: the program above forgets it (`note = nil`).
- `box.Add(func() { ... })` makes more objects in a box or a group later, in an event or a timer:

```go
list.Add(func() {
	ui.Label(ui.Props{Text: "One more"})
})
```

## When a box has more than it can show

`ScrollY: true` makes a box scroll its content up and down (`ScrollX: true` sideways). In a group, the `ScrollArea`
of chapter 3 does the same and takes the room it is given.

## Next

[Chapter 3: Layout](03-layout.md): groups place the objects for you.
