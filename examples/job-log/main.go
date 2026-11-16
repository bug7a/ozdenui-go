/* Bismillah */

// A job in the background and its output in the window: a goroutine, ui.Do() and the levels of a LogView.
package main

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageLight = "#FFFFFF"
	pageDark  = "#1E1E1E"
	stepTime  = 150 * time.Millisecond
)

var files = []string{"photo.jpg", "notes.txt", "disk.bin", "report.pdf", "song.mp3", "backup.zip"}

var (
	runButton *ui.Object
	bar       *ui.Object
	output    *ui.Object
	status    *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A job and its log", Width: 680, Height: 520})
}

func start(page *ui.Object) {
	page.Set("color", map[bool]string{false: pageLight, true: pageDark}[ui.IsDark()])
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 10, Padding: 12, Align: "left top", Shrink: 0}, func() {
			ui.HGroup(ui.Props{Hug: true, Gap: 8, Align: "left center"}, func() {
				runButton = ui.SmartButton(ui.SmartButtonProps{LabelText: "Run the job", Tone: "cyan", OnClick: runJob})
				ui.SmartButton(ui.SmartButtonProps{LabelText: "Clear", OnClick: func() { output.Call("clearLines") }})
				ui.CheckBox(ui.CheckBoxProps{LabelText: "Only warnings and errors", OnChange: showLevels})
			})
			bar = ui.ProgressBar(ui.ProgressBarProps{TitleText: "Files", MaxValue: len(files), Width: "100%"})
		})
		output = ui.LogView(ui.LogViewProps{Fill: true, Unframed: true, EmptyText: "Run the job: its output comes here."})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Ready", Items: []ui.P{{"key": "files", "text": "0 files"}}})
	})
}

func showLevels(box *ui.Object) {
	levels := []string{}
	if box.Checked() {
		levels = []string{"warn", "error"}
	}
	output.Call("setLevels", levels)
}

func runJob() {
	runButton.Set("enabled", false)
	bar.SetValue(0)
	status.Call("setBusy", true)
	status.Call("setMessageText", "The job is running...")
	go job()
}

// write adds a line to the view from the goroutine of the job: level is "info", "success", "warn" or "error".
func write(level string, format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	ui.Do(func() { output.Call(level, text) })
}

func job() {
	write("info", "The job started")
	for index, name := range files {
		time.Sleep(stepTime) // in the goroutine of the job: the window does not wait
		switch {
		case strings.HasSuffix(name, ".bin"):
			write("error", "%s can not be read: permission denied", name)
		case strings.HasSuffix(name, ".zip"):
			write("warn", "%s is very large", name)
		default:
			write("info", "Copied %s", name)
		}
		count := index + 1
		ui.Do(func() {
			bar.SetValue(count)
			status.Call("setItemText", "files", fmt.Sprintf("%d files", count))
		})
	}
	write("info", "The job ended")
	write("success", "Done")
	ui.Do(func() {
		status.Call("setBusy", false)
		status.Call("setMessageText", "Ready") // the lasting message: "Done" goes back to it
		status.Call("setMessageText", "Done", "success", 3000)
		runButton.Set("enabled", true)
	})
}
