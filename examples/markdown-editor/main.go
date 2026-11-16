/* Bismillah */

// A small Markdown editor: write on the left, read on the right. Open a file with its pictures, save it, a dark theme.
// A click on a link is told at the bottom, not opened.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	lightGround  = "#FFFFFF"
	darkGround   = "#1C1C1B"
	typingRestMs = 120 // the text is drawn again when the typing rests
)

const firstText = "# A small report\n\n" +
	"The **sales** went up this month. Read the [notes](https://example.com/notes).\n\n" +
	"## The numbers\n\n" +
	"| Product | Count | Total |\n|:---|:--:|---:|\n| Olive oil | 120 | 7,200 TL |\n| Honey | 45 | 3,150 TL |\n\n" +
	"## Next steps\n\n- [x] Write the report\n- [ ] Send it\n\n" +
	"> The numbers of the last week are not final.\n\n" +
	"```python\ntotal = sum(row[\"total\"] for row in rows)\nprint(f\"Total: {total:,} TL\")\n```\n"

// ![alt](address): a picture of the text.
var picturePattern = regexp.MustCompile(`(!\[[^\]]*\]\()([^)\s]+)`)

var (
	path          = ""
	pictureFolder = "" // the folder of the opened file: its pictures are next to it

	page    *ui.Object
	tools   *ui.Object
	split   *ui.Object
	editor  *ui.Object
	preview *ui.Object
	status  *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A small Markdown editor", Width: 960, Height: 640})
}

func start(window *ui.Object) {
	page = window
	page.Set("color", lightGround)
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		tools = ui.ToolBar(ui.ToolBarProps{Items: []any{
			ui.P{"key": "open", "icon": "open", "hint": "Open...", "shortcut": "Ctrl+O", "onClick": openDocument},
			ui.P{"key": "save", "icon": "save", "hint": "Save", "shortcut": "Ctrl+S", "onClick": saveDocument},
			" ",
			ui.P{"key": "dark", "text": "Dark theme", "toggle": true, "onClick": toggleTheme},
		}})
		split = ui.Splitter(ui.SplitterProps{Fill: true, FirstSize: "45%", MinFirst: 160, MinSecond: 160})
		split.First(func() {
			editor = ui.CodeEditor(ui.CodeEditorProps{Height: "100%", Value: firstText, ShowLineNumbers: false, ChangeDelay: typingRestMs,
				Unframed: true, OnChange: showText})
		})
		split.Second(func() {
			preview = ui.MarkdownView(ui.MarkdownViewProps{MarkdownText: firstText, OnLinkClick: showLink})
		})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Ready", Items: []any{ui.P{"key": "words", "text": wordsOf(firstText)}}})
	})
}

func wordsOf(text string) string {
	return fmt.Sprintf("%d words", len(strings.Fields(text)))
}

// The text with the pictures of this computer as addresses: the page can not read a file by its path.
func withPictures(text string) string {
	if pictureFolder == "" {
		return text
	}
	return picturePattern.ReplaceAllStringFunc(text, func(found string) string {
		parts := picturePattern.FindStringSubmatch(found)
		address := parts[2]
		if strings.Contains(address, ":") {
			return found // an address of the web or a data: address
		}
		return parts[1] + ui.Picture(filepath.Join(pictureFolder, address))
	})
}

func showText(code *ui.Object) {
	text := fmt.Sprint(code.Ask("getValue")) // WHY Ask: the editor does not send its text back by itself
	preview.Set("markdownText", withPictures(text))
	status.Call("setItemText", "words", wordsOf(text))
}

func showLink(view *ui.Object, address string) {
	status.Call("setMessageText", "Link: "+address, "info", 4000)
}

// WHY not ui.SetTheme(): it would draw the page again, and the text would be made again too.
func toggleTheme(bar *ui.Object) {
	theme, ground := "light", lightGround
	if contains(bar.Get("checkedKeys"), "dark") {
		theme, ground = "dark", darkGround
	}
	page.Set("color", ground)
	for _, obj := range []*ui.Object{tools, split, editor, preview, status} {
		obj.Call("setTheme", theme)
	}
}

func contains(list any, key string) bool {
	items, _ := list.([]any)
	for _, item := range items {
		if item == key {
			return true
		}
	}
	return false
}

func openDocument() {
	chosen := ui.OpenFile(ui.FileOptions{Types: []string{"md"}})
	if chosen == "" {
		return
	}
	content, err := os.ReadFile(chosen)
	if err != nil {
		status.Call("setMessageText", err.Error(), "error")
		return
	}
	text := string(content)
	path = chosen
	pictureFolder = filepath.Dir(chosen)
	editor.SetValue(text)
	preview.Set("markdownText", withPictures(text))
	preview.Call("scrollToTop")
	status.Set("messageText", path)
	status.Call("setItemText", "words", wordsOf(text))
}

func saveDocument() {
	if path == "" {
		path = ui.SaveFile(ui.FileOptions{FileName: "report.md", Types: []string{"md"}})
	}
	if path == "" {
		return
	}
	text := fmt.Sprint(editor.Ask("getValue"))
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		status.Call("setMessageText", err.Error(), "error")
		return
	}
	status.Call("setMessageText", "Saved", "success", 3000)
}
