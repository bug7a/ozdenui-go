/* Bismillah */

// The gallery: the main components on one page, in two columns of cards. The line at the bottom says what was changed
// last. OZDENUI_THEME=dark opens it dark; the Dark / Light button of the ToolBar draws it again in the other theme.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	ui "github.com/bug7a/ozdenui-go"
)

const cardWidth = 440

// The colors of this page in each theme. (The components take their look from the theme of the window.)
type palette struct{ page, card, soft, line string }

var (
	lightColors = palette{page: "#F4F5F7", card: "#FFFFFF", soft: "#5A6270", line: "#E1E3E6"}
	darkColors  = palette{page: "#141413", card: "#232322", soft: "rgba(255, 255, 255, 0.6)", line: "rgba(255, 255, 255, 0.12)"}
)

type person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

type product struct {
	Name    string `json:"name"`
	Price   int    `json:"price"`
	InStock bool   `json:"in_stock"`
}

var people = []person{
	{1, "Ayşe Yılmaz", "İstanbul", 34},
	{2, "Mehmet Kaya", "Ankara", 41},
	{3, "Zeynep Demir", "İzmir", 28},
	{4, "Can Öztürk", "Bursa", 52},
}

var columns = []ui.P{
	{"name": "Name", "dataTitle": "name", "width": 180, "sortable": 1},
	{"name": "City", "dataTitle": "city", "width": 120, "sortable": 1},
	{"name": "Age", "dataTitle": "age", "width": 70, "dataType": "integer", "sortable": 1},
}

var products = []product{{"Honey", 150, true}, {"Olive oil", 420, false}, {"Tea", 95, true}}

const picture = `<svg xmlns="http://www.w3.org/2000/svg" width="320" height="200"><rect width="320" height="200" fill="#DCEFF1"/>` +
	`<circle cx="160" cy="100" r="60" fill="#5BBAC1"/></svg>`

const code = "def hello(name):\n    print(\"Hello\", name)  # a comment\n\nhello(\"world\")\n"

// A page of HTML with its own script: a click sends a message to the program (OnMessage of the WebView).
const webPage = `<body style="margin: 16px; font: 15px -apple-system, 'Segoe UI', sans-serif">
<p style="margin: 0 0 12px">An HTML page with its own script.</p><button id="send">Send to Go</button>
<script>let count = 0; send.onclick = function () { parent.postMessage({count: ++count}, "*"); };</script></body>`

const markdownNote = "### A small report\n\nThe **sales** went up. Read the [notes](https://example.com/notes).\n\n" +
	"| Product | Total |\n|:---|---:|\n| Honey | 3,150 TL |\n\n- [x] Write it\n- [ ] Send it\n\n" +
	"```python\ntotal = sum(rows)\n```\n"

// The state that stays when the page is drawn again for the other theme.
var (
	wrapLines = true
	size      = "M"
)

var (
	colors palette
	status *ui.Object
	logs   *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "OzdenUI gallery", Width: 1000, Height: 820})
}

// A small picture made here (an SVG): a ground and a sun.
func pictureOf(ground, sun string) string {
	return svgAddress(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200"><rect width="200" height="200" fill="%s"/>`+
		`<circle cx="100" cy="90" r="46" fill="%s"/><rect y="150" width="200" height="50" fill="%s" opacity="0.5"/></svg>`, ground, sun, sun))
}

// WHY not ui.Picture([]byte(svg)): it names an SVG "text/xml", and the page does not draw it as a picture.
func svgAddress(svg string) string {
	return "data:image/svg+xml;utf8," + strings.NewReplacer("#", "%23", "\"", "'", "\n", " ").Replace(svg)
}

func start(page *ui.Object) {
	colors = lightColors
	if ui.IsDark() {
		colors = darkColors
	}
	page.Set("color", colors.page)
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{
			{"text": "File", "items": []any{
				ui.P{"text": "New", "shortcut": "Ctrl+N", "onClick": func() { tell("New was chosen.") }},
				"-",
				ui.P{"text": "Recent", "items": []string{"a.txt", "b.txt"}},
			}},
			{"text": "View", "items": []ui.P{{"text": "Word wrap", "key": "wrap", "checked": wrapLines, "onClick": toggleWrap}}},
			{"text": "Help", "items": []ui.P{{"text": "About", "onClick": func() { tell("About was chosen.") }}}},
		}})
		themeText := map[bool]string{false: "Dark", true: "Light"}[ui.IsDark()] // it names the other theme
		ui.ToolBar(ui.ToolBarProps{OnClick: pressTool, Items: []any{
			ui.P{"key": "new", "icon": "new", "hint": "New"},
			ui.P{"key": "save", "icon": "save", "hint": "Save"},
			"-",
			ui.P{"key": "bold", "icon": "bold", "hint": "Bold", "toggle": true},
			ui.P{"key": "left", "icon": "alignLeft", "hint": "Align left", "group": "align", "checked": true},
			ui.P{"key": "right", "icon": "alignRight", "hint": "Align right", "group": "align"},
			"-",
			ui.P{"key": "insert", "text": "Insert", "items": []string{"A picture", "A table"}},
			" ",
			ui.P{"key": "theme", "text": themeText},
			ui.P{"key": "run", "icon": "play", "text": "Run", "shortcut": "F9"},
		}})
		ui.ScrollArea(ui.ScrollAreaProps{Align: "center top", Gap: 20, Padding: 30}, drawCards)
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Change something.", Items: []any{"UTF-8", ui.P{"key": "count", "text": "53 components"}}})
	})
}

func tell(text string) {
	status.Set("messageText", text)
}

func pressTool(bar *ui.Object, item ui.P) {
	if item["key"] == "theme" {
		ui.SetTheme(map[bool]string{false: "dark", true: "light"}[ui.IsDark()]) // start() runs again
		return
	}
	name := item["hint"]
	if name == nil {
		name = item["text"]
	}
	tell(fmt.Sprintf("%v, pressed: %v", name, bar.Get("checkedKeys")))
}

func toggleWrap(bar *ui.Object) {
	wrapLines = !wrapLines
	logs.Set("wrapLines", wrapLines)
	bar.Call("setItemChecked", "wrap", wrapLines)
	tell(fmt.Sprintf("Word wrap: %v", wrapLines))
}

// A card with a title; the objects of content go under the title.
func card(title string, content func()) {
	ui.VGroup(ui.Props{Width: cardWidth, Height: "auto", Align: "left top", Gap: 14, Padding: 20, Color: colors.card, Round: 12}, func() {
		ui.Label(ui.Props{Text: title, FontSize: 13, Bold: true, TextColor: colors.soft})
		content()
	})
}

// A column of cards: a new card goes into the shorter column.
func column(cards func()) {
	ui.VGroup(ui.Props{Width: cardWidth, Height: "auto", Align: "left top", Gap: 20}, cards)
}

func drawCards() {
	ui.Label(ui.Props{Text: "OzdenUI: the components", FontSize: 24, Bold: true})
	// WHY one row of two columns: a new row started under the taller column and left an empty area beside the shorter one.
	ui.HGroup(ui.Props{Width: "auto", Height: "auto", Align: "left top", Gap: 20}, func() {
		column(drawFirstColumn)
		column(drawSecondColumn)
	})
}

func drawFirstColumn() {
	card("SmartButton", func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan", OnClick: func() { ui.ToastSuccess("Saved.") }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Cancel", OnClick: func() { tell("Cancel was clicked.") }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Delete", Tone: "orange", OnClick: askDelete})
		})
	})

	card("CheckBox, Toggle, RadioButton", func() {
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Remember me", Checked: true,
			OnChange: func(box *ui.Object) { tell(fmt.Sprintf("Remember me: %v", box.Checked())) }})
		ui.HGroup(ui.Props{Hug: true, Gap: 10, Align: "left center"}, func() {
			ui.Toggle(ui.ToggleProps{Value: true, OnChange: func(toggle *ui.Object) { tell(fmt.Sprintf("Notifications: %v", toggle.Bool("value"))) }})
			ui.Label(ui.Props{Text: "Notifications", FontSize: 15})
		})
		ui.HGroup(ui.Props{Hug: true, Gap: 18}, func() {
			for _, choice := range []string{"S", "M", "L"} {
				ui.RadioButton(ui.RadioButtonProps{Group: "size", Value: choice, LabelText: choice, Checked: choice == size,
					OnChange: func() {
						size = choice
						tell("Size: " + size)
					}})
			}
		})
	})

	card("SelectBox, SliderField", func() {
		ui.SelectBox(ui.SelectBoxProps{Options: []string{"İstanbul", "Ankara", "İzmir", "Bursa"}, Value: "Ankara", Width: "100%",
			OnChange: func(box *ui.Object) { tell(fmt.Sprintf("City: %v", box.Value())) }})
		ui.SliderField(ui.SliderFieldProps{LabelText: "Volume", Value: 6, MinValue: 0, MaxValue: 10, Width: "100%",
			OnChangeEnd: func(slider *ui.Object) { tell(fmt.Sprintf("Volume: %v", slider.Value())) }})
	})

	card("ProgressBar", func() {
		bar := ui.ProgressBar(ui.ProgressBarProps{TitleText: "Uploading", Value: 40, Width: "100%"})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "+10", OnClick: func() { bar.Call("addValue", 10) }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Reset", OnClick: func() { bar.SetValue(0) }})
		})
	})

	card("LogView", func() {
		logs = ui.LogView(ui.LogViewProps{Height: 170, WrapLines: wrapLines, Lines: []ui.P{
			{"text": "Started", "level": "info"},
			{"text": "14 files found", "level": "success"},
			{"text": "A file is very large", "level": "warn"},
			{"text": "A file can not be read", "level": "error"},
		}})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Write a line", OnClick: func() { logs.Call("info", "A new line") }})
	})

	card("Splitter", func() {
		ui.Box(ui.Props{Width: "100%", Height: 120, Color: "transparent", Border: 1, BorderColor: colors.line, Round: 8}, func() {
			split := ui.Splitter(ui.SplitterProps{FirstSize: "40%"})
			split.First(func() {
				ui.Label(ui.Props{Text: "Drag the line", FontSize: 14, TextColor: colors.soft, Padding: 10})
			})
			split.Second(func() {
				ui.Label(ui.Props{Text: "The other area", FontSize: 14, TextColor: colors.soft, Padding: 10})
			})
		})
	})

	card("Tabs", func() {
		tabs := ui.Tabs(ui.TabsProps{Tabs: []string{"General", "Advanced"}, OnChange: func(t *ui.Object) { tell(fmt.Sprintf("Tab: %v", t.Value())) }})
		tabs.Panel("General", func() { ui.Label(ui.Props{Text: "The general settings.", FontSize: 14}) })
		tabs.Panel("Advanced", func() { ui.Label(ui.Props{Text: "The advanced settings.", FontSize: 14}) })
	})

	card("TextField, EmailField, PasswordField, NumberField", func() {
		ui.TextField(ui.TextFieldProps{TitleText: "Name", Placeholder: "Ada", Width: "100%", OnChange: func(field *ui.Object) { tell("Name: " + field.Text()) }})
		ui.EmailField(ui.EmailFieldProps{TitleText: "E-mail", Width: "100%"})
		ui.PasswordField(ui.PasswordFieldProps{TitleText: "Password", Width: "100%"})
		ui.NumberField(ui.NumberFieldProps{TitleText: "Count", Value: 3, Width: "100%", OnChange: func(field *ui.Object) { tell("Count: " + field.Text()) }})
	})

	card("SearchInput, Slider, Stepper", func() {
		ui.SearchInput(ui.SearchInputProps{Width: "100%", OnChange: func(search *ui.Object) { tell("Search: " + search.Text()) }})
		ui.Slider(ui.SliderProps{Value: 30, Width: "100%", OnChangeEnd: func(slider *ui.Object) { tell(fmt.Sprintf("Slider: %v", slider.Value())) }})
		ui.Stepper(ui.StepperProps{Value: 2, MinValue: 0, MaxValue: 10, OnChange: func(stepper *ui.Object) { tell(fmt.Sprintf("Stepper: %v", stepper.Value())) }})
	})

	card("StepBar", func() {
		ui.StepBar(ui.StepBarProps{Steps: []string{"Cart", "Address", "Payment"}, CurrentStep: 1,
			OnStepClick: func(bar *ui.Object, index int) {
				bar.Call("setStep", index)
				tell(fmt.Sprintf("Step: %d", index+1))
			}})
	})

	card("RangeSlider", func() {
		ui.RangeSlider(ui.RangeSliderProps{LowValue: 200, HighValue: 800, MaxValue: 1000, Step: 50, Width: "100%", FormatValue: "{:,.0f} TL",
			OnChangeEnd: func(slider *ui.Object) { tell(fmt.Sprintf("Range: %v - %v", slider.Get("low"), slider.Get("high"))) }})
	})

	card("DatePicker, TimePicker, ColorPicker", func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 10}, func() {
			ui.DatePicker(ui.DatePickerProps{Date: "2026-11-20", Width: 140, OnChange: func(picker *ui.Object) { tell("Date: " + picker.Str("value")) }})
			ui.TimePicker(ui.TimePickerProps{Time: "19:30", Width: 110, OnChange: func(picker *ui.Object) { tell("Time: " + picker.Str("value")) }})
			ui.ColorPicker(ui.ColorPickerProps{SelectedColor: "#3871E0", Width: 130, OnChange: func(picker *ui.Object) { tell("Color: " + picker.Str("value")) }})
		})
	})

	card("GridGroup", func() {
		grid := ui.GridGroup(ui.GridGroupProps{Columns: 3, Gap: 8})
		grid.Add(func() {
			for number := 1; number <= 6; number++ {
				button := ui.SmartButton(ui.SmartButtonProps{LabelText: fmt.Sprintf("Cell %d", number), Width: "100%",
					OnClick: func() { tell("A cell was clicked.") }})
				if number == 4 {
					grid.Call("place", button, ui.P{"columnSpan": 2})
				}
			}
		})
	})

	card("UrlField, PhoneField, CurrencyField", func() {
		ui.UrlField(ui.UrlFieldProps{TitleText: "Site", Width: "100%", OnChange: func(field *ui.Object) { tell("Site: " + field.Text()) }})
		ui.PhoneField(ui.PhoneFieldProps{TitleText: "Phone", Width: "100%", OnChange: func(field *ui.Object) { tell("Phone: " + field.Text()) }})
		ui.CurrencyField(ui.CurrencyFieldProps{TitleText: "Price", UnitText: "$", Width: "100%", OnChange: func(field *ui.Object) { tell("Price: " + field.Text()) }})
	})

	card("SortableList", func() {
		ui.SortableList(ui.SortableListProps{Width: "100%", Height: "auto",
			Items:    []any{"Sales", "Stock", ui.P{"key": "pay", "text": "Payments", "description": "Every Monday"}},
			OnChange: func(list *ui.Object) { tell(fmt.Sprintf("Order: %v", list.Get("order"))) }})
	})

	card("WebView", func() {
		ui.WebView(ui.WebViewProps{HtmlText: webPage, Width: "100%", Height: 110,
			OnMessage: func(view *ui.Object, data struct{ Count int }) { tell(fmt.Sprintf("WebView: %d", data.Count)) }})
	})

	card("ChatView", func() {
		ui.ChatView(ui.ChatViewProps{Width: "100%", Height: 300,
			Messages: []any{ui.P{"role": "user", "text": "Hello!"}, "Hi! **How** can I help?"},
			OnSend:   func(chat *ui.Object, text string) { chat.Call("addMessage", "You wrote: "+text) }})
	})
}

func drawSecondColumn() {
	card("TextAreaField", func() {
		ui.TextAreaField(ui.TextAreaFieldProps{TitleText: "Note", Placeholder: "Write something...", Width: "100%", Height: 110,
			OnChange: func(note *ui.Object) { tell(fmt.Sprintf("Note: %d letters", utf8.RuneCountInString(note.Text()))) }})
	})

	card("ItemList", func() {
		ui.ItemList(ui.ItemListProps{Width: "100%", Height: 170,
			Items: []ui.P{
				{"key": "a", "text": "Apples", "description": "Red"},
				{"key": "p", "text": "Pears", "description": "Green"},
				{"key": "g", "text": "Grapes", "description": "Purple"},
			},
			OnChange: func(list *ui.Object) { tell(fmt.Sprintf("Fruit: %v", list.Get("selectedKeys"))) }})
	})

	card("ItemList: pictures", func() {
		ui.ItemList(ui.ItemListProps{Layout: "grid", Width: "100%", Height: "auto", SelectedKeys: []string{"sun"},
			StyleName: "gallery",                                         // WHY no "galleryDark": a dark window takes it
			Style:     ui.P{"tile": ui.P{"width": 120, "iconSize": 104}}, // WHY: three in the row of this card
			Items: []ui.P{
				{"key": "sea", "text": "Sea", "iconFile": pictureOf("#DCEFF1", "#3F8D94")},
				{"key": "sun", "text": "Sun", "iconFile": pictureOf("#FBEBD3", "#E0913A")},
				{"key": "wood", "text": "Wood", "iconFile": pictureOf("#E3EEDC", "#5E8C4A")},
			},
			OnChange: func(list *ui.Object) { tell(fmt.Sprintf("Picture: %v", list.Get("selectedKeys"))) }})
	})

	card("Canvas", func() {
		ui.Canvas(ui.CanvasProps{Width: "100%", Height: 170,
			Shapes: []ui.P{
				{"key": "ground", "type": "rect", "x": 0, "y": 125, "width": cardWidth, "height": 45, "fill": "#E3EEDC"},
				{"key": "hill", "type": "polygon", "points": [][]int{{20, 125}, {80, 55}, {140, 125}}, "fill": "#9DBE8C"},
				{"key": "path", "type": "line", "x1": 170, "y1": 150, "x2": 330, "y2": 150, "color": "#5E8C4A", "lineWidth": 3, "dash": []int{8, 5}},
				{"key": "sun", "type": "circle", "x": 300, "y": 55, "radius": 24, "fill": "#E0913A", "draggable": 1},
				{"key": "title", "type": "text", "x": 16, "y": 12, "text": "Drag the sun", "fontSize": 15, "bold": 1},
			},
			OnShapeMove: func(canvas *ui.Object, shape struct{ Key string }) { tell("Moved: " + shape.Key) },
			OnClick: func(canvas *ui.Object, x, y float64, shape *struct{ Key string }) {
				if shape == nil {
					tell("Canvas: empty")
					return
				}
				tell("Canvas: " + shape.Key)
			}})
	})

	card("SmartTable", func() {
		ui.SmartTable(ui.SmartTableProps{TitleDataList: columns, ItemDataList: people, Width: "100%", Height: 230, FitLineCount: true,
			OnSelect: func(table *ui.Object, row person) { tell("Row: " + row.Name) }})
	})

	card("ImageView", func() {
		ui.ImageView(ui.ImageViewProps{Height: 170, ImageFile: svgAddress(picture), EnlargeToFit: true, AltText: "A circle on a square"})
	})

	card("CodeEditor", func() {
		ui.CodeEditor(ui.CodeEditorProps{Height: 150, Language: "python", Value: code})
	})

	card("GraphBox", func() {
		ui.GraphBox(ui.GraphBoxProps{Type: "bar", Labels: []string{"Jan", "Feb", "Mar", "Apr"}, Values: []int{1200, 1900, 700, 1500},
			Width: "100%", Height: 210, FormatValue: "{:,.0f} TL",
			OnClick: func(graph *ui.Object, item struct{ Label string }) { tell("Bar: " + item.Label) }})
	})

	card("TreeView, ContextMenu", func() {
		tree := ui.TreeView(ui.TreeViewProps{Width: "100%", Height: 130,
			Items: []ui.P{
				{"key": "photos", "text": "photos", "expanded": 1, "children": []string{"sea.jpg", "hill.jpg"}},
				{"key": "notes", "text": "notes.txt", "description": "4 KB"},
			},
			OnSelect: func(t *ui.Object, node struct{ Text string }) { tell("Node: " + node.Text) }})
		ui.ContextMenu(ui.ContextMenuProps{Target: tree, Items: []any{
			ui.P{"text": "Open them all", "onClick": func() { tree.Call("expandAll") }},
			"-",
			ui.P{"text": "Close them all", "onClick": func() { tree.Call("collapseAll") }},
		}})
	})

	card("Tooltip", func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			hintButton := ui.SmartButton(ui.SmartButtonProps{LabelText: "A hint"})
			ui.Tooltip(ui.TooltipProps{Target: hintButton, HintText: "A hint on the button"})
			var bubbleButton *ui.Object
			bubbleButton = ui.SmartButton(ui.SmartButtonProps{LabelText: "A bubble", OnClick: func() { showBubble(bubbleButton) }})
			popoverButton := ui.SmartButton(ui.SmartButtonProps{LabelText: "A popover"})
			ui.Tooltip(ui.TooltipProps{Target: popoverButton, Trigger: "click", HintPosition: "bottom", Look: bubbleLook(),
				HintText: "It stays until a click outside."})
		})
	})

	card("Accordion", func() {
		folds := ui.Accordion(ui.AccordionProps{
			Sections: []any{ui.P{"key": "more", "text": "More settings", "description": "2", "opened": 1}, "Closed at the start"},
			OnToggle: func(a *ui.Object) { tell(fmt.Sprintf("Open: %v", a.Get("openKeys"))) }})
		folds.Section("more", func() { ui.CheckBox(ui.CheckBoxProps{LabelText: "Save it by itself", Checked: true}) })
		folds.Section("Closed at the start", func() {
			ui.Label(ui.Props{Text: "It was made while it was closed.", FontSize: 14, TextColor: colors.soft})
		})
	})

	card("MarkdownView", func() {
		ui.MarkdownView(ui.MarkdownViewProps{Height: 250, MarkdownText: markdownNote,
			OnLinkClick: func(view *ui.Object, address string) { tell("Link: " + address) },
			Style:       ui.P{"box": ui.P{"color": "transparent"}, "page": ui.P{"padding": []int{0, 0}}}})
	})

	card("JsonView", func() {
		ui.JsonView(ui.JsonViewProps{Height: 170, ExpandLevel: 2,
			Value:    json.RawMessage(`{"status": "ok", "count": 2, "users": [{"name": "Ada", "active": true}]}`), // WHY JSON: a map of Go has no order
			OnSelect: func(view *ui.Object, entry struct{ PathText string }) { tell("Path: " + entry.PathText) },
			Style:    ui.P{"box": ui.P{"color": "transparent"}}})
	})

	card("JsonEditor", func() {
		ui.JsonEditor(ui.JsonEditorProps{Height: 340, Value: inOrder(products), TitleText: "Products", OnChange: countProducts})
	})

	card("FilePicker", func() {
		ui.FilePicker(ui.FilePickerProps{Width: "100%", ZoneHeight: 110, Multiple: true,
			OnChange: func(picker *ui.Object) {
				files, _ := picker.Get("files").([]any)
				tell(fmt.Sprintf("Files: %d", len(files)))
			}})
	})
}

// The frame of a bubble: the bubble of the kit is white with no line, so on a white card only its text was seen.
func bubbleLook() ui.P {
	if ui.IsDark() {
		return ui.P{"color": "#2A2A29", "textColor": "rgba(255, 255, 255, 0.86)", "border": 1, "borderColor": "rgba(255, 255, 255, 0.14)",
			"round": 10, "padding": []int{16, 14}, "shadow": "0 10px 30px rgba(0, 0, 0, 0.4)"}
	}
	return ui.P{"color": "#FFFFFF", "textColor": "rgba(0, 0, 0, 0.85)", "border": 1, "borderColor": "rgba(0, 0, 0, 0.12)",
		"round": 10, "padding": []int{16, 14}, "shadow": "0 10px 30px rgba(0, 0, 0, 0.12)"}
}

// A bubble at the button that goes by itself.
func showBubble(button *ui.Object) {
	ui.ShowHint(button, "It goes by itself.", 1500)
}

// The value as JSON: a struct goes to the page as a map, and the columns of the editor would come in the order of A-Z.
func inOrder(value any) json.RawMessage {
	text, _ := json.Marshal(value)
	return text
}

// OnChange of the JsonEditor: the whole value, as a list of products.
func countProducts(editor *ui.Object, list []product) {
	inStock := 0
	for _, item := range list {
		if item.InStock {
			inStock++
		}
	}
	tell(fmt.Sprintf("Products: %d, %d in stock", len(list), inStock))
}

func askDelete() {
	if ui.Confirm("The file will be deleted.", ui.DialogOptions{Title: "Delete it?", ConfirmText: "Delete", Danger: true}) {
		tell("Deleted.")
	} else {
		tell("Not deleted.")
	}
}
