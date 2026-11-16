# 8. Data and pictures

The components of this chapter show what a Go program usually has: numbers, rows, logs, texts, code and pictures.
None of them needs another module: the page draws them.

## Charts: GraphBox

```go
ui.GraphBox(ui.GraphBoxProps{Type: "bar", Labels: []string{"Jan", "Feb", "Mar"}, Values: []int{1200, 1900, 700},
	FormatValue: "{:,.0f} TL"})
```

- `Type`: `"line"`, `"area"`, `"bar"`, `"pie"`, `"doughnut"`, `"radar"`, `"scatter"`, `"bubble"`.
- `Values` is a slice of any numbers. More lines or bars:
  `Datasets: []ui.P{{"label": "2025", "data": []int{8, 11, 9}}, {"label": "2026", "data": []int{10, 14, 12}}}`.
- `FormatValue` is a format text, not a function, because the page needs the text of a number while it draws:
  `"{}"`, `"{:.2f}"`, `"{:,}"`, `"{:,.0f} TL"`, `"{:.1%}"`. A method that takes a format takes it through `ui.Format`:
  `graph.Call("setFormatValue", ui.Format("{:.0f} kg"))`.
- Later: `graph.Call("setValues", []int{...})`, `graph.Call("setType", "line")`, `graph.Call("addValue", "Apr", 15)`.
  `OnClick: func(graph *ui.Object, item ui.P)`: `item["label"]`, `item["value"]`.

### A live chart and a log

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	points    = 30
	firstOnes = 12 // values drawn at the start, one second apart
)

var (
	chart   *ui.Object
	logView *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A sensor", Width: 560, Height: 520})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 12, Padding: 20}, func() {
		chart = ui.GraphBox(ui.GraphBoxProps{Type: "area", TitleText: "Temperature", Labels: []string{}, Values: []float64{},
			Curve: "smooth", ShowPoints: false, MaxCount: points, FormatValue: "{:.1f}°", Width: "100%", Height: 240})
		logView = ui.LogView(ui.LogViewProps{Fill: true})
	})
	for back := firstOnes; back > 0; back-- {
		measure(time.Now().Add(-time.Duration(back) * time.Second))
	}
	ui.Every(1000, func() { measure(time.Now()) })
}

func measure(at time.Time) {
	value := math.Round((20+rand.Float64()*6-3)*10) / 10
	chart.Call("addValue", at.Format("15:04:05"), value) // the oldest goes after `points` values
	if value > 22 {
		logView.Call("warn", fmt.Sprintf("Warm: %.1f°", value))
		return
	}
	logView.Call("info", fmt.Sprintf("%.1f°", value))
}
```

`examples/charts` has a live chart, bars you can click, a doughnut, two lines with a target, and a chart saved as a
PNG file (`graph.Ask("getImage", "image/png", 2)` gives its data address).

## Tables

A row may be a struct of your program: its `json` tags are the keys the columns name. `OnSelect` gives the row back
in your type:

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

type person struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Orders int    `json:"orders"`
}

var people = []person{{"Ada", "London", 12}, {"Alan", "Manchester", 7}, {"Grace", "New York", 21}}

var columns = []ui.P{
	{"name": "Name", "dataTitle": "name", "width": 180, "sortable": 1},
	{"name": "City", "dataTitle": "city", "width": 160, "sortable": 1},
	{"name": "Orders", "dataTitle": "orders", "dataType": "integer", "width": 100, "sortable": 1},
}

var (
	table  *ui.Object
	chosen *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "People", Width: 520, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 12, Padding: 20, Align: "left top"}, func() {
		table = ui.SmartTable(ui.SmartTableProps{TitleDataList: columns, ItemDataList: people, Width: "100%",
			Height: 260, OnSelect: showPerson})
		ui.HGroup(ui.Props{Hug: true, Gap: 12, Align: "left center"}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Add one", OnClick: addPerson})
			chosen = ui.Label(ui.Props{Text: "Click a row.", FontSize: 15})
		})
	})
}

func showPerson(table *ui.Object, row person) {
	chosen.SetText(fmt.Sprintf("%s, %s: %d orders", row.Name, row.City, row.Orders))
}

func addPerson() {
	people = append(people, person{fmt.Sprintf("Person %d", len(people)+1), "İzmir", 0})
	table.Set("itemDataList", people) // the whole list again: the table draws it
}
```

- **`SmartTable`** makes only the rows that are seen, so a very big list stays fast (`examples/csv-viewer` has a button
  for 100000 rows). It has a filter box and sorting on its titles. A row may also be a `ui.P`.
- A column: `name` (its title), `dataTitle` (the key of the row), `width`, `dataType` (`"string"`, `"integer"`...),
  `sortable`. The values are shown as text: format them in Go (`fmt.Sprintf("%.2f", price)`).
- **`DataTable`** is for up to a few hundred rows; its columns fit their texts.
- `examples/csv-viewer` opens a CSV file in a SmartTable.

## Pictures

A picture may be a file path, bytes (a PNG or a JPEG), an `image.Image` made by your program, or a web address.
`ui.Picture()` turns any of them into an address the page can show:

```go
ui.Icon(ui.Props{Width: 48, Height: 48, Picture: "logo.png"})                // Picture takes a path, bytes or an image.Image
view := ui.ImageView(ui.ImageViewProps{Height: 400, ImageFile: "photo.jpg"}) // a field that ends in File: a path
view.Call("setImage", ui.Picture(img))                                       // in a method: ui.Picture()
```

- An `ImageView` zooms with the wheel or two fingers and moves by dragging. A new picture keeps the zoom and the place,
  so a program can draw its plot again and again.
- The page can not read a file of this computer by its path, so Go reads it and gives its content as a `data:`
  address. A field that ends in `File` and the `Picture` of an `Icon` do it by themselves; in any other place (an item
  of a list, a style, a method) give `ui.Picture(path)`.
- The address carries the whole picture: for one you show many times, keep `ui.Picture(...)` in a variable.

```go
package main

import (
	"image"
	"image/color"

	ui "github.com/bug7a/ozdenui-go"
)

const size = 240

var (
	view  *ui.Object
	shade = 0
)

func main() {
	ui.Run(start, ui.Window{Title: "A picture of Go", Width: 420, Height: 420})
}

// gradient draws a picture pixel by pixel, as a program of yours would.
func gradient(shade int) image.Image {
	picture := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			picture.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(shade), 255})
		}
	}
	return picture
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 12, Padding: 20}, func() {
		view = ui.ImageView(ui.ImageViewProps{Width: 300, Height: 300, ImageFile: ui.Picture(gradient(shade))})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Draw it again", OnClick: drawAgain})
	})
}

func drawAgain() {
	shade = (shade + 80) % 256
	view.Call("setImage", ui.Picture(gradient(shade)))
}
```

`examples/image-viewer` opens a picture file, or draws a plot again and again.

## Logs: LogView

```go
logView := ui.LogView(ui.LogViewProps{Height: 300})
logView.Call("info", "Reading the file")
logView.Call("success", "Done")
logView.Call("error", "The file is not there")
log.SetOutput(logView.Writer("info"))                    // the log package writes into it
fmt.Fprintln(logView.Writer("warn"), "The disk is full") // any io.Writer: also from a goroutine
```

The levels are `"info"`, `"success"`, `"warn"`, `"error"` and `"debug"`. The view follows the last line until the
reader scrolls up. A text is never read as HTML. `logView.Writer(level)` may be used from any goroutine: it goes
through `ui.Do()` by itself, and every `Write` is a line (or more). `examples/job-log`.

## Texts, JSON and code

| | |
|---|---|
| `ui.MarkdownView(ui.MarkdownViewProps{MarkdownText: "# A report\n\nSome **bold** words"})` | `view.Call("appendText", piece)` adds to the end: the answer of an AI as it comes; `view.Set("markdownText", text)` |
| `ui.JsonView(ui.JsonViewProps{Value: data, Height: 400, OnSelect: fn})` | any Go value (a struct keeps the order of its fields); a click: `OnSelect: func(view *ui.Object, entry ui.P)`, `entry["pathText"]` (`users[0].name`) |
| `ui.JsonEditor(ui.JsonEditorProps{Value: rows, Height: 360, OnChange: fn})` | a list or a map of JSON written in a table; `OnChange: func(editor *ui.Object, list []product)` gets the whole value in your type |
| `ui.CodeEditor(ui.CodeEditorProps{Value: "x := 1", Language: "go", Height: 300, OnChange: fn})` | its text is asked: `editor.Ask("getValue")`; Go, Python, JavaScript, JSON, HTML, CSS, SQL, Bash |
| `ui.ChatView(ui.ChatViewProps{EmptyText: "Ask something", OnSend: answer})` | `OnSend: func(chat *ui.Object, text string)`; an answer: `chat.Call("addMessage", ui.P{"key": key, "role": "assistant"})`, then `chat.Call("appendText", key, piece)`; `chat.Set("busy", true)` shows the dots and Stop |
| `ui.WebView(ui.WebViewProps{HtmlText: page, Height: 400})` | a web page (`Url`) or an HTML text with its scripts; `OnMessage: func(view *ui.Object, data T)` gets what its page sends with `parent.postMessage` |

- **A big or fast value is asked, not kept**: a `CodeEditor` does not send its text with every key. Ask it when you
  need it, in its `OnChange` or when the user saves (`ChangeDelay: 120` waits for the typing to rest).
- A Markdown text is never read as HTML and no script of it runs. The page can not read a file by its path: a picture
  of this computer goes into the text as `ui.Picture(path)` (`examples/markdown-editor` does it). A picture with an
  `http` address stays a link (`RemotePictures: true` for a text you trust).
- The pieces of an answer of a `ChatView` may come from a goroutine through `ui.Do()`: the window keeps working, and
  the view follows the end until the reader scrolls up (`examples/chat`).
- `examples/markdown-editor`, `examples/json-viewer` and `examples/editor` (a small editor that runs its Python
  program) show them.

## Drawing: Canvas

```go
ui.Canvas(ui.CanvasProps{Width: "100%", Height: 300, Shapes: []ui.P{
	{"key": "sun", "type": "circle", "x": 100, "y": 80, "radius": 30, "fill": "orange", "draggable": 1},
	{"key": "path", "type": "line", "x1": 20, "y1": 200, "x2": 260, "y2": 200, "color": "tomato", "lineWidth": 4},
	{"key": "title", "type": "text", "x": 20, "y": 12, "text": "Drag the sun", "fontSize": 16},
}, OnShapeMove: func(canvas *ui.Object, shape struct{ Key string }) { fmt.Println("moved:", shape.Key) }})
```

The page moves a `draggable` shape by itself; `OnShapeMove` tells Go where it stopped. `Drawing: true` gives the
mouse a pen. Its methods (`addShape`, `moveShape`, `removeShape`, `getImage`...) are called with `Call`.

## Next

[Chapter 9: Your own component](09-your-own-component.md).
