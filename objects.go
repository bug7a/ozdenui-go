/* Bismillah */

/*

objects.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The basic.js objects on the Go side: Box, Label, Button, Input, Icon, the groups, and what every object does.

An object here is a proxy: it keeps an id and the values it was given, and sends commands to the page, where the
real basic.js object lives. Reading a value it was given needs no question to the page (label.Text()).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// Props are the values of a basic object (Box, Label, Button, Input, Icon) and of a group (VGroup, HGroup). A value
// that is not given keeps the default of basic.js. A size is a number (px) or a text: "100%", "auto".
//
//	ozdenui.Label(ozdenui.Props{Text: "Hello", FontSize: 24, TextColor: "#1F2933"})
//
// Flags take true / false (or 1 / 0, as basic.js writes them).
type Props struct {
	Left   any // with Top: a place of its own in a Box or on the page (in a group, the group places the object)
	Top    any
	Right  any
	Bottom any
	Width  any // 120, "100%", "auto"
	Height any

	Text string // shown exactly as written (a "<b>" stays "<b>"); the text of an Input
	Html string // HTML: "<b>3</b> new"

	Color       string // the background
	TextColor   string
	FontSize    any
	FontFamily  string
	Bold        any
	Italic      any
	TextAlign   string // "left", "center", "right"
	LineHeight  any
	Ellipsis    any // a text that does not fit ends with "..."
	Border      any // px
	BorderColor string
	Round       any    // px: the corners
	BoxShadow   string // CSS
	Opacity     any    // 0 .. 1
	Visible     any
	Clickable   any
	Selectable  any // the text can be selected
	ClipContent any // what goes out of the box is cut
	ScrollX     any
	ScrollY     any
	Cursor      string
	Rotate      any
	ZIndex      any
	Position    string
	Title       string // the small text the system shows under the mouse
	Enabled     any
	TabIndex    any
	Role        string
	Aria        P
	Css         P // CSS of its own: ozdenui.P{"letterSpacing": "1px"}

	// An Input:
	Placeholder string
	InputType   string // "text", "password", "number"...
	MaxLength   any
	ReadOnly    any

	// An Icon:
	Picture  any    // a path, bytes or an image.Image (see Picture())
	ImageFit string // "contain", "cover"...
	Alt      string

	// A group (VGroup, HGroup):
	Flow    string
	Align   string // "center" is the default; "left top", "right center"...
	Gap     any
	Padding any // 20 or []int{x, y}
	Wrap    any
	Justify string
	Hug     any // as large as its objects
	Stretch any // the objects without a width / height fill the group across

	// Inside a group:
	Fill   bool // the object takes the room the others leave: the height in a VGroup, the width in an HGroup
	Grow   any
	Shrink any

	// Events: func(), func(obj *ozdenui.Object) or func(obj *ozdenui.Object, event ozdenui.Event).
	// Any other event of the page: obj.On("dblclick", fn)
	OnClick      any
	OnInput      any // an Input: at every change of its text
	OnEnter      any // an Input: the Enter key
	OnKeyDown    any
	OnMouseEnter any
	OnMouseLeave any
	OnFocus      any
	OnBlur       any

	// Extra: any other value of basic.js, by its own name: ozdenui.P{"textSize": 12}
	Extra P
}

// Event is what the page tells of an event of an object.
type Event struct {
	Name   string // "click", "keydown"...
	Key    string // a key event: "Enter", "a", "ArrowLeft"
	Code   string // "KeyA"
	Shift  bool
	Ctrl   bool
	Alt    bool
	Meta   bool
	X      float64 // a mouse event: its place in the window
	Y      float64
	Text   string // an Input: its text
	Detail any    // an event of a component with values of its own
}

type objectKind int

const (
	kindBasic objectKind = iota
	kindBox
	kindGroup
	kindPart
	kindPage
	kindComponent
)

// Object is every object of the window: a basic object, a group or a component.
type Object struct {
	app      *application
	id       string
	typeName string // the function of basic.js or of the kit that makes it: "Label", "CheckBox"
	kind     objectKind
	isRow    bool // an HGroup: Fill takes the width that is left
	spec     *componentSpec

	values         map[string]any // the values given, by their kit names
	handlers       map[string][]*handlerEntry
	listened       map[string]bool
	propHandlers   map[string]func()   // event -> the remover of the handler given as a value (OnClick)
	paramCallbacks map[string][]string // param -> the ids of the functions given to it (also in its items)
	parts          map[string]*Object
}

type handlerEntry struct {
	handler any
}

const darkText = "rgba(255, 255, 255, 0.86)"

var (
	positionNames = []string{"left", "top", "width", "height"}
	darkValues    = map[string]P{
		"Box":   {"color": "#232322", "textColor": darkText},
		"Label": {"textColor": darkText},
		"Input": {"color": "#2A2A29", "textColor": darkText, "borderColor": "rgba(255, 255, 255, 0.2)"},
	}
	// WHY the basis and the minimum: a flex item is as large as its content unless it may shrink (no scroll bar came)
	fillInColumn = struct{ props, css P }{P{"width": "100%", "height": "auto", "grow": 1}, P{"flexBasis": "0", "minHeight": "0"}}
	fillInRow    = struct{ props, css P }{P{"width": "auto", "height": "100%", "grow": 1}, P{"flexBasis": "0", "minWidth": "0"}}
	groupParams  = []string{"flow", "align", "gap", "padding", "wrap", "justify", "hug", "fit", "stretch"}
)

// *** THE BASIC OBJECTS:

// Box is a rectangle that holds other objects: the objects made in content go into it, each with a place of its own.
//
//	ozdenui.Box(ozdenui.Props{Left: 20, Top: 20, Width: 200, Height: 120, Color: "white", Round: 8}, func() {
//	    ozdenui.Label(ozdenui.Props{Left: 12, Top: 12, Text: "In the box"})
//	})
func Box(props Props, content ...func()) *Object {
	return makeBasic("Box", kindBox, props, content)
}

// VGroup holds objects one under the other (flexbox). A group is 100% x 100% by default; Hug: true wraps its
// objects, Fill: true takes the room the other objects of its group leave.
// Unlike in basic.js, a group does not cut what goes out of it (ClipContent: true cuts).
func VGroup(props Props, content ...func()) *Object {
	return makeBasic("VGroup", kindGroup, props, content)
}

// HGroup holds objects side by side (flexbox). See VGroup.
func HGroup(props Props, content ...func()) *Object {
	return makeBasic("HGroup", kindGroup, props, content)
}

// Label is a text. Text is shown as it is written; Html takes HTML.
func Label(props Props) *Object { return makeBasic("Label", kindBasic, props, nil) }

// Button is the plain button of basic.js. SmartButton is the modern one.
func Button(props Props) *Object { return makeBasic("Button", kindBasic, props, nil) }

// Input is a plain text box: obj.Text() is always the text on the page. TextField is the modern one.
func Input(props Props) *Object { return makeBasic("Input", kindBasic, props, nil) }

// Icon is a picture: ozdenui.Icon(ozdenui.Props{Width: 48, Height: 48, Picture: "photo.png"}).
// ImageView is the one with a zoom.
func Icon(props Props) *Object { return makeBasic("Icon", kindBasic, props, nil) }

// What a Props type gave, by the names of the kit.
type given struct {
	values     P
	fill       bool
	unframed   bool
	picture    any
	hasPicture bool
}

// Reads the fields of a Props type that were given (not the zero value). The name of a value is the ui tag of its
// field, or the name of the field as the kit writes it (LabelText -> labelText).
func readProps(props any) given {
	result := given{values: P{}}
	value := reflect.ValueOf(props)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.IsZero() {
			continue
		}
		name := valueType.Field(i).Name
		switch name {
		case "Fill":
			result.fill = field.Bool()
			continue
		case "Unframed":
			result.unframed = field.Bool()
			continue
		case "Extra":
			continue // after the fields, so the fields are read first
		}
		tag := valueType.Field(i).Tag.Get("ui")
		if tag == "" {
			tag = kitName(name)
		}
		if tag == "picture" {
			result.picture, result.hasPicture = field.Interface(), true
			continue
		}
		result.values[tag] = field.Interface()
	}
	if extra := value.FieldByName("Extra"); extra.IsValid() && !extra.IsNil() {
		for key, item := range extra.Interface().(P) {
			result.values[key] = item
		}
	}
	return result
}

func makeBasic(typeName string, kind objectKind, props Props, content []func()) *Object {
	app := currentApp()
	app.check(typeName + "()")
	read := readProps(props)
	values := read.values
	obj := app.newObject(typeName, kind)
	obj.isRow = (typeName == "HGroup")
	if typeName == "Input" {
		obj.listened["input"] = true // the page always tells the text of an Input
	}

	parent := app.objects[app.containerID()]
	if isFreePlace(parent) {
		// WHY: basic.js puts an object without a place at -1000, -1000 (out of sight until it is placed). Outside a
		//      group it gets the top left corner: the commands go together, so nothing is seen early.
		setDefault(values, "left", 0)
		setDefault(values, "top", 0)
	}
	if read.fill {
		takeFill(values, parent)
	}
	handlers := takeDomHandlers(values)
	obj.warnAboutUnknownNames(values)
	if kind == kindGroup && !isOn(values["scrollY"]) && !isOn(values["scrollX"]) {
		// WHY: a group of basic.js cuts a button that grows, a focus ring or a shadow at its edge, so nearly every
		//      group of a program would say ClipContent: false.
		setDefault(values, "clipContent", 0)
	}
	if app.theme == "dark" && kind != kindGroup {
		for key, value := range darkValues[typeName] {
			setDefault(values, key, value)
		}
	}
	obj.keepValues(values)
	app.objects[obj.id] = obj
	app.send(map[string]any{
		"op": "create", "id": obj.id, "type": typeName, "args": []any{}, "props": obj.pageValues(values),
		"parent": app.containerID(), "src": sourceLine(),
	})
	if typeName == "Input" {
		app.send(map[string]any{"op": "listen", "id": obj.id, "event": "input"})
	}
	for _, name := range sortedKeys(handlers) {
		obj.setDomHandler(name, handlers[name])
	}
	if read.hasPicture {
		obj.setPicture(read.picture)
	}
	for _, fill := range content {
		obj.Add(fill)
	}
	return obj
}

func (app *application) newObject(typeName string, kind objectKind) *Object {
	return &Object{
		app: app, id: app.newID(), typeName: typeName, kind: kind,
		values: P{}, handlers: map[string][]*handlerEntry{}, listened: map[string]bool{},
		propHandlers: map[string]func(){}, paramCallbacks: map[string][]string{}, parts: map[string]*Object{},
	}
}

func newPage(app *application) *Object {
	page := app.newObject("page", kindPage)
	page.id = "page"
	page.values["width"], page.values["height"] = app.pageWidth, app.pageHeight
	page.listened["resize"] = true
	app.objects["page"] = page
	// WHY always: the width and the height of the page come with an event of the page; without a handler of the
	//      program they would keep the size of the start
	app.send(map[string]any{"op": "listen", "id": "page", "event": "resize"})
	return page
}

func setDefault(values P, key string, value any) {
	if _, has := values[key]; !has {
		values[key] = value
	}
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// The page or a Box: its objects have a place of their own (in a group, the group places them).
func isFreePlace(container *Object) bool {
	return container != nil && (container.kind == kindPage || container.kind == kindBox)
}

func takeFill(values P, parent *Object) {
	look := fillInColumn
	if parent != nil && parent.isRow {
		look = fillInRow
	}
	for key, value := range look.props {
		setDefault(values, key, value)
	}
	css := P{}
	for key, value := range look.css {
		css[key] = value
	}
	if own, ok := values["css"].(P); ok {
		for key, value := range own {
			css[key] = value
		}
	}
	values["css"] = css
}

// onClick: fn of a basic object is an event of the page, not a value: {"click": fn}, out of values.
func takeDomHandlers(values P) map[string]any {
	handlers := map[string]any{}
	for name, value := range values {
		if !isCallbackName(name) {
			continue
		}
		checkHandler(value, name)
		handlers[name] = value
		delete(values, name)
	}
	return handlers
}

// *** VALUES:

func (obj *Object) isPlainText() bool {
	return obj.typeName == "Label" || obj.typeName == "Button"
}

// WHY plainText: the text of basic.js is HTML, so a typed "<b>" would be bold and a "<img onerror=...>" would run.
// Here Text is the text as it is written; Html is the HTML (a Label and a Button of basic.js take it as their text).
func (obj *Object) pageKey(key string) string {
	if obj.isPlainText() && key == "text" {
		return "plainText"
	}
	if obj.isPlainText() && key == "html" {
		return "text"
	}
	return key
}

// The values as they go to the page. The functions given to a param replace the ones given to it before.
func (obj *Object) pageValues(values P) P {
	result := P{}
	for key, value := range values {
		result[obj.pageKey(key)] = obj.pageValue(key, value)
	}
	return result
}

func (obj *Object) pageValue(key string, value any) any {
	if !obj.isDataParam(key) {
		value = flagAsNumber(value)
	}
	if text, ok := value.(string); ok && isFileParam(key) {
		value = pictureOfPath(text)
	}
	if obj.spec != nil && obj.spec.isFormatParam(key) {
		if text, ok := value.(string); ok {
			value = Format(text)
		}
	}
	ids := []string{}
	converted := obj.app.toPage(value, obj.id, &ids)
	obj.replaceCallbacks(key, ids)
	if marker, ok := converted.(map[string]any); ok && marker["__bpCallback"] != nil && obj.spec != nil {
		if places, has := obj.spec.WholeArgs[key]; has {
			marker["wholeArgs"] = places
		}
	}
	return converted
}

func (obj *Object) isDataParam(key string) bool {
	if obj.spec == nil {
		return key == "value"
	}
	return obj.spec.isDataParam(key)
}

// WHY: every Set("onChange", fn) would keep a new function for ever, and so would every new list of items whose
// items have an onClick (with what each function holds).
func (obj *Object) replaceCallbacks(slot string, ids []string) {
	for _, id := range obj.paramCallbacks[slot] {
		delete(obj.app.callbacks, id)
	}
	delete(obj.paramCallbacks, slot)
	if len(ids) > 0 {
		obj.paramCallbacks[slot] = ids
	}
}

// The values given (not the functions), kept to be read back without a question to the page.
func (obj *Object) keepValues(values P) {
	for key, value := range values {
		if value != nil && reflect.ValueOf(value).Kind() == reflect.Func {
			continue
		}
		if text, ok := value.(string); ok && text == Empty {
			value = ""
		}
		if obj.isPlainText() && key == "text" {
			delete(obj.values, "html") // WHY: after one of them, the other one is not what the object shows
		}
		if obj.isPlainText() && key == "html" {
			delete(obj.values, "text")
		}
		obj.values[key] = value
		if obj.spec == nil {
			continue
		}
		// A value with two names (the "text" of a field is its "value") is read by both.
		for short, long := range obj.spec.Aliases {
			if key == long {
				obj.values[short] = value
			} else if key == short {
				obj.values[long] = value
			}
		}
	}
}

func (obj *Object) warnAboutUnknownNames(values P) {
	known := knownNames(obj)
	if known == nil {
		return
	}
	for name := range values {
		if known[name] {
			continue
		}
		hint := ""
		if guess := closestName(name, known); guess != "" {
			hint = fmt.Sprintf(` (did you mean "%s"?)`, guess)
		}
		obj.app.warn(`%s: "%s" is not a value of %s%s`, obj.typeName, name, obj.typeName, hint)
	}
}

// Set changes values by their names in the kit: label.Set("text", "Hi", "fontSize", 20). A name is written as
// basic.js writes it (camelCase). On a component, a param it has a method for (setUnitText) is changed with that
// method: the component draws it again, and its callback is not called.
func (obj *Object) Set(namesAndValues ...any) *Object {
	obj.app.check(obj.typeName + ".Set()")
	if len(namesAndValues)%2 != 0 {
		panic(fmt.Sprintf("ozdenui: %s: Set() takes pairs: a name, its value, a name, its value...", sourceLine()))
	}
	values := P{}
	for i := 0; i < len(namesAndValues); i += 2 {
		name, ok := namesAndValues[i].(string)
		if !ok {
			panic(fmt.Sprintf("ozdenui: %s: Set(): a name is a text, not %T", sourceLine(), namesAndValues[i]))
		}
		value := namesAndValues[i+1]
		if obj.takeSpecialValue(name, value) {
			continue
		}
		values[name] = value
	}
	if len(values) == 0 {
		return obj
	}
	obj.warnAboutUnknownNames(values)
	obj.keepValues(values)
	obj.app.send(map[string]any{"op": "set", "id": obj.id, "props": obj.pageValues(values), "src": sourceLine()})
	return obj
}

// A value that is not a plain value of the page: an event, the picture of an Icon, a param with a method.
func (obj *Object) takeSpecialValue(name string, value any) bool {
	if obj.kind == kindComponent {
		return obj.setComponentValue(name, value)
	}
	if isCallbackName(name) {
		checkHandler(value, name)
		obj.setDomHandler(name, value)
		return true
	}
	if name == "picture" && obj.typeName == "Icon" {
		obj.setPicture(value)
		return true
	}
	if name == "fill" {
		panic(fmt.Sprintf("ozdenui: %s: Fill is given when the object is made", sourceLine()))
	}
	return false
}

// WHY a method and not a value: an Icon loads its file with load(); a value named picture would do nothing.
func (obj *Object) setPicture(picture any) {
	obj.values["picture"] = picture
	obj.Call("load", Picture(picture))
}

// Get is a value the object was given, or the state the page told (the checked of a CheckBox), by its name in the
// kit. A value that was never given is nil: obj.Measure("width") asks the page.
func (obj *Object) Get(name string) any {
	if value, ok := obj.lookup(name); ok {
		return value
	}
	obj.app.warn(`%s.Get("%s") was not given: obj.Measure("%s") reads it from the page`, obj.typeName, name, name)
	return nil
}

func (obj *Object) lookup(name string) (any, bool) {
	value, ok := obj.values[name]
	return value, ok
}

// Text is the text of the object: of a Label, of an Input or a field as it is on the page now, of a button.
func (obj *Object) Text() string {
	if value, ok := obj.lookup("text"); ok {
		return textOf(value)
	}
	if value, ok := obj.lookup("html"); ok {
		return textOf(value)
	}
	if value, ok := obj.lookup("value"); ok {
		return textOf(value)
	}
	return ""
}

// SetText changes the text: label.SetText("Saved").
func (obj *Object) SetText(text string) *Object {
	if text == "" {
		text = Empty
	}
	return obj.Set("text", text)
}

// Value is the value of a component: of a SelectBox, a Slider, a Toggle, a DatePicker... (nil: none).
func (obj *Object) Value() any {
	value, _ := obj.lookup("value")
	return value
}

// SetValue changes the value of a component; its OnChange is not called (only the user's change calls it).
func (obj *Object) SetValue(value any) *Object { return obj.Set("value", value) }

// Checked is true for a checked CheckBox or RadioButton, and for a Toggle that is on.
func (obj *Object) Checked() bool {
	if value, ok := obj.lookup("checked"); ok {
		return isOn(value)
	}
	return isOn(obj.Value())
}

// Str is a value as a text ("" when it was not given): field.Str("titleText").
func (obj *Object) Str(name string) string {
	value, _ := obj.lookup(name)
	return textOf(value)
}

// Num is a value as a number (0 when it was not given): slider.Num("value").
func (obj *Object) Num(name string) float64 {
	value, _ := obj.lookup(name)
	return numberOf(value)
}

// Int is a value as a whole number.
func (obj *Object) Int(name string) int { return int(obj.Num(name)) }

// Bool is a flag: true for 1, true and a text that is not empty.
func (obj *Object) Bool(name string) bool {
	value, _ := obj.lookup(name)
	return isOn(value)
}

// Show and Hide change the visible value: a hidden object keeps its place in the program, not on the page.
func (obj *Object) Show() *Object { return obj.Set("visible", true) }

// Hide: see Show.
func (obj *Object) Hide() *Object { return obj.Set("visible", false) }

// Measure is a value as basic.js reads it on the page now (a size of "auto", the scrollTop of a box...). It is a
// question to the page: keep it out of loops.
func (obj *Object) Measure(name string) any {
	obj.app.check("Measure()")
	obj.app.flush()
	value, err := obj.app.server.request(map[string]any{"kind": "read", "id": obj.id, "key": obj.pageKey(name)})
	if err != nil {
		obj.app.warn(`%s.Measure("%s"): %v`, obj.typeName, name, err)
		return nil
	}
	return obj.app.goValue(value)
}

// Call calls a method of the object on the page: box.Call("setMotion", "left 0.3s"). Its result is not waited for.
func (obj *Object) Call(method string, args ...any) *Object {
	obj.sendCall(method, args, false, "")
	return obj
}

func (obj *Object) sendCall(method string, args []any, refresh bool, slot string) {
	obj.app.check(obj.typeName + "." + method + "()")
	if obj.spec != nil && len(obj.spec.Methods) > 0 && !obj.spec.Methods[method] && !objectMethods[method] {
		hint := ""
		if guess := closestName(method, obj.spec.Methods); guess != "" {
			hint = fmt.Sprintf(` (did you mean "%s"?)`, guess)
		}
		obj.app.warn(`%s: "%s" is not a method of %s%s`, obj.typeName, method, obj.typeName, hint)
	}
	ids := []string{}
	pageArgs := make([]any, len(args))
	for i, value := range args {
		pageArgs[i] = obj.app.toPage(value, obj.id, &ids)
	}
	if slot != "" {
		obj.replaceCallbacks(slot, ids)
	}
	command := map[string]any{"op": "call", "id": obj.id, "method": method, "args": pageArgs, "src": sourceLine()}
	if refresh {
		command["refresh"] = 1
	}
	obj.app.send(command)
}

// Ask calls a method of the object on the page and returns its result: table.Ask("getSelectedData"). It is a
// question to the page: keep it out of loops.
func (obj *Object) Ask(method string, args ...any) any {
	obj.app.check("Ask()")
	pageArgs := make([]any, len(args))
	for i, value := range args {
		pageArgs[i] = obj.app.toPage(value, obj.id, nil)
	}
	obj.app.flush()
	value, err := obj.app.server.request(map[string]any{"kind": "ask", "id": obj.id, "method": method, "args": pageArgs})
	if err != nil {
		obj.app.warn(`%s.Ask("%s"): %v`, obj.typeName, method, err)
		return nil
	}
	return obj.app.goValue(value)
}

// Writer writes lines into a LogView, from any goroutine: every Write is a line (or more) of the level "info",
// "warn", "error", "debug" or "success".
//
//	log.SetOutput(logView.Writer("info"))
//	fmt.Fprintln(logView.Writer("error"), "the file is missing")
func (obj *Object) Writer(level string) io.Writer {
	obj.mustBe("LogView", "Writer")
	return lineWriter{obj: obj, level: level}
}

type lineWriter struct {
	obj   *Object
	level string
}

func (w lineWriter) Write(text []byte) (int, error) {
	lines := strings.Split(strings.TrimRight(string(text), "\n"), "\n")
	Do(func() {
		for _, line := range lines {
			w.obj.Call(w.level, line)
		}
	})
	return len(text), nil
}

// JS is the object in JavaScript, for Evaluate(): ozdenui.Evaluate(label.JS() + ".elem.offsetWidth")
func (obj *Object) JS() string { return "__bp.objectOf('" + obj.id + "')" }

// TypeName is the name of the object in the kit: "Label", "SmartButton".
func (obj *Object) TypeName() string { return obj.typeName }

func (obj *Object) String() string { return "<" + obj.typeName + " " + obj.id + ">" }

// Focus puts the keyboard into the object (an Input, a field).
func (obj *Object) Focus() *Object {
	if obj.spec != nil && obj.spec.FocusPart != "" {
		obj.Part(obj.spec.FocusPart).Call("focus")
		return obj
	}
	return obj.Call("focus")
}

// *** EVENTS (of the page: click, input, keydown, mouseenter...):

// On listens to an event of the object: handler is func(), func(obj *Object) or func(obj *Object, event Event).
// It returns a function that removes the handler. The page listens to its window:
// page.On("keydown", fn), page.On("resize", fn).
func (obj *Object) On(event string, handler any) (remove func()) {
	obj.app.check(obj.typeName + ".On()")
	checkHandler(handler, `On("`+event+`")`)
	entry := &handlerEntry{handler: handler}
	obj.handlers[event] = append(obj.handlers[event], entry)
	if !obj.listened[event] {
		obj.listened[event] = true
		obj.app.send(map[string]any{"op": "listen", "id": obj.id, "event": event, "src": sourceLine()})
	}
	return func() {
		entries := obj.handlers[event]
		for i, candidate := range entries {
			if candidate == entry {
				obj.handlers[event] = append(append([]*handlerEntry{}, entries[:i]...), entries[i+1:]...)
				return
			}
		}
	}
}

// The handler given as a value (OnClick: fn, Set("onClick", fn)): a new one replaces it, nil takes it away; the
// handlers of On() stay.
func (obj *Object) setDomHandler(name string, handler any) {
	if remove := obj.propHandlers[name]; remove != nil {
		remove()
		delete(obj.propHandlers, name)
	}
	if handler == nil || reflect.ValueOf(handler).IsNil() {
		return
	}
	if name == "onEnter" {
		// The Enter key of an Input
		obj.propHandlers[name] = obj.On("keydown", func(self *Object, event Event) {
			if event.Key == "Enter" {
				callHandler(handler, []any{self, event})
			}
		})
		return
	}
	obj.propHandlers[name] = obj.On(eventNameOf(name), handler)
}

func (obj *Object) receive(name string, data map[string]any) {
	if text, ok := data["text"].(string); ok {
		obj.values["text"] = text
	}
	if state, ok := data["state"].(map[string]any); ok {
		obj.takeState(state)
	}
	entries := obj.handlers[name]
	if len(entries) == 0 {
		return
	}
	event := Event{Name: name}
	event.Key, _ = data["key"].(string)
	event.Code, _ = data["code"].(string)
	event.Shift, _ = data["shift"].(bool)
	event.Ctrl, _ = data["ctrl"].(bool)
	event.Alt, _ = data["alt"].(bool)
	event.Meta, _ = data["meta"].(bool)
	event.X, _ = data["x"].(float64)
	event.Y, _ = data["y"].(float64)
	event.Text, _ = data["text"].(string)
	event.Detail = obj.app.goValueOf(data["detail"])
	for _, entry := range append([]*handlerEntry{}, entries...) {
		callHandler(entry.handler, []any{obj, event})
	}
}

// The state the page sent with an event: the values the program reads without a question (checked, value...).
func (obj *Object) takeState(state map[string]any) {
	for key, value := range state {
		value = obj.app.goValueOf(value)
		obj.values[key] = value
		if obj.spec != nil && obj.spec.Aliases[key] != "" {
			obj.values[obj.spec.Aliases[key]] = value // the "text" of a field is its "value"
		}
	}
	obj.uncheckOthersOfGroup()
}

// *** A BOX HOLDS OBJECTS:

// Add makes the objects of content inside this box or group, also later (in an event, a timer):
//
//	list.Add(func() {
//	    ozdenui.Label(ozdenui.Props{Text: "One more"})
//	})
//
// Only a Box, a group and the areas of the components hold objects (a ScrollArea, the panel of a Tabs...).
func (obj *Object) Add(content func()) *Object {
	obj.app.check(obj.typeName + ".Add()")
	target := obj.container()
	obj.app.containers = append(obj.app.containers, target.id)
	defer func() { obj.app.containers = obj.app.containers[:len(obj.app.containers)-1] }()
	content()
	return obj
}

// Current is the box that takes the objects made now: inside the function of a group, the group itself (its
// variable is not set yet there).
//
//	ozdenui.GridGroup(ozdenui.GridGroupProps{Columns: 4}, func() {
//	    wide := ozdenui.SmartButton(ozdenui.SmartButtonProps{LabelText: "0"})
//	    ozdenui.Current().Place(wide, ozdenui.P{"columnSpan": 2})
//	})
func Current() *Object {
	app := currentApp()
	app.check("Current()")
	return app.objects[app.containerID()]
}

// The box that takes the objects: the object itself, or a box inside a component (the content of a ScrollArea).
func (obj *Object) container() *Object {
	switch obj.kind {
	case kindBox, kindGroup, kindPart, kindPage:
		return obj
	case kindComponent:
		if obj.spec.Container == "self" {
			return obj
		}
		if obj.spec.Container != "" {
			return obj.Part(obj.spec.Container)
		}
	}
	hint := ": only a Box, a group or an area of a component can. Make the object next to it, in the same box."
	if obj.spec != nil && obj.spec.ContainerHint != "" {
		hint = ": " + obj.spec.ContainerHint
	}
	panic(fmt.Sprintf("ozdenui: %s: a %s can not hold other objects%s", sourceLine(), obj.typeName, hint))
}

// Part is a box inside a component, to put objects into or to call: field.Part("input").Call("select").
// A step of the path is the name of a property, a place in a list (0) or a call of a method:
// ozdenui.P{"call": "getPanel", "args": []any{"details"}}. The same path gives the same object.
func (obj *Object) Part(path ...any) *Object {
	obj.app.check(obj.typeName + ".Part()")
	key := fmt.Sprint(path...)
	if part := obj.parts[key]; part != nil && obj.app.objects[part.id] != nil {
		return part
	}
	part := obj.app.newObject("Part", kindPart)
	obj.app.objects[part.id] = part
	obj.parts[key] = part
	obj.app.send(map[string]any{"op": "part", "id": part.id, "of": obj.id, "path": obj.app.toPage(path, obj.id, nil), "src": sourceLine()})
	return part
}

// Remove removes the object and every object inside it, for good. To hide it for a while: obj.Hide().
func (obj *Object) Remove() {
	obj.app.check(obj.typeName + ".Remove()")
	if obj.kind == kindPage {
		panic("ozdenui: the page can not be removed: close the window with ozdenui.Close()")
	}
	obj.app.send(map[string]any{"op": "remove", "id": obj.id, "src": sourceLine()})
}

// Width and Height of the page are current (the page tells them at every resize); of another object they are the
// values it was given, as numbers (a size of "auto" or "100%": Measure("width")).
func (obj *Object) Width() float64 { return obj.Num("width") }

// Height: see Width.
func (obj *Object) Height() float64 { return obj.Num("height") }

func knownNames(obj *Object) map[string]bool {
	switch obj.kind {
	case kindComponent:
		return obj.spec.knownNames()
	case kindPart:
		return nil // what it is, is known only on the page
	case kindGroup:
		return groupNames
	}
	return objectProps
}

var groupNames = func() map[string]bool {
	names := map[string]bool{}
	for name := range objectProps {
		names[name] = true
	}
	for _, name := range groupParams {
		names[name] = true
	}
	return names
}()

// A method every basic.js object has (not in the list of a component).
var objectMethods = map[string]bool{
	"focus": true, "blur": true, "setMotion": true, "center": true, "centerBy": true, "aline": true, "props": true,
	"load": true, "select": true, "remove": true, "setFocusRing": true, "withMotion": true, "dontMotion": true,
	"add": true, "scrollTo": true, "click": true,
}
