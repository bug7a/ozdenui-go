# OzdenUI for Go: rules for an AI assistant

This project uses **OzdenUI** (`go get github.com/bug7a/ozdenui-go`): desktop windows written in plain Go. It is a
new library, so do not guess from Fyne, Gio, Wails or a web framework: there is no widget tree to build, no layout
manager object, no HTML, no CSS file and no JavaScript to write. Everything below is the whole way of working.

## Where the truth is

- **The params of a component:** `go doc github.com/bug7a/ozdenui-go SelectBoxProps` (every field with its default and
  what it takes), and `go doc github.com/bug7a/ozdenui-go SelectBox` (its looks, the values it gives back). Read it
  before you use a component for the first time; a field that does not exist does not compile.
- **Every function:** `go doc -short github.com/bug7a/ozdenui-go`; the methods of an object: `go doc github.com/bug7a/ozdenui-go Object`.
- **Example programs:** https://github.com/bug7a/ozdenui-go/tree/main/examples (`gallery` shows the main components).

## A program

```go
package main

import (
	"fmt"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const pageColor = "#F4F5F7"

var (
	notes []string   // the data of the program is plain Go
	name  *ui.Object // the objects the program changes later
	list  *ui.Object
	empty *ui.Object
	count *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Notes", Width: 480, Height: 420}) // once, at the end of main()
}

func start(page *ui.Object) { // runs when the window is ready: make the objects here
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 0}, func() { // a group is as large as its box: here the whole window
		ui.HGroup(ui.Props{Height: 60, Padding: []int{16, 0}, Gap: 10, Align: "left center", Color: "white"}, func() {
			name = ui.TextField(ui.TextFieldProps{Placeholder: "A new note", Width: 260, OnEnter: addNote})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Add", Tone: "cyan", OnClick: addNote})
		})
		list = ui.ScrollArea(ui.ScrollAreaProps{Padding: 16, Gap: 10, Align: "left top"}, func() { // takes the room that is left
			empty = ui.Label(ui.Props{Text: "No notes yet", FontSize: 15})
		})
		count = ui.Label(ui.Props{Text: "Notes: 0", FontSize: 13, Height: 28})
	})
}

func addNote() {
	typed := strings.TrimSpace(name.Text()) // the state is in the program: reading asks nothing of the window
	if typed == "" {
		ui.ToastWarning("Write a note first.")
		return
	}
	notes = append(notes, typed)
	list.Add(func() { // objects made later go into a box the same way
		ui.Label(ui.Props{Text: typed, FontSize: 15})
	})
	empty.Hide()
	name.SetValue("")
	count.SetText(fmt.Sprintf("Notes: %d", len(notes)))
}
```

- `ui.Window{Title, Width, Height, Theme: "light" | "dark" | "auto", Style: "" | "minimal", Accent: "#2A78D6"}`;
  `Width` and `Height` are the size of the page. `ui.Close()`, `ui.SetTitle()`, `ui.Theme()`, `ui.IsDark()`.
  `Language: "tr"` when the texts of the window are not English (the capital letters of Turkish: "İ").
- Objects are made in `start()` or later (an event, a timer), never before `ui.Run()`.
- Every object is a `*ui.Object`: keep the ones the program changes later in variables (or in a struct of your own).

## Objects and values

- Five plain objects: `ui.Box`, `ui.Label`, `ui.Button`, `ui.Input`, `ui.Icon`, all with `ui.Props{...}`. The rest are
  components (below), each with its own type: `ui.CheckBox(ui.CheckBoxProps{...})`.
- A value is a field when the object is made and a name later: `label.Set("text", "Hi", "fontSize", 20)`. The names of
  `Set()` and `Get()` are the names of the fields with a small first letter (`fontSize`, `textColor`, `titleText`).
- Short ways: `obj.SetText("Hi")`, `obj.Text()`, `obj.SetValue(5)`, `obj.Value()`, `obj.Checked()`, `obj.Show()`,
  `obj.Hide()`; a value by its type: `obj.Str("titleText")`, `obj.Num("value")`, `obj.Int("value")`, `obj.Bool("enabled")`.
- Every object: `Color` (the background), `Border`, `BorderColor`, `Round`, `Opacity`, `Visible`, `Width`, `Height`,
  `Css: ui.P{...}`. A size is a number (px) or a text: `"100%"`, `"auto"`. Flags take `true` / `false`.
- **A field that is not written keeps the default of the kit**; a zero that is written is a value (`Gap: 0`,
  `Checked: false`). The one exception is a text: `""` is "not given", so an empty text on purpose is `ui.Empty`
  (`LabelText: ui.Empty`: a button with only an icon).
- `ui.P` is `map[string]any`: a style, an item of a list (`ui.P{"key": "tr", "text": "Türkiye"}`), CSS.
- `Text` is shown exactly as written (a `<b>` stays `<b>`); HTML goes in as `Html`.
- `obj.Remove()` removes it for good; to hide it for a while: `obj.Hide()`.

## Layout

- **Only a `Box`, a group and the areas of some components hold objects**: the objects made inside the function given
  to them go into them. Later: `box.Add(func() { ... })`.
- `ui.VGroup` (one under the other), `ui.HGroup` (side by side), `ui.GridGroup(ui.GridGroupProps{Columns: 4})`:
  `Gap`, `Padding` (`20` or `[]int{x, y}`), `Align` (`"center"` is the default; `"left top"`, `"right center"`...), `Wrap: true`.
- **A group is `Width: "100%"`, `Height: "100%"` by default.** A group inside another one almost always needs one of:
  `Hug: true` (as large as its objects), `Height: "auto"` / a number, or `Fill: true` (the room the others leave: the
  height in a `VGroup`, the width in an `HGroup`). Two groups of 100% in one `VGroup` squeeze each other.
- A typical window: `ui.VGroup(ui.Props{Gap: 0}, ...)` holding a bar of a fixed height, a middle part with
  `Fill: true` (or a `ui.ScrollArea`, which takes the room by itself and scrolls), a line at the bottom.
- Inside the function of a group its own variable is not set yet (`grid := ui.GridGroup(props, func() { grid... })`
  is nil there): the group is `ui.Current()`. A wide cell of a grid: `ui.Current().Place(obj, ui.P{"columnSpan": 2})`.
- In a group an object needs only its size. In a `Box` or on the page it has a place: `Left`, `Top`, `Width`,
  `Height` (for a game or a drawing).

## Events

- `OnClick: fn` when the object is made, `obj.Set("onClick", fn)` later, `obj.On("keydown", fn)` for one more handler
  (it returns a function that removes it). The page listens to the window: `page.On("keydown", fn)`, `page.On("resize", fn)`.
- A handler takes what it needs: `func()`, `func(obj *ui.Object)` or `func(obj *ui.Object, event ui.Event)`
  (`event.Key`, `event.X`, `event.Y`, `event.Shift`...).
- A component tells a change with a callback param: `OnChange: func(box *ui.Object) { ... box.Checked() ... }`. The
  values after the object are the ones the kit names in the comment of the field; each takes the type the function
  asks for (`func(table *ui.Object, row Person)`, `func(self *ui.Object, value string)`). **A change made by the code
  does not call it**, only the user's does.
- In a loop every closure has its own variable (Go 1.22+): `OnClick: func() { pick(number) }`.
- Everything runs one at a time on one goroutine, so no locks; but a slow function freezes the window. Never
  `time.Sleep()` in an event: `ui.After(ms, fn)` once, `ui.Every(ms, fn)` again and again (both return a function that
  stops it). A slow job (a download, a big file, an AI call) runs in a goroutine of its own and changes the objects
  only inside `ui.Do(func() { ... })`:

```go
go func() {
	text := download()                        // the window does not wait
	ui.Do(func() { label.SetText(text) })     // back on the goroutine of the window
}()
```

## Components

Each is a function of `ozdenui` with its `...Props` type; its state is read from the program (`field.Text()`,
`selectBox.Value()`, `box.Checked()`, `obj.Get("selectedKeys")`: the function's doc names them).

| For | Functions |
|---|---|
| Buttons | `SmartButton` (`LabelText`, `Tone: "cyan"` main / `"orange"` delete / `"accent"`, `OnClick`) |
| Fields | `TextField`, `EmailField`, `PasswordField`, `NumberField`, `CurrencyField`, `PhoneField`, `UrlField`, `TextAreaField`, `SearchInput`, `TagInput` (`TitleText`, `Placeholder`, `Value`, `OnEnter`, `OnChange`) |
| Choices | `CheckBox`, `Toggle`, `RadioButton`, `SelectBox` (`Options`, `Value`), `Slider`, `SliderField`, `RangeSlider` (`LowValue`, `HighValue`; read `Get("low")`, `Get("high")`), `Stepper`, `DatePicker` (`Date`), `TimePicker` (`Time`), `ColorPicker` (`SelectedColor`); the others give `obj.Value()` |
| Lists and tables | `ItemList` (`Items`), `SortableList`, `TreeView`, `SmartTable` (`TitleDataList`, `ItemDataList`: very many rows), `DataTable` |
| Areas that hold objects | `ScrollArea`, `GridGroup`, `Tabs` (`tabs.Panel(key, func() {...})`), `Accordion` (`acc.Section(key, func() {...})`), `Splitter` (`split.First(...)`, `split.Second(...)`), `SidePanel` (`panel.Body(...)`) |
| Bars | `MenuBar`, `ToolBar`, `StatusBar`, `ContextMenu` |
| Messages | `ui.Alert`, `ui.Confirm`, `ui.Prompt`, `ui.Ask`, `ui.ToastSuccess` / `ToastError` / `ToastWarning` / `ToastInfo`, `Tooltip`, `LoadingScreen`, `ProgressBar`, `StepBar` |
| Views | `GraphBox` (charts), `LogView`, `ImageView`, `MarkdownView`, `JsonView`, `JsonEditor`, `CodeEditor`, `Canvas`, `ChatView`, `WebView`, `FilePicker`, `Gauge` |

- A dialog returns its answer, like a question in a terminal: `if ui.Confirm("Delete it?", ui.DialogOptions{Danger: true}) {`,
  `name, ok := ui.Prompt("Name:")` (`ok` false: cancelled), `ui.Alert(...)`, `ui.Ask("Save?", "Save", "Don't save", "Cancel")`
  (the text of the pressed button). The window keeps working while the code waits.
- The file dialogs of the system: `ui.OpenFile()`, `ui.OpenFiles()`, `ui.SaveFile()`, `ui.ChooseFolder()` return a real
  path or `""`; `ui.FileOptions{Title, Directory, FileName, Types: []string{"png", "jpg"}}`.
  A click on a `FilePicker` opens that dialog too: `picker.Paths()` are the chosen files (a dropped file has no path).
- A picture (`ui.Icon(ui.Props{Picture: "photo.png"})`, a field whose name ends with `File`: `IconFile`, `ImageFile`)
  takes a path of this computer; in another place (an item of a list): `ui.Picture("photo.png")`, also of bytes or an
  `image.Image`.
- A row of a table and an item of a list may be a struct of the program: its `json` tags are its keys. A callback gives
  it back in the type the function asks for (`OnSelect: func(table *ui.Object, row Person)`).
- A param of the kit changes later with `obj.Set("unitText", "kg")`; a method of the kit: `obj.Call("addItem", item)`
  (not waited for), `obj.Ask("getSelectedKeys")` (its result; a question to the page, not in a loop). `go doc` of the
  function of a component lists its methods; what each takes is in its file of the kit:
  https://github.com/bug7a/ozdenui-go/tree/main/web/comp
- What comes from the page (the result of `Ask()`, a value of a callback taken as `any`): a number is a `float64`, a
  flag of the kit is `1` / `0` (not a `bool`), a list is `[]any`, an object `map[string]any`. Read them with
  `ui.IsOn(v)`, `ui.Number(v)`, `ui.Text(v)`, or let the callback ask for the type (`func(self *ui.Object, on bool)`).
- A big or fast value is asked, not kept: the text of a `CodeEditor` (`editor.Ask("getValue")`), the value of a
  `JsonEditor`, the selected row of a `SmartTable` (`table.Measure("selectedItemData")`).
- The columns of a `SmartTable`: `TitleDataList: []ui.P{{"name": "NAME", "dataTitle": "name", "dataType": "string",
  "width": 160, "sortable": 1}}`; its rows: `ItemDataList` (structs or `ui.P`), later `table.Set("itemDataList", rows)`.
- The keys of a `ui.P` go to the page in the order of the alphabet; a struct keeps the order of its fields (a
  `JsonView`, a `JsonEditor` show it).
- Small helpers: `field.Check()` (shows the warning of a field, true when its text is right), `ui.RadioValue(group)`
  (the `Value` of the checked `RadioButton`), `ui.ShowHint(obj, "Copied", 1500)` (a bubble for a while),
  `logView.Writer("info")` (an `io.Writer`, from any goroutine), `ui.AskWithTitle(title, message, buttons...)`,
  `obj.Set("theme", "dark")` (the look of one bar, editor or view while the program runs).
- A picture given to a method is an address: `view.Call("setImage", ui.Picture(path))`.
- A number a component writes (`FormatValue` of a `GraphBox`, `ValueFormat` of a `Stepper`) takes a text like
  `"{:,.2f} TL"` (1,234.50 TL), `"{:.0f} kg"`, `"{:.1%}"`: not a function.
- A component with no fields for something: `Extra: ui.P{"name": value}` by its name in the kit, and any component
  by its name: `ui.Component("Gauge", ui.P{"value": 72})`.

## Looks

- `ui.Window{Theme: "dark"}`: every component takes its dark look by itself. For your own colors:
  `textColor := map[bool]string{false: "#5A6270", true: "#9AA3B5"}[ui.IsDark()]`.
- `ui.Window{Style: "minimal", Accent: "#2A78D6"}`: a small look made of lines, in one color; the main button is
  `Tone: "accent"`.
- `ui.SetTheme("dark")` changes the theme while the program runs: the page is made again and `start()` runs again,
  so every object is new and the timers set before are gone. Keep the data in variables outside `start()`, fill the
  objects from them in `start()`, and set the timers there.
- One component: `StyleName: "..."` (its packages are named in the doc of its function), or
  `Style: ui.P{"box": ui.P{"color": "#FFF"}}`; `Unframed: true`: no border and no rounded corners (a table or an
  editor that fills the room between two bars).
- Keep the colors of a program as constants at the top of its file.

## Mistakes to avoid

1. A nested group without `Hug: true`, `Fill: true` or a height (see Layout).
2. Changing an object from a goroutine of your own: it panics with the fix (`ui.Do(func() { ... })`).
3. Waiting in an event (`time.Sleep`, a download): the window freezes.
4. Reading a value that was never given (`label.Get("width")` of an `"auto"` label) is `nil`:
   `obj.Measure("width")` asks the window. Keep `Measure()` and `Ask()` out of loops.
5. Putting an object into a `Label` or a button. Make it next to it, in the same box.
6. A name of `Set()`, `Extra` or `ui.P` that does not exist. The terminal says it with the line of the program
   (`ozdenui: main.go:12: Label: "fontSzie" is not a value of Label (did you mean "fontSize"?)`): run the program
   and **read the terminal**.
7. Colors written for a light window only, in a program that may open dark.
8. A callback that must answer at once (a validator that returns `false`) can not be a function of the program.
9. Saying a window is right without a look at it (below).

## Running and checking

- **Look at the window you made:** `OZDENUI_LOOK=app.png go run .` runs the program as it is, saves a picture of its
  window as `app.png` and ends. Read the picture: a program that runs without an error can still have a part covered,
  cut or squeezed. `OZDENUI_THEME=dark OZDENUI_LOOK=dark.png OZDENUI_LOOK_WAIT=3 go run .`: the dark window, more
  seconds before the picture. (It needs Chrome, Edge or Chromium on the computer; no window opens.)
- `go run .` opens a real window and returns when it is closed. To check a program without a person, close it from
  the code: `ui.After(3000, ui.Close)` as the first line of `start()`, then read the terminal.
- A check that plays the program from code must not reach a dialog: it waits for a person.
- `ui.Window{Debug: true}` gives the developer tools of the web view (right click > Inspect).
- macOS: the window is the web view of the system (it needs cgo: the command line tools of Xcode). Windows: the web
  view of the system (WebView2, no C compiler). Linux: the app window of Chrome, Chromium or Edge. The first build takes a while and may print warnings of the web view's C code.
