/* Bismillah */

/*

values.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The values of the program as they go to the page, the functions of the program as the page calls them, and the
line of the program a message names.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// P is a short name for a map of values: a style, an item of a list, a row of a table.
//
//	ozdenui.P{"key": "tr", "text": "Türkiye"}
type P = map[string]any

// Empty is an empty text that is given on purpose. WHY: a text field of a Props type that is "" was not given, so
// the component keeps its own default; LabelText: ozdenui.Empty is a button without a text.
const Empty = "\x00empty"

const modulePath = "github.com/bug7a/ozdenui-go"

var objectType = reflect.TypeOf((*Object)(nil))

// *** THE LINE OF THE PROGRAM: "main.go:12", the first line of the stack that is not of this module.

func sourceLine() string {
	callers := make([]uintptr, 32)
	count := runtime.Callers(2, callers)
	frames := runtime.CallersFrames(callers[:count])
	for {
		frame, more := frames.Next()
		isOurs := strings.HasPrefix(frame.Function, modulePath+".") || strings.HasPrefix(frame.Function, modulePath+"/internal")
		isGo := strings.HasPrefix(frame.Function, "runtime.") || strings.HasPrefix(frame.Function, "reflect.")
		// WHY the tests: they are of this module, and their lines are the lines of a program
		if frame.File != "" && (!isOurs || strings.HasSuffix(frame.File, "_test.go")) && !isGo {
			return fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		}
		if !more {
			return ""
		}
	}
}

// Tells once in the terminal, at the line of the program.
func (app *application) warn(format string, args ...any) {
	text := sourceLine() + ": " + fmt.Sprintf(format, args...)
	if app.warned[text] {
		return
	}
	app.warned[text] = true
	fmt.Fprintln(os.Stderr, "ozdenui: "+text)
}

// *** NAMES: LabelText (a field of Go) <-> labelText (the kit)

func kitName(fieldName string) string {
	runes := []rune(fieldName)
	// "URL" -> "url", "ID" -> "id", "LabelText" -> "labelText"
	for i := range runes {
		isLastOfUpperRun := i > 0 && i+1 < len(runes) && unicode.IsLower(runes[i+1])
		if !unicode.IsUpper(runes[i]) || isLastOfUpperRun {
			break
		}
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// "onClick" -> "click", "onKeyDown" -> "keydown"
func eventNameOf(name string) string {
	return strings.ToLower(strings.TrimPrefix(name, "on"))
}

func isCallbackName(name string) bool {
	return len(name) > 2 && strings.HasPrefix(name, "on") && unicode.IsUpper(rune(name[2]))
}

// style, look, buttonStyle, searchParams...: a param of looks, not of data.
func isStyleName(name string) bool {
	return name == "style" || name == "look" || strings.HasSuffix(name, "Style") || strings.HasSuffix(name, "Params")
}

func closestName(name string, known map[string]bool) string {
	best, bestDistance := "", 3
	lower := strings.ToLower(name)
	for candidate := range known {
		distance := editDistance(lower, strings.ToLower(candidate))
		if distance < bestDistance || (distance == bestDistance && best != "" && candidate < best) {
			best, bestDistance = candidate, distance
		}
	}
	return best
}

func editDistance(a, b string) int {
	previous := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		row := make([]int, len(b)+1)
		row[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			row[j] = min(previous[j]+1, row[j-1]+1, previous[j-1]+cost)
		}
		previous = row
	}
	return previous[len(b)]
}

// *** VALUES TO THE PAGE:

// A value as it goes to the page (JSON). A function becomes a callback of the owner ({ __bpCallback: id }), an object
// its id; a struct goes as its JSON (its json tags name the keys). ids takes the ids of the callbacks made.
func (app *application) toPage(value any, ownerID string, ids *[]string) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		if typed == Empty {
			return ""
		}
		return typed
	case bool, int, float64, int64, float32, int32, uint, json.Number:
		return typed
	case *Object:
		if typed == nil {
			return nil
		}
		return map[string]any{"__bpObject": typed.id}
	case numberFormat:
		return map[string]any{"__bpFormat": typed}
	case time.Time:
		return typed.Format(time.RFC3339)
	case json.RawMessage:
		return typed
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Func:
		if reflected.IsNil() {
			return nil
		}
		id := app.registerCallback(value, ownerID)
		if ids != nil {
			*ids = append(*ids, id)
		}
		return map[string]any{"__bpCallback": id}
	case reflect.Map:
		if reflected.Type().Key().Kind() != reflect.String {
			break
		}
		result := make(map[string]any, reflected.Len())
		iterator := reflected.MapRange()
		for iterator.Next() {
			result[iterator.Key().String()] = app.toPage(iterator.Value().Interface(), ownerID, ids)
		}
		return result
	case reflect.Slice, reflect.Array:
		if reflected.Kind() == reflect.Slice && reflected.Type().Elem().Kind() == reflect.Uint8 {
			break // bytes: as JSON writes them
		}
		result := make([]any, reflected.Len())
		for i := range result {
			result[i] = app.toPage(reflected.Index(i).Interface(), ownerID, ids)
		}
		return result
	case reflect.Pointer, reflect.Interface:
		if reflected.IsNil() {
			return nil
		}
		return app.toPage(reflected.Elem().Interface(), ownerID, ids)
	}
	// A struct of the program (a row, an item) or another kind of number: its JSON, as it is written.
	// WHY not a map: the keys of a map go in the order of the alphabet; a struct shown in a JsonView or a table
	//     without columns keeps the order of its fields.
	text, err := json.Marshal(value)
	if err != nil {
		app.warn("a value of the type %T can not go to the page: %v", value, err)
		return nil
	}
	return json.RawMessage(text)
}

// true / false -> 1 / 0, as basic.js takes its flags (a Toggle compares value === 1).
func flagAsNumber(value any) any {
	if flag, ok := value.(bool); ok {
		if flag {
			return 1
		}
		return 0
	}
	return value
}

// *** FUNCTIONS OF THE PROGRAM, CALLED BY THE PAGE:

// Calls a function of the program with the values of the page. The function takes what it needs: func(),
// func(obj *Object), func(obj *Object, value string)... A value becomes the type of its place (a number of the page is
// an int, a float64...; a list or a map fills a slice, a map or a struct of the program).
func callHandler(handler any, args []any) {
	function := reflect.ValueOf(handler)
	if function.Kind() != reflect.Func || function.IsNil() {
		return
	}
	functionType := function.Type()
	count := functionType.NumIn()
	inputs := make([]reflect.Value, 0, count)
	for i := 0; i < count; i++ {
		place := functionType.In(i)
		if functionType.IsVariadic() && i == count-1 {
			for _, value := range args[min(i, len(args)):] {
				inputs = append(inputs, valueFor(value, place.Elem()))
			}
			break
		}
		var value any
		if i < len(args) {
			value = args[i]
		}
		inputs = append(inputs, valueFor(value, place))
	}
	function.Call(inputs)
}

func valueFor(value any, place reflect.Type) reflect.Value {
	if value == nil {
		return reflect.Zero(place)
	}
	given := reflect.ValueOf(value)
	if given.Type().AssignableTo(place) {
		return given
	}
	if place == objectType {
		return reflect.Zero(place) // a part of a component (a row) is not an *Object: take it as `any`
	}
	isNumber := given.Kind() == reflect.Float64 && place.Kind() >= reflect.Int && place.Kind() <= reflect.Float64 && place.Kind() != reflect.Uintptr
	if isNumber {
		return given.Convert(place)
	}
	if place.Kind() == reflect.String {
		return reflect.ValueOf(textOf(value)).Convert(place)
	}
	if place.Kind() == reflect.Bool {
		return reflect.ValueOf(isOn(value)).Convert(place)
	}
	filled := reflect.New(place).Elem()
	switch place.Kind() {
	case reflect.Struct:
		// WHY not through JSON: a flag of the kit is 1 / 0, and JSON does not put a number into a bool field (the whole
		//     struct stayed empty without a word)
		fields, isMap := value.(map[string]any)
		if !isMap {
			break
		}
		for i := 0; i < place.NumField(); i++ {
			field := place.Field(i)
			if !field.IsExported() {
				continue
			}
			name := strings.Split(field.Tag.Get("json"), ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = field.Name
			}
			item, has := fields[name]
			if !has {
				item, has = fields[kitName(field.Name)]
			}
			if has {
				filled.Field(i).Set(valueFor(item, field.Type))
			}
		}
	case reflect.Slice:
		items, isList := value.([]any)
		if !isList {
			break
		}
		filled = reflect.MakeSlice(place, len(items), len(items))
		for i, item := range items {
			filled.Index(i).Set(valueFor(item, place.Elem()))
		}
	case reflect.Map:
		fields, isMap := value.(map[string]any)
		if !isMap || place.Key().Kind() != reflect.String {
			break
		}
		filled = reflect.MakeMapWithSize(place, len(fields))
		for key, item := range fields {
			filled.SetMapIndex(reflect.ValueOf(key).Convert(place.Key()), valueFor(item, place.Elem()))
		}
	case reflect.Pointer:
		filled.Set(reflect.New(place.Elem()))
		filled.Elem().Set(valueFor(value, place.Elem()))
	}
	return filled
}

func checkHandler(handler any, where string) {
	if handler == nil {
		return
	}
	if reflect.ValueOf(handler).Kind() != reflect.Func {
		panic(fmt.Sprintf("ozdenui: %s: %s takes a function (func(), func(obj *ozdenui.Object)...), not %T", sourceLine(), where, handler))
	}
}

// *** READING VALUES:

func textOf(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		if typed == Empty {
			return ""
		}
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	}
	return fmt.Sprint(value)
}

func numberOf(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case bool:
		if typed {
			return 1
		}
	case string:
		number, _ := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return number
	case nil:
	default:
		reflected := reflect.ValueOf(value)
		if reflected.CanFloat() {
			return reflected.Float()
		}
		if reflected.CanInt() {
			return float64(reflected.Int())
		}
		if reflected.CanUint() {
			return float64(reflected.Uint())
		}
	}
	return 0
}

// 1, true and a text that is not empty are on (the flags of basic.js are 1 / 0).
func isOn(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		return typed != "" && typed != "0" && typed != "false"
	}
	return numberOf(value) != 0
}

// IsOn reads a flag that came from the page (the result of Ask(), a value of a callback): the kit writes its flags
// as 1 / 0, which come as numbers. It is true for 1, true and a text that is not empty.
func IsOn(value any) bool { return isOn(value) }

// Number reads a number that came from the page as a float64 (0 for nil; a text of a number is read too).
func Number(value any) float64 { return numberOf(value) }

// Text reads a value that came from the page as a text ("" for nil; a whole number without ".0").
func Text(value any) string { return textOf(value) }

// *** THE TEXT OF A NUMBER, made by the page at once: Format("{:,.2f} TL")

type numberFormat struct {
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Thousands int    `json:"thousands"`
	Decimals  *int   `json:"decimals"`
	Percent   int    `json:"percent"`
}

var formatPattern = regexp.MustCompile(`^(.*?)\{:?(,?)(?:\.(\d+))?([fd%]?)\}(.*)$`)

// Format is how a component writes its numbers (the FormatValue of a GraphBox, the ValueFormat of a Stepper): a
// text around one {} place, the page fills it at once.
//
//	"{:,.2f} TL"   1,234.50 TL      "{:.0f} kg"   12 kg      "{:.1%}"   12.3%      "$ {:,}"   $ 1,234
//
// WHY not a function of the program: the page draws at once; the answer of the program would come too late.
func Format(pattern string) any {
	parts := formatPattern.FindStringSubmatch(pattern)
	if parts == nil {
		panic(fmt.Sprintf(`ozdenui: %s: Format(%q): a text around one place, Ex: "{:,.2f} TL", "{:.0f} kg", "{:.1%%}"`, sourceLine(), pattern))
	}
	format := numberFormat{Prefix: parts[1], Suffix: parts[5]}
	if parts[2] == "," {
		format.Thousands = 1
	}
	if parts[3] != "" {
		decimals, _ := strconv.Atoi(parts[3])
		format.Decimals = &decimals
	} else if parts[4] == "d" {
		zero := 0
		format.Decimals = &zero
	}
	if parts[4] == "%" {
		format.Percent = 1
	}
	return format
}

// *** PICTURES:

var pictureEnding = regexp.MustCompile(`(?i)\.(png|jpe?g|gif|webp|bmp)$`)

var svgStart = regexp.MustCompile(`(?s)^\s*(<\?xml.*?\?>\s*)?(<!--.*?-->\s*)*<svg[\s>]`)

var addressPattern = regexp.MustCompile(`^(https?:|data:|blob:|file:)`)

// Picture is the address of a picture for the page: of a file of this computer ("photo.png"), of bytes (a PNG or a
// JPEG read from somewhere) or of an image.Image made by the program. An address of the web stays as it is.
// A param of a component that names a file (IconFile, ImageFile) and the Picture of an Icon take a path by
// themselves; Picture() is for the other places (the "icon" of an item of a list, a picture in a style).
func Picture(source any) string {
	switch typed := source.(type) {
	case string:
		return pictureOfPath(typed)
	case []byte:
		kind := http.DetectContentType(typed)
		if svgStart.Match(typed) {
			kind = "image/svg+xml" // WHY: Go names the bytes of an SVG "text/xml", which the page does not draw
		}
		return dataAddress(typed, kind)
	case image.Image:
		var buffer bytes.Buffer
		if err := png.Encode(&buffer, typed); err != nil {
			currentApp().warn("Picture(): %v", err)
			return ""
		}
		return dataAddress(buffer.Bytes(), "image/png")
	}
	panic(fmt.Sprintf("ozdenui: %s: Picture() takes a path, bytes or an image.Image, not %T", sourceLine(), source))
}

func dataAddress(content []byte, kind string) string {
	return "data:" + kind + ";base64," + base64.StdEncoding.EncodeToString(content)
}

// The address of a file of this computer; a text that is not a file stays as it is (an address, an icon of the kit).
// WHY a data address: the page comes from the program, not from the disk, so it can not read a file by its path.
func pictureOfPath(path string) string {
	if path == "" || addressPattern.MatchString(path) || strings.HasPrefix(path, "comp/") {
		return path
	}
	content, err := os.ReadFile(path)
	if err != nil {
		if pictureEnding.MatchString(path) && current != nil {
			current.warn("the picture %s could not be read: %v", path, err)
		}
		return path
	}
	kind := mime.TypeByExtension(strings.ToLower(filepath.Ext(path)))
	if kind == "" {
		kind = http.DetectContentType(content)
	}
	return dataAddress(content, kind)
}

// A param that names a picture file: iconFile, imageFile, checkIconFile...
func isFileParam(name string) bool {
	return strings.HasSuffix(name, "File")
}
