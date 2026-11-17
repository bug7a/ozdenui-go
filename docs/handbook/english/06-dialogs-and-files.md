# 6. Dialogs, messages and files

## Dialogs return their answer

A dialog in OzdenUI is a function call that returns the answer, as a question in a terminal does:

```go
if ui.Confirm("The file will be deleted.", ui.DialogOptions{Title: "Delete it?", ConfirmText: "Delete", Danger: true}) {
	deleteTheFile()
}

name, ok := ui.Prompt("A name for the folder:", ui.DialogOptions{Value: "Photos", Required: true}) // ok false: cancelled
choice := ui.Ask("Save the changes?", "Save", "Don't save", "Cancel")                              // the text of the button
choice = ui.AskWithTitle("Unsaved changes", "Save them?", "Save", "Don't save")
ui.Alert("The copy is done.", ui.DialogOptions{Title: "Done"})
```

- **The window keeps working while your code waits**: other events and timers run, a clock keeps ticking. Only the
  function that asked waits for the answer.
- `ui.DialogOptions` holds the texts: `Title`, `ConfirmText`, `CancelText`; for `Prompt` also `Value`, `Placeholder`,
  `FieldText` and `Required` (an empty answer is not taken). What you do not write keeps the text of the kit.
- `Danger: true` draws the confirm button as a delete. `ui.Ask()` makes its first button the main one.
- Escape or a click outside is "no": `false` for `Confirm`, `ok == false` for `Prompt`, `""` for `Ask`.
- A dialog is asked from an event or a timer, as here. A goroutine of your own reaches the window only through
  `ui.Do()` (chapter 4).

`examples/dialogs` has a button for each of them and shows the answer as Go writes it.

## Short messages: toasts

```go
ui.ToastSuccess("Saved.")
ui.ToastError("Could not save the file.", "Error") // a title may follow the message
ui.ToastWarning("The disk is almost full.")
ui.ToastInfo("3 new messages")
```

A toast comes in a corner and goes by itself. It never waits: the code goes on at once.

## Bubbles: ShowHint and Tooltip

```go
save := ui.SmartButton(ui.SmartButtonProps{LabelText: "Save"})
ui.Tooltip(ui.TooltipProps{Target: save, HintText: "Saves the file (Ctrl+S)"}) // under the mouse, by itself

ui.ShowHint(save, "Saved.", 2000)                 // at the button now, for 2 s ("top" is its side)
hint := ui.ShowHint(save, "Working", 0, "bottom") // 0: it stays until hint.Remove()
```

- A `Tooltip` with `Trigger: "click"` is a popover: a click on its target opens it, and it stays until a click outside.
  `HintPosition` is its side: `"top"`, `"bottom"`, `"left"`, `"right"`.
- In v0.2 a bubble shows a text. A bubble that holds objects of your own is not there yet.

## The file dialogs of the system

```go
path := ui.OpenFile(ui.FileOptions{Types: []string{"png", "jpg"}}) // a path, or "" (cancelled)
paths := ui.OpenFiles()                                            // a []string, empty when cancelled
target := ui.SaveFile(ui.FileOptions{FileName: "report.csv", Types: []string{"csv"}})
folder := ui.ChooseFolder(ui.FileOptions{Title: "Where are the photos?"})
```

- They are the dialogs of the computer (on Linux: `zenity`), so the paths are real: read and write them with `os`.
- `ui.FileOptions{Title, Directory, FileName, Types}`: `Directory` is the folder it opens in, `FileName` the name
  `SaveFile` offers, `Types` the endings it shows.
- `SaveFile` only gives the path: your program writes the file.

### A small editor

```go
package main

import (
	"os"
	"path/filepath"

	ui "github.com/bug7a/ozdenui-go"
)

const pageColor = "#F4F5F7"

var fileTypes = []string{"txt", "md"}

var (
	path string // the file of the note; "" until it is opened or saved
	note *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Note", Width: 560, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 12, Padding: 20, Align: "left top"}, func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Open", OnClick: openNote})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan", OnClick: saveNote})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Clear", Tone: "orange", OnClick: clearNote})
		})
		note = ui.TextAreaField(ui.TextAreaFieldProps{TitleText: "Note", Width: "100%", Height: 260})
	})
}

func openNote() {
	chosen := ui.OpenFile(ui.FileOptions{Types: fileTypes})
	if chosen == "" {
		return
	}
	content, err := os.ReadFile(chosen)
	if err != nil {
		ui.ToastError(err.Error(), "Could not open it")
		return
	}
	path = chosen
	note.SetText(string(content))
	note.Set("titleText", filepath.Base(path))
}

func saveNote() {
	if path == "" {
		path = ui.SaveFile(ui.FileOptions{FileName: "note.txt", Types: fileTypes})
	}
	if path == "" {
		return
	}
	if err := os.WriteFile(path, []byte(note.Text()), 0o644); err != nil {
		ui.ToastError(err.Error(), "Could not save it")
		return
	}
	note.Set("titleText", filepath.Base(path))
	ui.ToastSuccess("Saved.", filepath.Base(path))
}

func clearNote() {
	if ui.Confirm("The text in the window goes away.", ui.DialogOptions{Title: "Clear the note?", ConfirmText: "Clear", Danger: true}) {
		note.SetText("")
	}
}
```

`note.Text()` asks nothing of the page: the field tells the program its text while the user types.
`examples/csv-viewer`, `examples/image-viewer` and `examples/markdown-editor` open their files the same way.

## A wait over the window: LoadingScreen

```go
loading := ui.LoadingScreen(ui.LoadingScreenProps{TitleText: "Copying", MessageText: "Starting...", Progress: 0})
go func() {
	for step := 1; step <= 4; step++ {
		copyOneFile(step) // the slow work, in the goroutine
		ui.Do(func() {
			loading.Call("setProgress", float64(step)/4) // 0 .. 1
			loading.Call("setMessageText", fmt.Sprintf("File %d of 4", step))
		})
	}
	ui.Do(func() { loading.Call("close") })
}()
```

It covers the whole window, wherever it is made. `Progress: -1` (its default) shows no bar; `MinDuration` (ms) keeps it
a moment, so a fast job does not blink; `RemoveOnClose: true` removes it after `close`, for a screen used once.

## A place to drop files: FilePicker

`FilePicker` is a place to drop files on, or to click for a file. `picker.Get("files")` is its list, a `[]any` of
maps: `name`, `path`, `size`, `sizeText`, `type`, `extension`, `lastModified`.

A click opens the file dialog of the system, so the chosen files come with their real paths; `picker.Paths()` lists
them. A file dropped on it comes from the page, and the page never knows where a file is: its `path` is only its name,
and it is not in `Paths()`.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	list *ui.Object
	rows []*ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Files", Width: 520, Height: 460})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14, Padding: 24, Align: "left top"}, func() {
		ui.FilePicker(ui.FilePickerProps{Width: "100%", ZoneHeight: 140, Multiple: true, ShowList: false,
			TitleText: "Drop files here", OnChange: showFiles})
		list = ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 6, Align: "left top"}, func() {
			rows = append(rows, ui.Label(ui.Props{Text: "No files yet.", FontSize: 14, TextColor: "#5A6270"}))
		})
	})
}

func showFiles(picker *ui.Object) {
	for _, row := range rows {
		row.Remove()
	}
	rows = nil
	files, _ := picker.Get("files").([]any)
	list.Add(func() {
		for _, item := range files {
			file, _ := item.(map[string]any)
			text := fmt.Sprintf("%v (%v)", file["path"], file["sizeText"])
			rows = append(rows, ui.Label(ui.Props{Text: text, FontSize: 14}))
		}
	})
}
```

- `Accept: ".csv,.txt"` or `"image/*"`; `MaxSize` (bytes) and `MaxFiles`: what does not pass goes to
  `OnError`. `Multiple: true` adds to the list; without it a new file takes the place of the old one.
- `picker.Call("clear")` empties it; `TitleText`, `DescriptionText`: its texts; `IconFile`: a picture of yours.
- `picker.AddFiles("report.pdf", "notes.txt")` adds files of this computer from the code. `OnOpenDialog:
  func(picker *ui.Object)` puts a dialog of your own in the place of the one of the system.
- The bytes of a dropped file do not reach Go. For the content of a file, read it from its path (`os.ReadFile`): a
  click gives one, a drop does not.

## Next

[Chapter 7: Looks](07-looks.md): dark windows, the minimal look and your colors.
