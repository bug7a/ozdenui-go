/* Bismillah */

// The minimal look of every component that has it, on one page: ui.Window{Style: "minimal"}. Small (28 px controls,
// 13 px text) and made of lines; the main action is the only filled button. OZDENUI_THEME=dark: the minimalDark looks.
package main

import (
	"fmt"
	"net/url"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	columnWidth = 420
	fieldWidth  = 260
	loadingTime = 1500 // ms: the loading screen of this page stays this long
)

// The colors of this page in each theme. (The components take their look from the window.)
type palette struct{ ground, caption, text, line string }

var (
	lightColors = palette{ground: "#FFFFFF", caption: "rgba(0, 0, 0, 0.45)", text: "rgba(0, 0, 0, 0.85)", line: "rgba(0, 0, 0, 0.2)"}
	darkColors  = palette{ground: "#1A1A19", caption: "rgba(255, 255, 255, 0.45)", text: "rgba(255, 255, 255, 0.88)", line: "rgba(255, 255, 255, 0.18)"}
)

type person struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Orders int    `json:"orders"`
}

var people = []person{{"Ada Lovelace", "London", 12}, {"Alan Turing", "Manchester", 7}, {"Grace Hopper", "New York", 21}}

var (
	iconSave = lineIcon(`<path d="M5 4h11l3 3v13H5z"/><path d="M8 4v5h7V4M8 20v-6h8v6"/>`)
	iconMore = lineIcon(`<circle cx="5" cy="12" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/>`)
)

var (
	colors palette
	status *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "OzdenUI minimal", Width: 1000, Height: 820, Style: "minimal"})
}

// A line icon (24 x 24) as a data: address; SmartButton draws it white in a dark window.
func lineIcon(paths string) string {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#000000" stroke-width="2" ` +
		`stroke-linecap="round" stroke-linejoin="round">` + paths + `</svg>`
	return "data:image/svg+xml;utf8," + url.PathEscape(svg)
}

// The small title over a part of the page.
func caption(text string) {
	ui.Label(ui.Props{Text: text, FontSize: 11, Bold: true, TextColor: colors.caption})
}

func tell(text string) {
	status.Set("messageText", text)
}

func start(page *ui.Object) {
	colors = lightColors
	if ui.IsDark() {
		colors = darkColors
	}
	page.Set("color", colors.ground)
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{
			{"text": "File", "items": []any{ui.P{"text": "Open", "shortcut": "Ctrl+O"}, ui.P{"text": "Save", "shortcut": "Ctrl+S"}, "-",
				ui.P{"text": "Quit", "enabled": false}}},
			{"text": "Edit", "items": []any{"Undo", "Redo", "-", ui.P{"text": "Transform", "items": []string{"Upper case", "Lower case"}}}},
			{"text": "View", "items": []ui.P{{"text": "Word wrap", "checked": true}}},
		}})
		ui.ToolBar(ui.ToolBarProps{
			Items: []any{
				ui.P{"key": "new", "icon": "new", "hint": "New"},
				ui.P{"key": "save", "icon": "save", "hint": "Save", "shortcut": "Ctrl+S"},
				"-",
				ui.P{"key": "bold", "icon": "bold", "hint": "Bold", "toggle": true, "checked": true},
				ui.P{"key": "left", "icon": "alignLeft", "hint": "Left", "group": "align", "checked": true},
				ui.P{"key": "right", "icon": "alignRight", "hint": "Right", "group": "align"},
			},
			OnClick: func(bar *ui.Object, item ui.P) {
				tell(fmt.Sprintf("%v, pressed: %v", item["hint"], bar.Get("checkedKeys")))
			},
		})
		ui.ScrollArea(ui.ScrollAreaProps{}, func() {
			ui.HGroup(ui.Props{Width: "100%", Height: "auto", Align: "center top", Gap: 40, Padding: 30, Wrap: true}, func() {
				ui.VGroup(ui.Props{Width: columnWidth, Height: "auto", Align: "left top", Gap: 14}, func() {
					drawButtons()
					drawChoices()
					drawFields()
					drawLists()
					drawConversation()
				})
				ui.VGroup(ui.Props{Width: columnWidth, Height: "auto", Align: "left top", Gap: 14}, func() {
					drawPickers()
					drawStepsAndAreas()
					drawTreeAndTables()
					drawChart()
					drawOverThePage()
				})
			})
		})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Change something.", Items: []string{"UTF-8", "100%"}})
	})
}

func drawButtons() {
	caption("BUTTONS")
	ui.HGroup(ui.Props{Hug: true, Align: "left center", Gap: 8}, func() {
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", IconFile: iconSave, OnClick: func() { tell("Save was clicked.") }})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Send", Tone: "accent", OnClick: func() { ui.ToastSuccess("Sent.") }})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Delete", StyleName: "minimalDanger", Ask: "Click again to delete.",
			OnClick: func() { tell("Deleted.") }})
		ui.SmartButton(ui.SmartButtonProps{LabelText: ui.Empty, IconFile: iconMore, AriaLabel: "More", OnClick: func() { tell("More was clicked.") }})
	})
	ui.HGroup(ui.Props{Hug: true, Align: "left center", Gap: 8}, func() {
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Disabled", Enabled: false})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Waiting", Waiting: true})
	})
}

func drawChoices() {
	caption("CHOICES")
	ui.HGroup(ui.Props{Hug: true, Align: "left center", Gap: 20}, func() {
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Remember me", Checked: true, OnChange: func(box *ui.Object) { tell(fmt.Sprintf("Remember me: %v", box.Checked())) }})
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Send the news", OnChange: func(box *ui.Object) { tell(fmt.Sprintf("Send the news: %v", box.Checked())) }})
	})
	ui.HGroup(ui.Props{Hug: true, Align: "left center", Gap: 20}, func() {
		for _, plan := range []string{"Monthly", "Yearly"} {
			ui.RadioButton(ui.RadioButtonProps{Group: "plan", Value: plan, LabelText: plan, Checked: plan == "Monthly",
				OnChange: func() { tell("Plan: " + plan) }})
		}
	})
	ui.HGroup(ui.Props{Hug: true, Align: "left center", Gap: 10}, func() {
		ui.Toggle(ui.ToggleProps{Value: true, OnChange: func(toggle *ui.Object) { tell(fmt.Sprintf("Sounds: %v", toggle.Bool("value"))) }})
		ui.Label(ui.Props{Text: "Sounds", FontSize: 13, TextColor: colors.text})
	})
}

func drawFields() {
	caption("FIELDS")
	ui.EmailField(ui.EmailFieldProps{TitleText: "E-mail", Placeholder: "name@mail.com", Width: fieldWidth})
	ui.PasswordField(ui.PasswordFieldProps{TitleText: "Password", Value: "secret", Width: fieldWidth})
	ui.TextAreaField(ui.TextAreaFieldProps{TitleText: "Note", Placeholder: "A few words", Width: fieldWidth, Height: 96,
		OnChange: func(area *ui.Object) { tell(fmt.Sprintf("Note: %d letters", len([]rune(area.Text())))) }})
	ui.SearchInput(ui.SearchInputProps{Placeholder: "Search", Width: fieldWidth, OnChange: func(search *ui.Object) { tell("Search: " + search.Text()) }})
	ui.SelectBox(ui.SelectBoxProps{Options: []string{"Apple", "Banana", "Cherry", "Grape"}, Value: "Apple", Width: fieldWidth,
		OnChange: func(box *ui.Object) { tell(fmt.Sprintf("Fruit: %v", box.Value())) }})
	ui.Stepper(ui.StepperProps{Value: 3, MinValue: 0, MaxValue: 10, OnChange: func(stepper *ui.Object) { tell(fmt.Sprintf("Count: %v", stepper.Value())) }})
}

func drawLists() {
	caption("LIST, SLIDER, PROGRESS, TABS")
	ui.ItemList(ui.ItemListProps{Items: []ui.P{{"key": "p", "text": "Projects"}, {"key": "d", "text": "Designs"}, {"key": "a", "text": "Archive"}},
		SelectedKeys: []string{"d"}, Width: fieldWidth, Height: "auto",
		OnChange: func(list *ui.Object) { tell(fmt.Sprintf("Folder: %v", list.Get("selectedKeys"))) }})
	ui.Slider(ui.SliderProps{Value: 40, MinValue: 0, MaxValue: 100, Width: fieldWidth,
		OnChangeEnd: func(slider *ui.Object) { tell(fmt.Sprintf("Slider: %v", slider.Value())) }})
	ui.ProgressBar(ui.ProgressBarProps{TitleText: "Uploading files", Value: 62, Width: fieldWidth})
	ui.Tabs(ui.TabsProps{Tabs: []string{"General", "Sound", "Advanced"}, Panels: false, Width: 300,
		OnChange: func(tabs *ui.Object) { tell(fmt.Sprintf("Tab: %v", tabs.Value())) }})
}

func drawPickers() {
	caption("PICKERS, FIELDS")
	ui.HGroup(ui.Props{Width: "100%", Height: "auto", Align: "left center", Gap: 8, Wrap: true}, func() {
		ui.DatePicker(ui.DatePickerProps{Date: "2026-11-12", Width: 140, OnChange: func(picker *ui.Object) { tell("Date: " + picker.Str("value")) }})
		ui.TimePicker(ui.TimePickerProps{Time: "14:30", Width: 100, OnChange: func(picker *ui.Object) { tell("Time: " + picker.Str("value")) }})
		ui.ColorPicker(ui.ColorPickerProps{SelectedColor: "#3871E0", Width: 130, OnChange: func(picker *ui.Object) { tell("Color: " + picker.Str("value")) }})
	})
	ui.SliderField(ui.SliderFieldProps{LabelText: "Size", Value: 40, MinValue: 0, MaxValue: 100, UnitText: " px", Width: fieldWidth,
		OnChangeEnd: func(slider *ui.Object) { tell(fmt.Sprintf("Size: %v", slider.Value())) }})
	ui.TagInput(ui.TagInputProps{Width: "100%", Tags: []string{"audio", "guide"}, Suggestions: []string{"audio", "fitness", "guide", "travel"},
		OnChange: func(field *ui.Object) { tell(fmt.Sprintf("Tags: %v", field.Get("tags"))) }})
	ui.FilePicker(ui.FilePickerProps{Multiple: true, Width: "100%", OnChange: func(picker *ui.Object) {
		files, _ := picker.Get("files").([]any)
		tell(fmt.Sprintf("Files: %d", len(files)))
	}})
}

func drawStepsAndAreas() {
	caption("STEPS, SECTIONS, AREAS")
	ui.StepBar(ui.StepBarProps{Steps: []string{"Cart", "Address", "Payment", "Done"}, CurrentStep: 2, Width: "100%",
		OnStepClick: func(bar *ui.Object, index int) { tell(fmt.Sprintf("Back to step %d", index+1)) }})
	settings := ui.Accordion(ui.AccordionProps{Width: "100%", Height: "auto", Sections: []ui.P{
		{"key": "general", "text": "General", "description": "3 settings", "opened": 1},
		{"key": "advanced", "text": "Advanced"},
		{"key": "locked", "text": "Locked", "enabled": 0},
	}})
	settings.Section("general", func() {
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Open with the system"})
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Look for updates", Checked: true})
	})
	ui.Box(ui.Props{Width: "100%", Height: 80, Color: "transparent", Border: 1, BorderColor: colors.line, Round: 4}, func() {
		ui.Splitter(ui.SplitterProps{FirstSize: 160, OnChangeEnd: func() { tell("The line was moved.") }})
	})
}

func drawConversation() {
	caption("CONVERSATION")
	ui.ChatView(ui.ChatViewProps{Width: "100%", Height: 300,
		Messages: []any{
			ui.P{"role": "user", "text": "How do I count the lines of a file?"},
			"Read the file and count its lines:\n\n```python\nlines = open(\"notes.txt\").read().splitlines()\nprint(len(lines))\n```",
		},
		OnSend: func(chat *ui.Object, text string) { chat.Call("addMessage", "You wrote: "+text) }})
}

func drawTreeAndTables() {
	caption("TREE, TABLES, RANGE")
	ui.TreeView(ui.TreeViewProps{Value: "app", Width: fieldWidth, Height: "auto",
		Items: []ui.P{
			{"key": "src", "text": "src", "expanded": 1, "children": []any{ui.P{"key": "app", "text": "app.py", "description": "4 KB"}, "lib"}},
			{"key": "readme", "text": "README.md", "tag": "new"},
		},
		OnSelect: func(tree *ui.Object, node struct{ Text string }) { tell("File: " + node.Text) }})
	ui.SmartTable(ui.SmartTableProps{ItemLineCount: 3, Width: "100%", Height: "auto", ItemDataList: people,
		TitleDataList: []ui.P{
			{"name": "Name", "dataTitle": "name", "width": 150, "sortable": 1},
			{"name": "City", "dataTitle": "city", "width": 120, "sortable": 1},
			{"name": "Orders", "dataTitle": "orders", "width": 80, "dataType": "integer", "sortable": 1},
		},
		OnSelect: func(table *ui.Object, row person) { tell("Row: " + row.Name) }})
	ui.DataTable(ui.DataTableProps{Rows: people, Width: "100%", Height: 124,
		Columns: []ui.P{
			{"key": "name", "title": "Name", "width": 150},
			{"key": "city", "title": "City", "width": 120},
			{"key": "orders", "title": "Orders", "width": 80, "align": "right", "type": "number"},
		},
		OnSelect: func(table *ui.Object, row *struct{ Data person }) {
			if row == nil {
				tell("No person")
				return
			}
			tell("Person: " + row.Data.Name)
		}})
	ui.RangeSlider(ui.RangeSliderProps{LowValue: 20, HighValue: 70, Width: fieldWidth,
		OnChangeEnd: func(slider *ui.Object) { tell(fmt.Sprintf("Range: %v - %v", slider.Get("low"), slider.Get("high"))) }})
}

func drawChart() {
	caption("CHART")
	ui.GraphBox(ui.GraphBoxProps{Type: "line", TitleText: "Visits", Labels: []string{"Jun", "Jul", "Aug", "Sep", "Oct", "Nov"}, Width: "100%", Height: 220,
		Datasets: []ui.P{
			{"label": "2026", "data": []int{12, 19, 7, 15, 22, 18}},
			{"label": "2025", "data": []int{8, 11, 9, 12, 14, 13}},
		}})
}

func drawOverThePage() {
	caption("OVER THE PAGE")
	sidePanel := ui.SidePanel(ui.SidePanelProps{Side: "right", PanelSize: 300, TitleText: "Filters"})
	sidePanel.Body(func() {
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Only the new ones", Checked: true})
		ui.CheckBox(ui.CheckBoxProps{LabelText: "With a picture"})
	})
	ui.HGroup(ui.Props{Width: "100%", Height: "auto", Align: "left center", Gap: 8, Wrap: true}, func() {
		hint := ui.SmartButton(ui.SmartButtonProps{LabelText: "A hint"})
		ui.Tooltip(ui.TooltipProps{Target: hint, HintText: "Save the file (Ctrl+S)"})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "A dialog", OnClick: askDelete})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "A toast", OnClick: func() { ui.ToastSuccess("report.pdf is saved.", "Saved") }})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "A side panel", OnClick: func() { sidePanel.Call("open") }})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Loading", OnClick: showLoading})
	})
}

func askDelete() {
	if ui.Confirm("report.pdf goes to the trash.", ui.DialogOptions{Title: "Delete the file?", ConfirmText: "Delete", Danger: true}) {
		tell("The file is deleted.")
	}
}

// A loading screen over the whole window; close() lets it go after loadingTime.
func showLoading() {
	loading := ui.LoadingScreen(ui.LoadingScreenProps{TitleText: "My App", MessageText: "Loading the data...", MinDuration: loadingTime, RemoveOnClose: true})
	loading.Call("close")
}
