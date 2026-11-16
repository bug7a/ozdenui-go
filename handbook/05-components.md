# 5. Components

A component is a ready part of a window, made of the objects of chapter 2: a button with a waiting clock, a field with
a title and a warning, a list with a search box, a table of 100000 rows. 83 of them have a function of their own,
each with its own `...Props` type; `ui.Component("Name", ...)` makes any of them by its name.

## How every component works

```go
remember := ui.CheckBox(ui.CheckBoxProps{LabelText: "Remember me", Checked: true, OnChange: show})
```

- **Params are fields** of its Props type. A field is the param of the kit with a capital first letter: `TitleText`,
  `SelectedKeys`, `PanelHeight`. `go doc github.com/bug7a/ozdenui-go CheckBoxProps` lists every field with its default
  and what it takes, and your editor offers them. A wrong name does not compile.
- **A field that is not written keeps the default of the kit.** A zero that is written is a value
  (`Checked: false`, `Value: 0`). The one exception is a text: `""` is "not given", so an empty text on purpose is
  `ui.Empty` (`LabelText: ui.Empty`: a button with only an icon).
- **A callback is a field**: `OnChange: show`. Like an event, it takes what it needs: `show()`, `show(box *ui.Object)`,
  or more. The values after the object are the ones the comment of the field names, each in the type your function
  asks for: `OnSelect: func(table *ui.Object, row Person)` gets the row as a `Person` (by its `json` tags),
  `OnStepClick: func(bar *ui.Object, index int)` gets a number.
- **Its state is in Go**: `remember.Checked()`, `selectBox.Value()`, `field.Text()`, `list.Get("selectedKeys")`.
  Reading it asks nothing of the page; the doc of each function names what it keeps.
- **Changing it**: `remember.SetValue(...)`, `bar.SetValue(40)`, `field.SetText("")`, or a param by its name:
  `field.Set("unitText", "kg")`. A param the kit has a method for (`setUnitText`) is changed with that method, so the
  component draws itself again.
- **A method of the kit**: `list.Call("addItem", item)` (not waited for), `list.Ask("getSelectedKeys")` (waits for its
  result; a question to the page, so not in a loop). The doc of the function lists its methods; what each takes is in
  [its file of the kit](https://github.com/bug7a/ozdenui-go/tree/main/web/comp).
- **A change made by your code does not call the callback**; only the user's does.

## A form on two tabs

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const formWidth = 400

var name, language, news, volume, sounds *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "Settings", Width: 540, Height: 580})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 16, Padding: 24, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "Settings", FontSize: 24, Bold: true})
		ui.VGroup(ui.Props{Width: 460, Height: "auto", Padding: 20, Gap: 16, Align: "left top", Color: "white", Round: 12}, func() {
			tabs := ui.Tabs(ui.TabsProps{Tabs: []string{"Profile", "Sound"}, PanelHeight: 240})
			tabs.Panel("Profile", func() {
				ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 14, Padding: []int{0, 16}, Align: "left top"}, func() {
					name = ui.TextField(ui.TextFieldProps{TitleText: "Name", Required: true, Width: formWidth})
					language = ui.SelectBox(ui.SelectBoxProps{Options: []string{"English", "Türkçe", "Deutsch"}, Value: "English", Width: formWidth})
					news = ui.CheckBox(ui.CheckBoxProps{LabelText: "Send me the news"})
				})
			})
			tabs.Panel("Sound", func() {
				ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 16, Padding: []int{0, 16}, Align: "left top"}, func() {
					volume = ui.Slider(ui.SliderProps{Value: 7, MinValue: 0, MaxValue: 10, Width: formWidth})
					ui.HGroup(ui.Props{Hug: true, Gap: 18}, func() {
						for _, quality := range []string{"Low", "Normal", "High"} {
							ui.RadioButton(ui.RadioButtonProps{Group: "quality", Value: quality, LabelText: quality, Checked: quality == "Normal"})
						}
					})
					sounds = ui.Toggle(ui.ToggleProps{Value: true})
				})
			})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan", OnClick: save})
		})
	})
}

func save() {
	if !name.Check() { // shows the warning on the field
		return
	}
	fmt.Println(name.Text(), language.Value(), news.Checked(), volume.Value(), ui.RadioValue("quality"), sounds.Checked())
	ui.ToastSuccess("Saved.")
}
```

`Value()` is `any`: the value your code gave (the `int` 7 of the slider) until the user changes it, then the value
the page tells (a number of the page is a `float64`). `ui.Number(v)`, `ui.Text(v)` and `ui.IsOn(v)` read either kind:
`ui.Number(volume.Value())` is always a `float64`. `examples/settings` is a longer one with three tabs,
`examples/form` a form with dates, a time and a color.

## The components, by what they do

### Buttons

| | |
|---|---|
| `ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan", OnClick: save})` | `Tone`: `"cyan"` the main action, `"orange"` a delete, `"yellow"`, `"green"`, `"blue"`, `"white"`; `"accent"`: the main action of a minimal window (chapter 7) |
| `ui.SmartButtonProps{LabelText: ui.Empty, IconFile: "plus.svg", AriaLabel: "Add"}` | an icon button: a square |
| `button.Set("waiting", true)` | a clock on it; it takes no clicks until `button.Set("waiting", false)` |
| `ui.SmartButtonProps{LabelText: "Delete", Ask: "Click again to delete."}` | the first click asks, the second one runs `OnClick` |
| `button.Set("enabled", false)`, `button.SetText("Saved")` | |

### Fields

| | |
|---|---|
| `ui.TextField(ui.TextFieldProps{TitleText: "Name", Placeholder: "Ada", OnChange: fn, OnEnter: fn})` | `field.Text()`; `OnChange` at every change, `OnEnter` on Enter |
| `EmailField`, `PasswordField`, `NumberField`, `TextAreaField` | the same, for an e-mail, a password, a number, many lines |
| `Required: true`, then `field.Check()` | shows the warning on the field; `true` when there is none |
| `MaxLength: 40`, `ui.PasswordFieldProps{MinLength: 8}` | the most and the fewest letters |
| `ui.UrlField(ui.UrlFieldProps{TitleText: "Site"})` | an address, with its protocol |
| `ui.PhoneField(...)`, `ui.CurrencyField(ui.CurrencyFieldProps{UnitText: "$"})` | the number in its mask; the amount in groups |
| `ui.SearchInput(ui.SearchInputProps{Placeholder: "Find", OnChange: fn})` | `fn(search *ui.Object)` at every change; `search.Text()` |
| `ui.TagInput(ui.TagInputProps{Tags: []string{"audio"}, Suggestions: []string{...}})` | `field.Get("tags")` |

`Check()` shows the warning of one field. To show every warning of a form, check them all first, then decide:
`nameOK, mailOK := name.Check(), email.Check()` (with `&&` the second one would not be checked when the first fails).

### Choices

| | |
|---|---|
| `ui.CheckBox(ui.CheckBoxProps{LabelText: "Remember me", Checked: true})` | `box.Checked()` |
| `ui.Toggle(ui.ToggleProps{Value: true})` | `toggle.Checked()` |
| `ui.RadioButton(ui.RadioButtonProps{Group: "size", Value: "M", LabelText: "Medium"})` | `ui.RadioValue("size")` |
| `ui.SelectBox(ui.SelectBoxProps{Options: []string{"Apple", "Pear"}, Value: "Pear"})` | `selectBox.Value()` |
| `ui.Slider(ui.SliderProps{Value: 30, MinValue: 0, MaxValue: 100, OnChangeEnd: fn})` | `slider.Value()`; `OnChangeEnd`: when the mouse lets it go |
| `ui.SliderField(...)`, `ui.Stepper(ui.StepperProps{Value: 2, MinValue: 0, MaxValue: 10})` | `.Value()` |
| `ui.RangeSlider(ui.RangeSliderProps{LowValue: 200, HighValue: 800, MaxValue: 1000, Step: 50})` | `slider.Get("low")`, `slider.Get("high")` |
| `ui.DatePicker(ui.DatePickerProps{Date: "2026-11-20"})`, `ui.TimePicker(ui.TimePickerProps{Time: "19:30"})`, `ui.ColorPicker(ui.ColorPickerProps{SelectedColor: "#3871E0"})` | `.Value()`: `"2026-11-20"`, `"19:30"`, `"#3871E0"` |

### Lists, tables and trees

```go
fruits := ui.ItemList(ui.ItemListProps{
	Items:     []any{"Apple", ui.P{"key": "p", "text": "Pear", "description": "Green"}},
	Selection: "multiple", Height: 240,
	OnChange: func(list *ui.Object) { fmt.Println(list.Get("selectedKeys")) },
})
fruits.Call("addItem", ui.P{"key": "c", "text": "Cherry"})
fruits.Call("removeItem", "p")
```

- **`ItemList`**: `list.Get("selectedKeys")` (a `[]any` of the keys), `list.Get("selected")`. A text item gets its
  index as its key (`"0"`); `"enabled": false` can not be chosen. `Layout: "grid", StyleName: "gallery"` shows
  pictures. `examples/todo` is a to-do list made of one.
- **An item** is a text, a `ui.P` or a struct of yours: its `json` tags are its keys. The keys are the names of the
  kit: `text`, `description`, `iconFile`, `enabled`...
- **`SmartTable`**: very many rows (only the rows that are seen are made), a filter box and sorting. Its columns:
  `TitleDataList: []ui.P{{"name": "NAME", "dataTitle": "name", "width": 160, "sortable": 1}}`; its rows:
  `ItemDataList` (structs or `ui.P`); `OnSelect: func(table *ui.Object, row Person)`. Chapter 8 has more.
- **`DataTable`** (`Columns`, `Rows`): a smaller table; `table.Get("selectedRow")`. **`TreeView`**
  (`Items: []ui.P{{"key": "a", "text": "A", "children": []string{"b"}}}`): `tree.Value()` is the selected key.
  **`SortableList`**: the user changes the order; `list.Get("order")` gives the keys in the new order.

### Areas that hold your objects

Some components hold objects of yours. Each gives a function of the area, as a box does:

```go
tabs := ui.Tabs(ui.TabsProps{Tabs: []string{"Orders", "Customers"}})
tabs.Panel("Orders", func() { ui.Label(ui.Props{Text: "No orders yet"}) })

folds := ui.Accordion(ui.AccordionProps{Sections: []string{"General", "Advanced"}})
folds.Section("General", func() { ui.CheckBox(ui.CheckBoxProps{LabelText: "Start with the system"}) })

split := ui.Splitter(ui.SplitterProps{FirstSize: 240}) // two areas, a line to drag between them
split.First(func() { ui.Label(ui.Props{Text: "Left"}) })
split.Second(func() { ui.Label(ui.Props{Text: "Right"}) })

panel := ui.SidePanel(ui.SidePanelProps{TitleText: "Filters"}) // slides in from the side
panel.Body(func() { ui.CheckBox(ui.CheckBoxProps{LabelText: "Only the new ones"}) })
ui.SmartButton(ui.SmartButtonProps{LabelText: "Filters", OnClick: func() { panel.Call("open") }})
```

A tab or a section given as a text has that text as its key. A `Splitter` is as large as its box: in a group, give
it `Fill: true` or a size.

### The bars of a window

| | |
|---|---|
| `ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{{"text": "File", "items": []ui.P{{"text": "Save", "shortcut": "Ctrl+S", "onClick": save}}}}})` | in the page, in the look of the window; Ctrl is Cmd on a Mac |
| `ui.ToolBar(ui.ToolBarProps{Items: []any{ui.P{"key": "save", "icon": "save", "hint": "Save"}, "-"}, OnClick: fn})` | `"-"` a line; `"toggle": true` stays pressed; `fn(bar *ui.Object, item ui.P)` |
| `ui.StatusBar(ui.StatusBarProps{MessageText: "Ready", Items: []string{"UTF-8"}})` | `bar.Call("setMessageText", "Saved", "success", 3000)`, `bar.Call("setBusy", true)` |
| `ui.ContextMenu(ui.ContextMenuProps{Target: table, Items: []any{"Open", "-", ui.P{"text": "Delete", "onClick": fn}}})` | a right click on the target opens it |

A window with a menu, a tool bar, a list that takes the room and a status bar is a `VGroup` with `Gap: 0` of them,
with `Fill: true` on the middle part (chapter 3):

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var status *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "Files", Width: 620, Height: 420})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Gap: 0}, func() {
		ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{
			{"text": "File", "items": []any{
				ui.P{"text": "New", "shortcut": "Ctrl+N", "onClick": func() { tell("New") }},
				"-",
				ui.P{"text": "Close", "onClick": ui.Close},
			}},
		}})
		ui.ToolBar(ui.ToolBarProps{OnClick: pressTool, Items: []any{
			ui.P{"key": "new", "icon": "new", "hint": "New"},
			ui.P{"key": "save", "icon": "save", "hint": "Save"},
		}})
		ui.ItemList(ui.ItemListProps{Fill: true, Unframed: true, OnChange: choose, Items: []ui.P{
			{"key": "report", "text": "report.pdf", "description": "120 KB"},
			{"key": "photo", "text": "photo.jpg", "description": "2.4 MB"},
			{"key": "notes", "text": "notes.txt", "description": "4 KB"},
		}})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Ready", Items: []string{"3 files"}})
	})
}

func tell(text string) {
	status.Call("setMessageText", text, "info", 2000)
}

func pressTool(bar *ui.Object, item ui.P) { tell(fmt.Sprintf("%v was pressed", item["hint"])) }

func choose(list *ui.Object) { tell(fmt.Sprintf("Chosen: %v", list.Get("selectedKeys"))) }
```

`Unframed: true` takes the border and the round corners of the list away: it fills the room between two bars.
`examples/editor` and `examples/gallery` are bigger windows of this kind.

### Progress and steps

| | |
|---|---|
| `ui.ProgressBar(ui.ProgressBarProps{TitleText: "Copying", Value: 0})` | `bar.SetValue(40)` |
| `ui.LoadingScreen(ui.LoadingScreenProps{TitleText: "My App", MessageText: "Reading..."})` | covers the window; `loading.Call("setProgress", 0.4)`, `loading.Call("close")` |
| `ui.StepBar(ui.StepBarProps{Steps: []string{"Cart", "Address", "Payment"}, CurrentStep: 1})` | `bar.Call("next")`, `bar.Call("previous")`, `bar.Get("step")` |
| `ui.Tooltip(ui.TooltipProps{Target: button, HintText: "Saves the file"})` | a hint under the mouse; chapter 6 has its other forms |

### Views: chapter 8

`GraphBox` (charts), `LogView`, `ImageView`, `MarkdownView`, `JsonView`, `JsonEditor`, `CodeEditor`, `Canvas`,
`ChatView` and `WebView` show data, pictures, code, a conversation and web pages (`JsonEditor` and `CodeEditor` also
write them): chapter 8.

### The texts of a component

The texts a component writes itself (a warning, the buttons of a dialog, the message of an empty table) are in
English. Each one changes with a field:

```go
ui.TextField(ui.TextFieldProps{TitleText: "Name", Required: true, RequiredText: "Needed"}) // in the place of "Required"
ui.Confirm("The file will be deleted.", ui.DialogOptions{Title: "Delete it?", ConfirmText: "Delete", CancelText: "Keep it", Danger: true})
ui.SmartTable(ui.SmartTableProps{NoDataFoundAlert: "No records", AllTitlesText: "ALL"})
ui.DatePicker(ui.DatePickerProps{Date: "2026-11-20", Language: "tr"}) // the months and the days in Turkish
ui.JsonEditor(ui.JsonEditorProps{Texts: ui.P{"addRow": "Row", "search": "Search the rows", "table": "Table"}})
```

The text fields of a component end with `Text` or say what they are (`Placeholder`, `EmptyText`); `go doc` of its
Props type shows each one with its English default. The texts of a `JsonEditor` are in one `Texts` map, by the names
of the kit (`addRow`).

## Any component by its name: ui.Component()

A function with a Props type is the easy way. `ui.Component()` makes the same component from its name and a `ui.P`
of its params, by their names in the kit; it is the way for a name your program chooses while it runs:

```go
ui.Gauge(ui.GaugeProps{Value: 72, SubText: "Disk"})
ui.Component("Gauge", ui.P{"value": 35, "subText": "Memory"}) // the same component, by its name
```

- The page loads a component the first time your program names it, with what it needs: a program pays only for what
  it uses. Nothing is listed before `ui.Run()`.
- A name of `ui.P` that the component does not know is told in the terminal, with the line of your program.
- A param that its Props type has no field for: `Extra: ui.P{"name": value}`, by its name in the kit.
- Every component, with an example of each: [the catalog](https://bug7a.github.io/js-components/) and
  `go doc github.com/bug7a/ozdenui-go`. A box inside a component, to put objects into or to call: `obj.Part("name")`.

## Next

[Chapter 6: Dialogs, messages and files](06-dialogs-and-files.md).
