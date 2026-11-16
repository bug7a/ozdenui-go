/* Bismillah */

/*

OzdenUI for Go Website - Samples - v26.11

The Go code of the page (shown, not run) and the commands of the terminal.
NOTE: SAMPLES.counter is examples/counter/main.go of the module (its picture is images/counter.webp), SAMPLES.own the
code of tools/shots/card (images/card.webp) and SAMPLES.others the code of tools/shots/gauge (images/gauge.webp):
change them together. The pictures are made by tools/make-images.sh from the real windows.
Every sample compiles: put it into a main.go of a module that requires github.com/bug7a/ozdenui-go and run go vet.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

"use strict";

const SAMPLES = {};

// WHY spaces and not tabs (gofmt writes tabs): a tab is 8 columns wide in the code panel, and the lines went out of it
SAMPLES.counter = [
    'package main',
    '',
    'import (',
    '    "strconv"',
    '',
    '    ui "github.com/bug7a/ozdenui-go"',
    ')',
    '',
    'var (',
    '    count = 0',
    '    label *ui.Object',
    ')',
    '',
    'func main() {',
    '    ui.Run(start, ui.Window{Title: "Counter", Width: 360, Height: 260})',
    '}',
    '',
    'func start(page *ui.Object) {',
    '    page.Set("color", "#F4F5F7")',
    '    ui.VGroup(ui.Props{Gap: 16}, func() {',
    '        label = ui.Label(ui.Props{Text: "0", FontSize: 48})',
    '        ui.HGroup(ui.Props{Hug: true, Gap: 10}, func() {',
    '            ui.SmartButton(ui.SmartButtonProps{LabelText: "Minus", OnClick: func() { add(-1) }})',
    '            ui.SmartButton(ui.SmartButtonProps{LabelText: "Plus", Tone: "cyan", OnClick: func() { add(1) }})',
    '        })',
    '    })',
    '}',
    '',
    'func add(step int) {',
    '    count += step',
    '    label.SetText(strconv.Itoa(count))',
    '}',
].join("\n");

SAMPLES.dialogs = [
    'func replaceFile() {',
    '    if ui.Confirm("The old copy will be replaced.",',
    '        ui.DialogOptions{Title: "Replace it?", ConfirmText: "Replace"}) {',
    '        writeTheNewCopy()',
    '        ui.ToastSuccess("Replaced.")',
    '    }',
    '}',
    '',
    '// the text, and false when it is cancelled',
    'name, ok := ui.Prompt("A name for the folder:", ui.DialogOptions{Value: "Photos"})',
    '',
    '// the text of the pressed button',
    'choice := ui.Ask("Save the changes?", "Save", "Don\'t save", "Cancel")',
    '',
    '// the file dialog of the system: a real path, or ""',
    'path := ui.OpenFile(ui.FileOptions{Types: []string{"csv"}})',
].join("\n");

SAMPLES.others = [
    'func start(page *ui.Object) {',
    '    ui.HGroup(ui.Props{Gap: 24}, func() {',
    '        // the function of the component, with its own type',
    '        ui.Gauge(ui.GaugeProps{Value: 72, SubText: "Disk"})',
    '        // the same component by its name, with the params of the kit',
    '        ui.Component("Gauge", ui.P{"value": 35, "subText": "Memory"})',
    '    })',
    '}',
].join("\n");

SAMPLES.own = [
    'package main',
    '',
    'import ui "github.com/bug7a/ozdenui-go"',
    '',
    '// Card is a component of this program: a title, a text and a button.',
    'type Card struct {',
    '    Box   *ui.Object',
    '    title *ui.Object',
    '}',
    '',
    '// NewCard makes a card in the group it is called in.',
    'func NewCard(title, text string, onOpen func(card *Card)) *Card {',
    '    card := &Card{}',
    '    card.Box = ui.VGroup(ui.Props{Width: 220, Height: "auto", Align: "left top",',
    '        Padding: 16, Gap: 8, Color: "white", Round: 12}, func() {',
    '        card.title = ui.Label(ui.Props{Text: title, FontSize: 18, Bold: true})',
    '        ui.Label(ui.Props{Text: text, TextColor: "#5A6270"})',
    '        ui.SmartButton(ui.SmartButtonProps{LabelText: "Open", OnClick: func() { onOpen(card) }})',
    '    })',
    '    return card',
    '}',
    '',
    '// SetTitle changes its title.',
    'func (card *Card) SetTitle(title string) {',
    '    card.title.SetText(title)',
    '}',
    '',
    'func main() {',
    '    ui.Run(start, ui.Window{Title: "Cards", Width: 760, Height: 240})',
    '}',
    '',
    'func start(page *ui.Object) {',
    '    page.Set("color", "#F4F5F7")',
    '    ui.HGroup(ui.Props{Gap: 16, Padding: 24, Wrap: true}, func() {',
    '        for _, name := range []string{"Sales", "Stock", "Orders"} {',
    '            NewCard(name, "Updated today.", func(card *Card) {',
    '                card.SetTitle("Opened")',
    '            })',
    '        }',
    '    })',
    '}',
].join("\n");

// One line (or a few) for each component of the COMPONENTS section (the same order as the texts)
SAMPLES.componentLines = {
    SmartButton: 'ui.SmartButton(\n  ui.SmartButtonProps{\n    LabelText: "Save",\n    Tone: "cyan",\n    OnClick: save})',
    TextField: 'ui.TextField(\n  ui.TextFieldProps{\n    TitleText: "Name",\n    OnEnter: save})\nif email.Check() {\n  send()\n}',
    "UrlField, PhoneField": 'ui.UrlField(\n  ui.UrlFieldProps{\n    Value: "github.com"})\nui.CurrencyField(\n  ui.CurrencyFieldProps{\n    Value: 1234.5})',
    TextAreaField: 'ui.TextAreaField(\n  ui.TextAreaFieldProps{\n    TitleText: "Note",\n    ShowCount: true})',
    SearchInput: 'ui.SearchInput(\n  ui.SearchInputProps{\n    OnChange: find})',
    CheckBox: 'agree := ui.CheckBox(\n  ui.CheckBoxProps{\n    LabelText: "I agree"})\nagree.Checked()',
    Toggle: 'ui.Toggle(ui.ToggleProps{\n  Value: true,\n  OnChange: fn})',
    RadioButton: 'ui.RadioButton(\n  ui.RadioButtonProps{\n    Group: "size",\n    Value: "M"})\nui.RadioValue("size")',
    SelectBox: 'ui.SelectBox(\n  ui.SelectBoxProps{\n    Options: options,\n    Searchable: true})',
    "Slider, RangeSlider": 'ui.Slider(ui.SliderProps{\n  Value: 5,\n  MaxValue: 10})\nui.RangeSlider(\n  ui.RangeSliderProps{\n    LowValue: 2,\n    HighValue: 8})',
    Stepper: 'ui.Stepper(\n  ui.StepperProps{\n    Value: 2,\n    MaxValue: 10})',
    DatePicker: 'ui.DatePicker(\n  ui.DatePickerProps{\n    OnChange: fn})\nui.TimePicker(\n  ui.TimePickerProps{\n    Time: "19:30"})',
    FilePicker: 'ui.FilePicker(\n  ui.FilePickerProps{\n    Multiple: true,\n    OnChange: fn})\npath := ui.OpenFile()',
    "ProgressBar, StepBar": 'bar.SetValue(40)\nui.StepBar(\n  ui.StepBarProps{\n    Steps: steps,\n    CurrentStep: 1})',
    ItemList: 'ui.ItemList(\n  ui.ItemListProps{\n    Items: items})',
    SortableList: 'ui.SortableList(\n  ui.SortableListProps{\n    Items: items,\n    OnChange: fn})',
    SmartTable: 'ui.SmartTable(\n  ui.SmartTableProps{\n    TitleDataList: columns,\n    ItemDataList: people})',
    TreeView: 'ui.TreeView(\n  ui.TreeViewProps{\n    Items: nodes,\n    OnSelect: fn})',
    Tabs: 'tabs.Panel("general",\n  func() {\n    // its objects\n  })',
    Accordion: 'accordion.Section("more",\n  func() {\n    // its objects\n  })',
    GridGroup: 'ui.GridGroup(\n  ui.GridGroupProps{\n    Columns: 3},\n  func() {\n    // the cells\n  })',
    GraphBox: 'ui.GraphBox(\n  ui.GraphBoxProps{\n    Type: "bar",\n    Values: sales,\n    FormatValue:\n      "{:,.0f} TL"})',
    LogView: 'w := logView.Writer("info")\nlogView.Call("error",\n  "It failed")',
    ImageView: 'view.Call("setImage",\n  ui.Picture(img))',
    Canvas: 'canvas.Call("addCircle",\n  100, 80, 30,\n  ui.P{"draggable": 1})',
    "MarkdownView, JsonView": 'view.Call("appendText",\n  "The next words")\nui.JsonView(\n  ui.JsonViewProps{\n    Value: answer})',
    JsonEditor: 'ui.JsonEditor(\n  ui.JsonEditorProps{\n    Value: products,\n    OnChange: fn})',
    ChatView: 'chat.Call("appendText",\n  key, piece)',
    WebView: 'ui.WebView(\n  ui.WebViewProps{\n    HtmlText: html})',
    "MenuBar, ToolBar": 'ui.MenuBar(\n  ui.MenuBarProps{\n    Menus: menus})\nui.ToolBar(\n  ui.ToolBarProps{\n    Items: tools})',
    "StatusBar, Splitter": 'bar.Call("setMessageText",\n  "Saved")\nsplitter.First(func() {\n  // the left area\n})',
    CodeEditor: 'code := editor.Ask(\n  "getValue")',
    "ContextMenu, Tooltip": 'ui.Tooltip(\n  ui.TooltipProps{\n    Target: button,\n    HintText: "Saves"})',
    "Dialog, Toast": 'if ui.Confirm("Delete?") {\n  ui.ToastSuccess("Done")\n}',
};

SAMPLES.dark = [
    'func main() {',
    '    ui.Run(start, ui.Window{Theme: "dark"}) // "light", "dark", "auto"',
    '}',
    '',
    'func start(page *ui.Object) {',
    '    // the same objects',
    '}',
].join("\n");

SAMPLES.install = [
    'go mod init hello',
    'go get github.com/bug7a/ozdenui-go',
].join("\n");

SAMPLES.goGet = 'go get github.com/bug7a/ozdenui-go';

SAMPLES.hello = [
    'package main',
    '',
    'import ui "github.com/bug7a/ozdenui-go"',
    '',
    'func main() {',
    '    ui.Run(start, ui.Window{Title: "Hello"})',
    '}',
    '',
    'func start(page *ui.Object) {',
    '    ui.VGroup(ui.Props{}, func() {',
    '        ui.Label(ui.Props{Text: "Hello, window!", FontSize: 32})',
    '    })',
    '}',
].join("\n");

SAMPLES.run = [
    'go run .',
    '',
    '# one file to share',
    'go build -o hello',
].join("\n");

SAMPLES.agents = [
    '# the rules of OzdenUI for your assistant, into AGENTS.md',
    'go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents',
    '',
    '# a picture of the window, without a window',
    'OZDENUI_LOOK=app.png go run .',
].join("\n");
