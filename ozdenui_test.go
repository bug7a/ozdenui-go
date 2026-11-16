/* Bismillah */

/*

ozdenui_test.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The checks of the module: go test .

TestWindow runs one program in a hidden Chrome (or Edge) and plays it as a user would: the commands really make
basic.js objects, and the events of the page really come back. Without such a browser it is skipped; the other tests
need nothing.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestKitName(t *testing.T) {
	names := map[string]string{"LabelText": "labelText", "ID": "id", "URLText": "urlText", "X": "x", "OnClick": "onClick", "Css": "css"}
	for field, want := range names {
		if got := kitName(field); got != want {
			t.Errorf("kitName(%q) = %q, want %q", field, got, want)
		}
	}
	if eventNameOf("onKeyDown") != "keydown" || !isCallbackName("onChange") || isCallbackName("only") {
		t.Error("the names of the events")
	}
}

func TestReadProps(t *testing.T) {
	read := readProps(Props{Text: "Hi", Gap: 0, Visible: false, Fill: true, Extra: P{"textSize": 12}})
	want := P{"text": "Hi", "gap": 0, "visible": false, "textSize": 12}
	if !reflect.DeepEqual(read.values, want) || !read.fill {
		t.Errorf("readProps: %v fill %v", read.values, read.fill)
	}
	// A zero that is given (Gap: 0, Checked: false) is a value; a field that is not written is not.
	typed := readProps(CheckBoxProps{LabelText: Empty, Checked: false})
	if len(typed.values) != 2 || typed.values["checked"] != false {
		t.Errorf("readProps of a component: %v", typed.values)
	}
}

func TestFormat(t *testing.T) {
	cases := map[string]string{
		"{:,.2f} TL": `{"prefix":"","suffix":" TL","thousands":1,"decimals":2,"percent":0}`,
		"$ {:,}":     `{"prefix":"$ ","suffix":"","thousands":1,"decimals":null,"percent":0}`,
		"{:.1%}":     `{"prefix":"","suffix":"","thousands":0,"decimals":1,"percent":1}`,
		"{} kg":      `{"prefix":"","suffix":" kg","thousands":0,"decimals":null,"percent":0}`,
	}
	for pattern, want := range cases {
		text, _ := json.Marshal(Format(pattern))
		if string(text) != want {
			t.Errorf("Format(%q) = %s, want %s", pattern, text, want)
		}
	}
}

type testRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestPicture(t *testing.T) {
	svg := Picture([]byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"></svg>`))
	png := Picture([]byte("\x89PNG\r\n\x1a\n"))
	if !strings.HasPrefix(svg, "data:image/svg+xml;base64,") || !strings.HasPrefix(png, "data:image/png;base64,") {
		t.Errorf("Picture of bytes: %.40s %.40s", svg, png)
	}
	if Picture("https://example.com/a.png") != "https://example.com/a.png" {
		t.Error("an address of the web stays as it is")
	}
}

func TestCallHandler(t *testing.T) {
	obj := &Object{typeName: "Label"}
	got := ""
	callHandler(func() { got += "a" }, []any{obj, 3.0})
	callHandler(func(self *Object) { got += self.typeName }, []any{obj, 3.0})
	callHandler(func(self *Object, count int, name string, on bool) { got += fmt.Sprint(count, name, on) }, []any{obj, 3.0, "x", 1.0})
	callHandler(func(self *Object, row testRow) { got += row.Name + fmt.Sprint(row.ID) }, []any{obj, map[string]any{"id": 7.0, "name": "Ali"}})
	callHandler(func(self *Object, missing string, part *Object) { got += "[" + missing + "]" + fmt.Sprint(part == nil) }, []any{obj})
	callHandler(func(values ...any) { got += fmt.Sprint(len(values)) }, []any{obj, 1.0, 2.0})
	callHandler(func(item struct {
		Key     string
		Checked bool `json:"checked"`
		Tags    []string
	}) {
		got += fmt.Sprint(item.Key, item.Checked, item.Tags)
	}, []any{map[string]any{"key": "k", "checked": 1.0, "tags": []any{"a", "b"}}})
	if got != "aLabel3xtrueAli7[]true3ktrue [a b]" {
		t.Errorf("callHandler: %q", got)
	}
}

func TestKitTable(t *testing.T) {
	for name, spec := range kit {
		for _, callback := range spec.Sync {
			if !spec.Params[callback] {
				t.Errorf("%s: the state callback %s is not a param of the kit (renamed?)", name, callback)
			}
		}
		for short, long := range spec.Aliases {
			if !spec.Params[long] && !objectProps[long] && long != "styleName" {
				t.Errorf("%s: the short name %s names %s, which is not a param of the kit", name, short, long)
			}
		}
		for key := range spec.Defaults {
			if !spec.Params[key] && !objectProps[key] && key != "styleName" {
				t.Errorf("%s: the default %s is not a param of the kit", name, key)
			}
		}
	}
	if kit["TextField"].Setters["value"].Method != "setValue" || !kit["CheckBox"].Setters["checked"].Silent {
		t.Error("the set methods of the kit were not read")
	}
}

// The files of the module start with Bismillah and carry the notice (the rule of the project).
func TestHeaders(t *testing.T) {
	files, _ := os.ReadDir(".")
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		content, _ := os.ReadFile(name)
		text := string(content)
		start := text
		if strings.HasPrefix(text, "//go:build") {
			start = text[strings.Index(text, "\n\n")+2:]
		}
		if !strings.HasPrefix(start, "/* Bismillah */") {
			t.Errorf("%s does not start with Bismillah", name)
		}
		if !strings.Contains(text, "SPDX-License-Identifier: Apache-2.0") {
			t.Errorf("%s has no license line", name)
		}
	}
}

// The texts people and assistants read name only what the module has.
func TestDocuments(t *testing.T) {
	readme, _ := os.ReadFile("README.md")
	table := string(readme)
	table = table[strings.Index(table, "| For | Functions |"):]
	table = table[:strings.Index(table, "\n\n")]
	for _, match := range regexp.MustCompile("`([A-Z][A-Za-z]+)`").FindAllStringSubmatch(table, -1) {
		if kit[match[1]] == nil {
			t.Errorf("README.md: %s is not a component of the kit copy", match[1])
		}
	}
	rules, _ := os.ReadFile("cmd/ozdenui/AGENTS.md")
	text := string(rules)
	program := text[strings.Index(text, "```go\npackage main")+6:]
	program = program[:strings.Index(program, "```")]
	example, _ := os.ReadFile("examples/notes/main.go")
	if !strings.HasSuffix(string(example), program) {
		t.Error("the program of cmd/ozdenui/AGENTS.md is not examples/notes/main.go (the example is the one that is built and pictured)")
	}
	command, _ := os.ReadFile("cmd/ozdenui/main.go")
	if !strings.Contains(string(command), `version       = "`+Version+`"`) {
		t.Error("cmd/ozdenui/main.go: its version is not the Version of the module")
	}
}

// *** THE WINDOW:

type windowTest struct {
	t   *testing.T
	app *application
}

// Lets the page and the tasks of the program work for a while.
func (w *windowTest) pause(milliseconds int) {
	done := false
	After(milliseconds, func() { done = true })
	w.app.waitFor(func() bool { return done })
}

func (w *windowTest) click(obj *Object) {
	Evaluate(obj.JS() + ".elem.click(); 0")
	w.pause(60)
}

func (w *windowTest) equal(what string, got, want any) {
	w.t.Helper()
	if !reflect.DeepEqual(got, want) {
		w.t.Errorf("%s: got %#v, want %#v", what, got, want)
	}
}

func TestWindow(t *testing.T) {
	if findBrowser() == "" {
		t.Skip("no Chrome, Edge or Chromium on this computer")
	}
	t.Setenv("OZDENUI_WINDOW", "headless")
	finished := false
	guard := time.AfterFunc(60*time.Second, func() {
		t.Error("the window test did not end in a minute")
		Close()
	})
	defer guard.Stop()
	starts := 0
	Run(func(page *Object) {
		w := &windowTest{t: t, app: current}
		starts++
		if starts == 2 {
			// SetTheme() made the page again: start() runs again, in the dark window
			defer Close()
			box := CheckBox(CheckBoxProps{LabelText: "Dark"})
			w.equal("the theme of the second start", []any{Theme(), box.Str("styleName"), page.Measure("color")}, []any{"dark", "dark", darkPageColor})
			finished = true
			return
		}
		w.objects(page)
		w.components()
		w.dialogs()
		w.tasks()
		SetTheme("dark")
	}, Window{Title: "Test", Width: 900, Height: 700})
	if !finished {
		t.Error("the program of the test did not reach its end")
	}
}

func (w *windowTest) objects(page *Object) {
	w.equal("the size of the page", []float64{page.Width(), page.Height()}, []float64{900, 700})
	w.equal("the title of the page", Evaluate("document.title"), "Test")
	clicks := 0
	var label, button, input, inner *Object
	group := VGroup(Props{Gap: 10, Padding: 20, Align: "left top"}, func() {
		label = Label(Props{Text: "<b>plain</b>", FontSize: 20})
		button = Button(Props{Text: "Press", Width: 120, OnClick: func(self *Object, event Event) {
			clicks++
			self.SetText(fmt.Sprint("Pressed ", clicks, " ", event.Name))
		}})
		input = Input(Props{Width: 200, Placeholder: "A name"})
		Box(Props{Width: 100, Height: 50, Color: "#DDD"}, func() {
			inner = Label(Props{Text: "in"})
		})
	})
	w.equal("a text is not HTML", Evaluate(label.JS()+".elem.innerText"), "<b>plain</b>")
	w.equal("the font size", label.Measure("fontSize"), 20.0)
	w.equal("an object in a Box stands at its corner", []any{inner.Measure("left"), inner.Measure("top")}, []any{0.0, 0.0})
	w.equal("a group does not cut", Evaluate(group.JS()+".elem.style.overflow"), "visible")
	w.click(button)
	w.click(button)
	w.equal("the clicks", clicks, 2)
	w.equal("the text of the button", button.Text(), "Pressed 2 click")
	label.Set("html", "<i>x</i>", "textColor", "red")
	w.equal("html is HTML", Evaluate(label.JS()+".elem.querySelector('i') !== null"), true)
	Evaluate(input.JS() + ".inputElement.value = 'Ali'; " + input.JS() + ".inputElement.dispatchEvent(new Event('input', { bubbles: true })); 0")
	w.pause(60)
	w.equal("the text of an Input comes with its event", input.Text(), "Ali")
	group.Add(func() { Label(Props{Text: "later"}) })
	w.equal("an object made later is in its group", Evaluate(group.JS()+".elem.children.length"), 5.0)
	inner.Remove()
	Flush()
	w.equal("a removed object is forgotten", w.app.objects[inner.id] == nil, true)
	group.Remove()
}

func (w *windowTest) components() {
	changes := []string{}
	var check, field, slider, selectBox, tabs, table, toggle, list *Object
	ScrollArea(ScrollAreaProps{Gap: 10, Padding: 20, Align: "left top"}, func() {
		check = CheckBox(CheckBoxProps{LabelText: "Agree", OnChange: func(self *Object) {
			changes = append(changes, fmt.Sprint("check ", self.Checked()))
		}})
		toggle = Toggle(ToggleProps{Value: true})
		field = TextField(TextFieldProps{TitleText: "Name", Value: "Ali", Width: 200})
		slider = Slider(SliderProps{Value: 30, Width: 200})
		selectBox = SelectBox(SelectBoxProps{Options: []string{"Red", "Green"}, Value: "Green"})
		list = ItemList(ItemListProps{Width: 200, Height: 120, SelectedKeys: []string{"b"},
			Items: []P{{"key": "a", "text": "Apple"}, {"key": "b", "text": "Banana"}}})
		tabs = Tabs(TabsProps{Tabs: []P{{"key": "one", "text": "One"}, {"key": "two", "text": "Two"}}})
		tabs.Panel("one", func() { Label(Props{Text: "Panel one"}) })
		table = SmartTable(SmartTableProps{Width: 300, Height: 160, ItemLineCount: 3,
			TitleDataList: []P{{"name": "ID", "dataTitle": "id", "width": 60}, {"name": "NAME", "dataTitle": "name", "width": 140}},
			ItemDataList:  []testRow{{1, "Ali"}, {2, "Can"}}})
		Component("Gauge", P{"value": 72, "width": 120, "height": 100})
	})
	w.equal("the state of the start", []any{check.Checked(), toggle.Checked(), field.Text(), slider.Num("value"), selectBox.Value(), tabs.Value()},
		[]any{false, true, "Ali", 30.0, "Green", "one"})
	w.equal("the selected item of the list", list.Ask("getSelectedKeys"), []any{"b"})
	w.equal("the rows of the table", table.Measure("itemDataList"), []any{map[string]any{"id": 1.0, "name": "Ali"}, map[string]any{"id": 2.0, "name": "Can"}})
	w.equal("a panel of the tabs holds its label", Evaluate(tabs.JS()+".elem.innerText.includes('Panel one')"), true)

	// A change of the user tells the program; a change of the program does not call OnChange.
	Evaluate(check.JS() + ".elem.click(); 0")
	w.pause(80)
	w.equal("the click of the user", []any{check.Checked(), changes}, []any{true, []string{"check true"}})
	check.Set("checked", false)
	field.SetValue("Veli")
	field.Set("titleText", "Full name")
	slider.SetValue(250)
	selectBox.SetValue("Nope")
	w.pause(80)
	w.equal("OnChange is not called by the program", len(changes), 1)
	w.equal("the page keeps what the program gave", []any{check.Measure("checked"), field.Ask("getValue"), field.Measure("titleText")}, []any{0.0, "Veli", "Full name"})
	w.equal("the page corrects a value, and the program reads it", []any{slider.Value(), selectBox.Value()}, []any{100.0, nil})

	// What the reports of the first programs found:
	var low, high, date, number, grid, wide, editor *Object
	VGroup(Props{Hug: true, Gap: 6}, func() {
		low = RadioButton(RadioButtonProps{LabelText: "Low", Group: "level", Value: "low", Checked: true})
		high = RadioButton(RadioButtonProps{LabelText: "High", Group: "level", Value: "high"})
		date = DatePicker(DatePickerProps{Date: "2026-11-16"})
		number = NumberField(NumberFieldProps{Value: 250, Width: 120})
		grid = GridGroup(GridGroupProps{Columns: 2, Width: 200}, func() {
			wide = Label(Props{Text: "wide"})
			Current().Place(wide, P{"columnSpan": 2})
		})
		editor = CodeEditor(CodeEditorProps{Width: 200, Height: 60, Value: "x = 1"})
	})
	Evaluate(high.JS() + ".elem.click(); 0")
	w.pause(80)
	w.equal("the radio buttons of a group", []any{low.Checked(), high.Checked(), RadioValue("level")}, []any{false, true, "high"})
	low.Set("checked", true)
	w.equal("a radio checked by the program", []any{high.Checked(), RadioValue("level")}, []any{false, "low"})
	w.equal("the value of a DatePicker is its date", date.Value(), "2026-11-16")
	date.SetValue("2026-12-01")
	w.pause(80)
	w.equal("SetValue of a DatePicker", []any{date.Value(), Evaluate(date.JS() + ".getDate().getDate()")}, []any{"2026-12-01", 1.0})
	typedNumber := number.JS() + ".input.inputElement"
	Evaluate(typedNumber + ".value = '300'; " + typedNumber + ".dispatchEvent(new Event('input', { bubbles: true })); 0")
	w.pause(80)
	w.equal("the value of a NumberField follows the user", Text(number.Value()), "300")
	w.equal("a span in a grid", Evaluate(wide.JS()+".elem.style.gridColumn.includes('span 2')"), true)
	w.equal("Current() is the grid", grid.Measure("columns"), 2.0)
	editor.Set("theme", "dark")
	w.equal("the theme of one object", editor.Measure("styleName"), "dark")
	hint := ShowHint(low, "A hint", 0)
	w.equal("a hint is on the page", Evaluate("document.body.innerText.includes('A hint')"), true)
	hint.Remove()
	w.equal("a removed hint", Evaluate("document.body.innerText.includes('A hint')"), false)
	w.equal("a field that is right", field.Check(), true)

	// A text typed by the user is read without a question.
	typed := field.JS() + ".input.inputElement"
	Evaluate(typed + ".value = 'Typed'; " + typed + ".dispatchEvent(new Event('input', { bubbles: true })); 0")
	w.pause(80)
	w.equal("the text of a field", []any{field.Text(), field.Value()}, []any{"Typed", "Typed"})
}

func (w *windowTest) dialogs() {
	// The dialog waits for its answer while the window works: a timer of the program presses its button.
	press := func(text string) {
		After(250, func() {
			Evaluate(`Array.from(document.querySelectorAll("button, [role=button]")).filter(function (b) { return b.innerText.trim() === "` + text + `"; }).forEach(function (b) { b.click(); }); 0`)
		})
	}
	press("Confirm")
	w.equal("Confirm: yes", Confirm("Delete it?"), true)
	press("Cancel")
	w.equal("Confirm: no", Confirm("Delete it?"), false)
	press("Later")
	w.equal("Ask", Ask("Save the changes?", "Save", "Later"), "Later")
	press("OK")
	text, ok := Prompt("Name:", DialogOptions{Value: "Ali"})
	w.equal("Prompt", []any{text, ok}, []any{"Ali", true})
	ToastSuccess("Saved.")
	w.pause(100)
	w.equal("a toast is on the page", Evaluate("document.body.innerText.includes('Saved.')"), true)
}

func (w *windowTest) tasks() {
	order := []string{}
	stop := Every(20, func() { order = append(order, "tick") })
	After(70, func() { order = append(order, "once"); stop() })
	go func() { Do(func() { order = append(order, "do") }) }()
	w.pause(200)
	ticks := strings.Count(strings.Join(order, " "), "tick")
	if ticks < 2 || ticks > 5 || !strings.Contains(strings.Join(order, " "), "once") || order[0] != "do" {
		w.t.Errorf("the timers and Do(): %v", order)
	}
	// Another goroutine may not change an object: it is told how.
	told := make(chan any, 1)
	label := Label(Props{Text: "x"})
	go func() {
		defer func() { told <- recover() }()
		label.SetText("no")
	}()
	if problem := <-told; problem == nil || !strings.Contains(fmt.Sprint(problem), "ozdenui.Do") {
		w.t.Errorf("a change from another goroutine: %v", problem)
	}
}
