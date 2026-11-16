/* Bismillah */

// A small JSON viewer: open a file, click a row to see its path, right click to copy its value or its path.
// go run ./examples/json-viewer some.json: that file.
package main

import (
	"fmt"
	"os"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageLight = "#FFFFFF"
	pageDark  = "#1E1E1E"
	copyShown = 60 // letters of a copied text in the status bar
)

type address struct {
	City string `json:"city"`
	Zip  string `json:"zip"`
}

type user struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Active  bool     `json:"active"`
	Roles   []string `json:"roles"`
	Address *address `json:"address"`
}

type answer struct {
	Status   string  `json:"status"`
	Count    int     `json:"count"`
	PageTime float64 `json:"page_time"`
	Cached   bool    `json:"cached"`
	NextPage *int    `json:"next_page"`
	Users    []user  `json:"users"`
}

var firstValue = answer{
	Status:   "ok",
	Count:    2,
	PageTime: 0.042,
	Users: []user{
		{ID: 7, Name: "Ada Yılmaz", Email: "ada@example.com", Active: true, Roles: []string{"admin", "editor"},
			Address: &address{City: "İzmir", Zip: "35000"}},
		{ID: 12, Name: "Can Demir", Email: "can@example.com", Roles: []string{}},
	},
}

var (
	firstFile string
	jsonView  *ui.Object
	status    *ui.Object
)

func main() {
	if len(os.Args) > 1 {
		firstFile = os.Args[1]
	}
	ui.Run(start, ui.Window{Title: "A small JSON viewer", Width: 720, Height: 560})
}

func start(page *ui.Object) {
	page.Set("color", map[bool]string{false: pageLight, true: pageDark}[ui.IsDark()])
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		ui.ToolBar(ui.ToolBarProps{OnClick: toolClicked, Items: []any{
			ui.P{"key": "open", "icon": "open", "hint": "Open...", "shortcut": "Ctrl+O"},
			"-",
			ui.P{"key": "expand", "text": "Open all"},
			ui.P{"key": "collapse", "text": "Close all"},
		}})
		jsonView = ui.JsonView(ui.JsonViewProps{Fill: true, Value: firstValue, ExpandLevel: 2,
			OnSelect: showEntry, OnCopy: showCopied})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Click a row: its path comes here."})
	})
	if firstFile != "" {
		readFile(firstFile)
	}
}

func toolClicked(bar *ui.Object, item ui.P) {
	switch item["key"] {
	case "open":
		openJson()
	case "expand":
		jsonView.Call("expandAll")
	case "collapse":
		jsonView.Call("collapseAll")
	}
}

func showEntry(view *ui.Object, entry ui.P) {
	status.Call("setMessageText", fmt.Sprintf("%v  (%v)", entry["pathText"], entry["kind"]))
}

func showCopied(view *ui.Object, text string) {
	if runes := []rune(text); len(runes) > copyShown {
		text = string(runes[:copyShown])
	}
	status.Call("setMessageText", "Copied: "+text, "success", 3000)
}

func readFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		status.Call("setMessageText", err.Error(), "error")
		return
	}
	isJson, _ := jsonView.Ask("setJsonText", string(content)).(float64) // the kit answers 1 or 0
	if isJson == 1 {
		status.Call("setMessageText", path)
		return
	}
	status.Call("setMessageText", "It is not JSON: "+path, "error")
}

func openJson() {
	path := ui.OpenFile(ui.FileOptions{Types: []string{"json"}})
	if path != "" {
		readFile(path)
	}
}
