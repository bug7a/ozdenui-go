# 7. Looks

Three levels, from the widest to the smallest: the **theme** of the window (light or dark), the **style** of the window
(the usual look or the minimal one, in your color), and the **look of one component** (its package and its colors).
A look your code gives always wins over the one of the window.

## A dark window

```go
ui.Run(start, ui.Window{Theme: "dark"}) // "light" (the default), "dark", or "auto": what the system asks
```

- Every component takes its dark look; the dialogs and the toasts too.
- The page, a `Label`, a `Box` and an `Input` follow: a dark ground, a light text.
- A `Tone` of a `SmartButton` becomes its dark tone (`"cyan"` → `"cyanDark"`).
- `ui.IsDark()` is true in a dark window (also after `"auto"`), and `ui.Theme()` is `"light"` or `"dark"`: use them for
  the colors of your own objects.

```go
textSoft := map[bool]string{false: "#5A6270", true: "#9AA3B5"}[ui.IsDark()]
ui.Label(ui.Props{Text: "Updated today", TextColor: textSoft})
```

`OZDENUI_THEME=dark go run .` opens any program dark, without a change of its code: a quick way to see both.

## Changing the theme while the program runs

Fifteen components change their theme while they are on the page: `LogView`, `ImageView`, `MarkdownView`, `JsonView`,
`ChatView`, `RangeSlider`, `StepBar`, `Canvas`, `MenuBar`, `ToolBar`, `StatusBar`, `Splitter`, `CodeEditor`,
`Accordion` and `GraphBox`:

```go
logView.Set("theme", "dark")
```

The others take their look when they are made, so a whole window changes its theme by making its page again.
`ui.SetTheme("dark")` does it (`"light"`, `"dark"` or `"auto"`): the page is made again and **`start()` runs again**.
Every object is new, a text typed in a field is gone, and the timers set before are stopped. What your program keeps
in its own variables stays, so keep the data outside `start()`, fill the objects from it in `start()`, and set the
timers there:

```go
package main

import (
	"strconv"

	ui "github.com/bug7a/ozdenui-go"
)

// The colors of this program in each theme: the key is ui.IsDark().
var (
	pageColor = map[bool]string{false: "#F4F5F7", true: "#141413"}
	cardColor = map[bool]string{false: "#FFFFFF", true: "#232322"}
	textSoft  = map[bool]string{false: "#5A6270", true: "#9AA3B5"}
)

var clicks = 0 // outside start(): it stays when the page is made again

func main() {
	ui.Run(start, ui.Window{Title: "Theme", Width: 420, Height: 320})
}

func start(page *ui.Object) {
	dark := ui.IsDark()
	other := map[bool]string{false: "dark", true: "light"}[dark]
	page.Set("color", pageColor[dark])
	var count *ui.Object
	ui.VGroup(ui.Props{Gap: 16}, func() {
		ui.VGroup(ui.Props{Width: 300, Height: "auto", Padding: 20, Gap: 10, Color: cardColor[dark], Round: 12}, func() {
			ui.Label(ui.Props{Text: "Clicks", FontSize: 13, TextColor: textSoft[dark]})
			count = ui.Label(ui.Props{Text: strconv.Itoa(clicks), FontSize: 40, Bold: true})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Click", Tone: "cyan", OnClick: func() {
				clicks++
				count.SetText(strconv.Itoa(clicks))
			}})
		})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Make it " + other, OnClick: func() { ui.SetTheme(other) }})
	})
}
```

`examples/gallery` has a Dark / Light button in its ToolBar for it. `examples/editor` and `examples/markdown-editor`
change only their bars and views with `Call("setTheme", ...)`, so the text in the editor stays.

## The minimal look, in your color

```go
ui.Run(start, ui.Window{Style: "minimal"})                    // small and made of lines; Theme: "dark" works too
ui.Run(start, ui.Window{Style: "minimal", Accent: "#2A78D6"}) // the same, in your color
```

**`Style: "minimal"`** gives about 40 components their minimal look: 28 px controls, a 13 px text, no fill and no shadow
of their own, so they take the ground of the window. A tool, an editor, a panel of settings look calm in it.

**`Accent`** gives the minimal look one color (it goes only with `Style: "minimal"`). It goes where the look draws its
ink: an on toggle, a checked box, the fill of a slider and of a progress bar, the selected tab, row and step, the edge
of a field while typing, and the main button. The lines, the texts and the focus ring stay as they are.

```go
package main

import ui "github.com/bug7a/ozdenui-go"

const accent = "#2A78D6" // a color that reads on white: it is the text of a selected tab too

func main() {
	ui.Run(start, ui.Window{Title: "Minimal", Width: 420, Height: 380, Style: "minimal", Accent: accent})
}

func start(page *ui.Object) {
	page.Set("color", map[bool]string{false: "#FFFFFF", true: "#1A1A19"}[ui.IsDark()])
	ui.VGroup(ui.Props{Gap: 14, Padding: 30, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "Account", FontSize: 18, Bold: true})
		ui.Tabs(ui.TabsProps{Tabs: []string{"General", "Privacy", "Advanced"}, Panels: false, Width: 300})
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Remember me", Checked: true})
		ui.HGroup(ui.Props{Hug: true, Gap: 10}, func() {
			ui.Toggle(ui.ToggleProps{Value: true})
			ui.Label(ui.Props{Text: "Sounds", FontSize: 13})
		})
		ui.ProgressBar(ui.ProgressBarProps{TitleText: "Uploading", Value: 62, Width: 260})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Cancel"})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "accent"})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Delete", StyleName: "minimalDanger"})
		})
	})
}
```

- **`Tone: "accent"`** is the main action of any window: the accent color here, the ink of the minimal look without an
  accent, `"cyan"` in the usual look.
- **`StyleName: "minimalDanger"`**: a delete in the minimal look (in a dark window its dark twin).
- **In a dark window** the color is made lighter until it reads on the dark ground. `AccentDark: "#7FB6DE"` chooses it.
- A component without a minimal look keeps its usual one. `examples/minimal-theme` shows all of them.

## The look of one component

### A package: StyleName

Most components have a few looks, called packages. `go doc` of the function names them
(`go doc github.com/bug7a/ozdenui-go CheckBox`):

```go
ui.CheckBox(ui.CheckBoxProps{LabelText: "Remember me", StyleName: "modern"})
ui.ItemList(ui.ItemListProps{Items: photos, Layout: "grid", StyleName: "gallery"})
ui.SmartButton(ui.SmartButtonProps{LabelText: "Next", StyleName: "pill", Tone: "blue"})
```

A package you name stays in every window; in a dark window a light package with a dark twin takes the twin
(`"gallery"` → `"galleryDark"`, `"minimalDanger"` → `"minimalDangerDark"`). A component is made with its package: a
later `Set("styleName", ...)` does not change it (the terminal says so), except `"theme"` of the fifteen above.

### Your colors: Style

`Style` changes the parts of a component. Its keys are the names of the kit, in camelCase:

```go
ui.SmartButton(ui.SmartButtonProps{LabelText: "Go", Style: ui.P{
	"box":   ui.P{"color": "#2A78D6"},
	"hover": ui.P{"color": "#2367BA"},
	"label": ui.P{"textColor": "white"},
}})
ui.Toggle(ui.ToggleProps{Value: true, StyleName: "classic", BackgroundStyle: ui.P{"selectedColor": "#6AC966"}})
```

- `Style` goes over the package: what you give wins, the rest stays.
- The names of the parts (`box`, `hover`, `label`...) are the ones of the component in basic.js; its file in
  [web/comp](https://github.com/bug7a/ozdenui-go/tree/main/web/comp) and the
  [catalog](https://bug7a.github.io/js-components/) show them.
- A field that ends in `Style` (`BackgroundStyle`, `ButtonStyle`) is a style too.
- `Unframed: true` takes the frame of a component (no border, no rounded corners on its box): a `LogView` or a
  `CodeEditor` that fills the room between two bars.

### Plain objects

A `Box`, a `Label`, a group take their colors from their values (`Color`, `TextColor`, `BorderColor`): chapter 2.
Keep the colors of a program at the top of its file, as constants (or as a map by `ui.IsDark()`), so a change is
one line:

```go
const (
	pageColor = "#F4F5F7"
	cardColor = "#FFFFFF"
	textSoft  = "#5A6270"
)
```

## Next

[Chapter 8: Data and pictures](08-data-and-pictures.md): charts, big tables, logs, Markdown, JSON, code and pictures.
