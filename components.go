/* Bismillah */

/*

components.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
What every component of the kit (comp/) does on the Go side. The Props types and the functions of the components
themselves are in kit_gen.go, written by internal/gen from the files of the kit and from internal/gen/spec.json.

A component is driven by its entry in the table `kit`: its params, its methods, the state the page sends back, the
look it takes in a dark or a minimal window.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"fmt"
	"reflect"
)

type setter struct {
	Method string
	Silent bool // the method takes silent after the value (1: its callback is not called)
}

type componentSpec struct {
	Name          string
	Params        map[string]bool   // the params of the component (its defaults in the kit)
	Setters       map[string]setter // the params it has a method for: the way of the kit to change one later
	Methods       map[string]bool
	Defaults      P                 // params given unless the program gives them
	Aliases       map[string]string // a short name of a param in Set() and Get(): text -> labelText
	State         map[string]string // the values the page sends back with every event: {"checked": "checked"}
	Sync          []string          // callbacks always given to the page, so the state stays current
	StateDefaults P
	DataParams    []string // a true in them stays true (the others go as 1 / 0)
	FormatParams  []string // they take a Format() text
	WholeArgs     map[string][]int
	HasOwnPlace   bool   // it knows its place without one (a StatusBar: at the bottom of its box)
	Container     string // "self": objects go into it; a name: into that box inside it (the content of a ScrollArea)
	ContainerHint string // how objects go into it, when not directly: "tabs.Panel(key, func() {...})"
	FocusPart     string // the part that takes the keyboard (the input of a field)
	OnPage        bool   // it is made on the page, not in the box it was written in (a LoadingScreen covers the window)
	SizeInStyle   bool   // it reads its width and height from its style (SmartTable)
	DarkStyle     string // its package in a dark window
	HasMinimal    bool   // it has the minimal look (Window{Style: "minimal"})
	ThemeParam    bool   // it takes theme: "dark", not a package (GraphBox)

	known map[string]bool
}

var (
	defaultDataParams = []string{"value", "options", "items", "rows", "columns", "tabs", "selectedKeys", "group", "itemDataList", "titleDataList"}
	plainParams       = []string{"styleName", "style", "key"} // every component takes them besides its own
	// A light package the code names, and its dark twin: a dark window takes the twin.
	darkTwins = map[string]map[string]string{
		"ItemList":    {"gallery": "galleryDark"},
		"SmartButton": {"minimalDanger": "minimalDangerDark"},
	}
	darkTones = map[string]bool{"cyan": true, "orange": true, "yellow": true, "green": true, "blue": true, "ink": true, accentStyleName: true}
)

func (spec *componentSpec) isDataParam(key string) bool {
	params := spec.DataParams
	if params == nil {
		params = defaultDataParams
	}
	for _, name := range params {
		if name == key {
			return true
		}
	}
	return false
}

func (spec *componentSpec) isFormatParam(key string) bool {
	for _, name := range spec.FormatParams {
		if name == key {
			return true
		}
	}
	return false
}

// Every name the component takes: its params, the values of a box (not text and html: it fills its box itself).
func (spec *componentSpec) knownNames() map[string]bool {
	if len(spec.Params) == 0 {
		return nil // a component the kit copy does not describe: its params are not known either
	}
	if spec.known != nil {
		return spec.known
	}
	spec.known = map[string]bool{}
	for name := range objectProps {
		if name != "text" && name != "html" {
			spec.known[name] = true
		}
	}
	for name := range spec.Params {
		spec.known[name] = true
	}
	for _, name := range plainParams {
		spec.known[name] = true
	}
	for short := range spec.Aliases {
		spec.known[short] = true
	}
	return spec.known
}

// Component makes a component of the kit by its name, with its params by their names in the kit. The functions of
// kit_gen.go (ozdenui.CheckBox(ozdenui.CheckBoxProps{...})) do the same with a type that knows the params.
//
//	ozdenui.Component("Gauge", ozdenui.P{"value": 72, "width": 200})
func Component(typeName string, params P, content ...func()) *Object {
	values := P{}
	for key, value := range params {
		values[key] = value
	}
	read := given{values: values}
	if fill, has := values["fill"]; has {
		read.fill = isOn(fill)
		delete(values, "fill")
	}
	return makeComponent(typeName, read, content)
}

func makeTyped(typeName string, props any, content []func()) *Object {
	return makeComponent(typeName, readProps(props), content)
}

func makeComponent(typeName string, read given, content []func()) *Object {
	app := currentApp()
	app.check(typeName + "()")
	spec := kit[typeName]
	if spec == nil {
		spec = &componentSpec{Name: typeName}
	}
	values := read.values
	obj := app.newObject(typeName, kindComponent)
	obj.spec = spec

	parentID := app.containerID()
	if spec.OnPage {
		parentID = "page"
	}
	parent := app.objects[parentID]
	if !spec.HasOwnPlace && !spec.OnPage && isFreePlace(parent) {
		setDefault(values, "left", 0)
		setDefault(values, "top", 0)
	}
	if read.fill {
		takeFill(values, parent)
	}
	if read.unframed {
		values["style"] = unframedStyle(values["style"])
	}
	var onEnter any
	if spec.FocusPart != "" {
		onEnter = values["onEnter"] // WHY: not a param of the kit: the Enter key of the field
		delete(values, "onEnter")
	}
	obj.warnAboutUnknownNames(values)
	for short, long := range spec.Aliases {
		if value, has := values[short]; has && !spec.Params[short] {
			delete(values, short)
			values[long] = value
		}
	}
	isWindowLook := takeWindowLook(app, spec, values)
	for key, value := range spec.Defaults {
		if key == "styleName" && isWindowLook {
			continue // WHY: the look of a light window; the window chose the dark or the minimal one
		}
		setDefault(values, key, value)
	}
	if typeName == "FilePicker" && values["onOpenDialog"] == nil {
		values["onOpenDialog"] = openDialogOfPicker // WHY: the dialog of the page gives no path (files.go)
	}
	for _, name := range spec.Sync {
		if values[name] == nil {
			values[name] = keepState // WHY: the state of the component comes back with this callback
		}
	}
	for key, value := range values {
		if isCallbackName(key) {
			checkHandler(value, typeName+": "+key)
		}
	}
	if spec.SizeInStyle {
		sizeIntoStyle(values)
	}
	for key, value := range spec.StateDefaults {
		obj.values[key] = value
	}
	obj.keepValues(values)
	if typeName == "Tabs" && obj.values["value"] == nil {
		obj.values["value"] = firstKey(values["tabs"]) // the kit selects the first tab
	}
	app.objects[obj.id] = obj
	command := map[string]any{
		"op": "create", "id": obj.id, "type": typeName, "args": []any{}, "props": obj.pageValues(values),
		"parent": parentID, "src": sourceLine(),
	}
	if len(spec.State) > 0 {
		command["state"] = spec.State
	}
	app.send(command)
	if onEnter != nil {
		obj.setDomHandler("onEnter", onEnter)
	}
	for _, fill := range content {
		obj.Add(fill)
	}
	return obj
}

// The key of the first item of a list of the kit: its "key", else its "text", else the item itself.
func firstKey(list any) any {
	items := reflect.ValueOf(list)
	if !items.IsValid() || (items.Kind() != reflect.Slice && items.Kind() != reflect.Array) || items.Len() == 0 {
		return nil
	}
	first := items.Index(0).Interface()
	if item, ok := first.(P); ok {
		if key, has := item["key"]; has {
			return key
		}
		return item["text"]
	}
	return first
}

// The callback of a component that only keeps its state in the program (the page sends it with every call).
func keepState() {}

// The style of a component without its frame: no border and no rounded corners on its box. The keys of the code win.
func unframedStyle(style any) P {
	result := P{}
	if given, ok := style.(P); ok {
		for key, value := range given {
			result[key] = value
		}
	}
	box := P{"border": 0, "round": 0}
	if given, ok := result["box"].(P); ok {
		for key, value := range given {
			box[key] = value
		}
	}
	result["box"] = box
	return result
}

// WHY: a SmartTable reads its size from its style; a width given as a param would change only its box.
func sizeIntoStyle(values P) {
	style := P{}
	if given, ok := values["style"].(P); ok {
		for key, value := range given {
			style[key] = value
		}
	}
	for _, name := range []string{"width", "height"} {
		if value, has := values[name]; has {
			setDefault(style, name, value)
		}
	}
	if _, has := values["height"]; has {
		// WHY: with a height of its own the rows keep their height and their count follows the table; without it the
		//      16 rows of the kit share the height, each a few px high in a small table
		setDefault(values, "fitLineCount", 1)
	}
	if len(style) > 0 {
		values["style"] = style
	}
}

// A minimal window: the minimal look; a dark window: the dark look; unless the code chose one. It is true when the
// window chose the look.
func takeWindowLook(app *application, spec *componentSpec, values P) bool {
	isDark := app.theme == "dark"
	isMinimal := app.window.Style == "minimal"
	if spec.Name == "SmartButton" {
		takeWindowTone(app, values)
	}
	if name, has := values["styleName"].(string); has {
		if twin := darkTwins[spec.Name][name]; isDark && twin != "" {
			values["styleName"] = twin
		}
		return false
	}
	if spec.ThemeParam {
		if isMinimal && isDark {
			setDefault(values, "theme", "minimalDark")
		} else if isMinimal {
			setDefault(values, "theme", "minimal")
		} else if isDark {
			setDefault(values, "theme", "dark")
		}
		return false
	}
	if isMinimal && spec.HasMinimal {
		name := "minimal"
		if app.window.Accent != "" {
			name = accentStyleName
		}
		if isDark {
			name += "Dark"
		}
		values["styleName"] = name
		return true
	}
	if isDark && spec.DarkStyle != "" {
		values["styleName"] = spec.DarkStyle
		return true
	}
	return false
}

// The tone of a SmartButton in this window: "accent" is the main action in the color of the window, and a dark
// window takes the dark twin of a tone.
func takeWindowTone(app *application, values P) {
	tone, has := values["tone"].(string)
	_, hasStyle := values["styleName"]
	if !has || (hasStyle && tone != "accent") {
		return
	}
	if tone == "accent" {
		switch {
		case app.window.Accent != "":
			tone = accentStyleName
		case app.window.Style == "minimal":
			tone = "ink"
		default:
			tone = "cyan"
		}
	}
	if app.theme == "dark" && darkTones[tone] {
		tone += "Dark"
	}
	values["tone"] = tone
}

// A value of a component given later (Set): a callback replaces the one before it; a param with a method of the kit
// (setUnitText) is changed with it. It is true when the value is taken here.
func (obj *Object) setComponentValue(name string, value any) bool {
	spec := obj.spec
	if long, has := spec.Aliases[name]; has && !spec.Params[name] {
		name = long
	}
	if name == "onEnter" && spec.FocusPart != "" {
		obj.setDomHandler("onEnter", value)
		return true
	}
	if isCallbackName(name) {
		checkHandler(value, name)
		if value == nil || reflect.ValueOf(value).IsNil() {
			value = keepState // WHY: a nil would go to the page as null, and the component calls it
		}
		obj.app.send(map[string]any{"op": "set", "id": obj.id, "props": P{name: obj.pageValue(name, value)}, "src": sourceLine()})
		return true
	}
	if name == "styleName" && spec.Methods["setTheme"] {
		// The look of one object while the program runs (a dark switch of its own, without a new page)
		obj.keepValues(P{name: value})
		obj.sendCall("setTheme", []any{value}, false, "")
		return true
	}
	method, has := spec.Setters[name]
	if !has {
		if name == "styleName" || name == "style" {
			obj.app.warn("%s: %s is the look the component is made with; it does not change later", obj.typeName, name)
		}
		return false
	}
	obj.warnAboutUnknownNames(P{name: value})
	obj.keepValues(P{name: value})
	obj.uncheckOthersOfGroup()
	args := []any{obj.pageValue(name, value)}
	if method.Silent {
		args = append(args, 1)
	}
	command := map[string]any{"op": "call", "id": obj.id, "method": method.Method, "args": args, "src": sourceLine()}
	if len(spec.State) > 0 {
		command["refresh"] = 1 // WHY: the method may change the state too (a slider keeps its value inside its range)
	}
	obj.app.send(command)
	return true
}

// *** THE AREAS OF THE COMPONENTS THAT HOLD OBJECTS:

func (obj *Object) mustBe(typeName, method string) {
	if obj.typeName != typeName {
		panic(fmt.Sprintf("ozdenui: %s: %s() is of a %s, not of a %s", sourceLine(), method, typeName, obj.typeName))
	}
}

func (obj *Object) fillPart(part *Object, content []func()) *Object {
	for _, fill := range content {
		part.Add(fill)
	}
	return part
}

// Panel is the area of a tab of a Tabs, by the key of the tab; the objects made in content go into it.
//
//	tabs.Panel("general", func() { ozdenui.Label(ozdenui.Props{Text: "General"}) })
func (obj *Object) Panel(key string, content ...func()) *Object {
	obj.mustBe("Tabs", "Panel")
	if obj.parts["panel:"+key] == nil {
		// WHY once: it makes the wrapper of the panel, and a second call would clear it
		obj.sendCall("setPanelContent", []any{key, P{"__bpNoop": 1}}, false, "")
		obj.parts["panel:"+key] = obj.Part(P{"call": "getPanel", "args": []any{key}}, "wrapper")
	}
	return obj.fillPart(obj.parts["panel:"+key], content)
}

// Section is the content of a section of an Accordion, by its key.
func (obj *Object) Section(key string, content ...func()) *Object {
	obj.mustBe("Accordion", "Section")
	return obj.fillPart(obj.Part(P{"call": "getContent", "args": []any{key}}), content)
}

// First is the first area of a Splitter (the left or the top one).
func (obj *Object) First(content ...func()) *Object {
	obj.mustBe("Splitter", "First")
	return obj.fillPart(obj.Part("first"), content)
}

// Second is the second area of a Splitter.
func (obj *Object) Second(content ...func()) *Object {
	obj.mustBe("Splitter", "Second")
	return obj.fillPart(obj.Part("second"), content)
}

// Place gives an object a cell of its own in a GridGroup (from 0), or a span:
//
//	grid.Place(wide, ozdenui.P{"columnSpan": 2})      // also "column", "row", "rowSpan", "align"
//
// Inside the function of the grid, the grid itself is ozdenui.Current().
func (obj *Object) Place(item *Object, where P) *Object {
	obj.mustBe("GridGroup", "Place")
	obj.sendCall("place", []any{item, where}, false, "")
	return item
}

// Check shows the warning of a field (TextField and its kin) when its text is not right, and is true when it is
// right: a required field that is empty, an e-mail that is not one...
func (obj *Object) Check() bool {
	if obj.spec == nil || obj.spec.FocusPart == "" {
		panic(fmt.Sprintf("ozdenui: %s: Check() is of a field (TextField, EmailField...), not of a %s", sourceLine(), obj.typeName))
	}
	obj.sendCall("checkInputValue", nil, false, "")
	return numberOf(obj.Measure("status")) == 0
}

// RadioValue is the Value of the checked RadioButton of a group ("default" when no Group was given), or nil.
func RadioValue(group string) any {
	app := currentApp()
	for _, obj := range app.objects {
		if obj.typeName == "RadioButton" && obj.Str("group") == group && obj.Checked() {
			return obj.values["value"]
		}
	}
	return nil
}

// WHY: the page tells only of the radio the user pressed; the one that was checked before would stay checked in
// the program.
func (obj *Object) uncheckOthersOfGroup() {
	if obj.typeName != "RadioButton" || !isOn(obj.values["checked"]) {
		return
	}
	for _, other := range obj.app.objects {
		if other != obj && other.typeName == "RadioButton" && other.Str("group") == obj.Str("group") {
			other.values["checked"] = 0
		}
	}
}

// ShowHint shows a small bubble with a text at an object, for a while (milliseconds; 0: until hint.Remove()).
// side: "top" (default), "bottom", "left", "right". A bubble that comes under the mouse by itself is a Tooltip.
func ShowHint(target *Object, text string, milliseconds int, side ...string) *Object {
	app := currentApp()
	app.check("ShowHint()")
	look := P{"color": "#FFFFFF", "textColor": "rgba(0, 0, 0, 0.85)", "border": 1, "borderColor": "rgba(0, 0, 0, 0.12)", "shadow": "0 10px 30px rgba(0, 0, 0, 0.12)"}
	if app.theme == "dark" {
		// WHY a frame: the bubble of the kit has no line and no shadow, so it is not seen on a card of its own color
		look = P{"color": "#2A2A29", "textColor": darkText, "border": 1, "borderColor": "rgba(255, 255, 255, 0.14)", "shadow": "0 10px 30px rgba(0, 0, 0, 0.4)"}
	}
	place := "top"
	if len(side) > 0 {
		place = side[0]
	}
	bubble := app.newObject("Part", kindPart)
	app.objects[bubble.id] = bubble
	app.callStatic("Tooltip", "show", []any{app.toPage(target, "", nil), text, look, place}, false, bubble.id)
	if milliseconds > 0 {
		After(milliseconds, func() {
			if app.objects[bubble.id] != nil {
				bubble.Remove()
			}
		})
	}
	return bubble
}

// Body is the area of a SidePanel: its objects stand one under the other.
func (obj *Object) Body(content ...func()) *Object {
	obj.mustBe("SidePanel", "Body")
	return obj.fillPart(obj.Part(P{"call": "getBody"}), content)
}
