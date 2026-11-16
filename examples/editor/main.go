/* Bismillah */

// A small code editor: write Python, run it (F5) and read its output. A menu, a row of buttons, two areas with a line
// between them, the line and the column of the cursor, a light and a dark theme, the keys Ctrl + S and F5 (Cmd on a Mac).
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	lightGround  = "#FFFFFF"
	darkGround   = "#1C1C1B"
	outputHeight = 170
)

const firstText = `"""Press F5 to run this program."""

names = ["Ada", "Can", "Deniz"]

for number, name in enumerate(names, start=1):
    print(number, name)

print("Done")
`

var (
	isDark      = false
	path        = ""
	lineNumbers = true

	page   *ui.Object
	menu   *ui.Object
	tools  *ui.Object
	split  *ui.Object
	editor *ui.Object
	output *ui.Object
	status *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A small code editor", Width: 760, Height: 600})
}

func start(window *ui.Object) {
	page = window
	page.Set("color", lightGround)
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		menu = ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{
			{"text": "File", "items": []ui.P{
				{"text": "Open...", "key": "open", "shortcut": "Ctrl+O", "onClick": openProgram},
				{"text": "Save", "key": "save", "shortcut": "Ctrl+S", "onClick": saveProgram},
			}},
			{"text": "Run", "items": []ui.P{
				{"text": "Run the program", "key": "run", "shortcut": "F5", "onClick": runProgram},
				{"text": "Clear the output", "key": "clear", "onClick": clearOutput},
			}},
			{"text": "View", "items": []ui.P{
				{"text": "Dark theme", "key": "dark", "checked": false, "onClick": toggleTheme},
				{"text": "Line numbers", "key": "numbers", "checked": true, "onClick": toggleLineNumbers},
			}},
		}})
		// WHY the same shortcuts as the menu: the hint of a button tells them; a key that is in both bars runs once.
		tools = ui.ToolBar(ui.ToolBarProps{Items: []any{
			ui.P{"key": "open", "icon": "open", "hint": "Open...", "shortcut": "Ctrl+O", "onClick": openProgram},
			ui.P{"key": "save", "icon": "save", "hint": "Save", "shortcut": "Ctrl+S", "onClick": saveProgram},
			"-",
			ui.P{"key": "run", "icon": "play", "text": "Run", "shortcut": "F5", "onClick": runProgram},
			ui.P{"key": "clear", "icon": "delete", "hint": "Clear the output", "onClick": clearOutput},
		}})
		split = ui.Splitter(ui.SplitterProps{Fill: true, Direction: "vertical", SecondSize: outputHeight, MinFirst: 80, MinSecond: 60})
		split.First(func() {
			editor = ui.CodeEditor(ui.CodeEditorProps{Height: "100%", Language: "python", Value: firstText, Unframed: true,
				OnCursorChange: showCursor})
		})
		split.Second(func() {
			output = ui.LogView(ui.LogViewProps{Height: "100%", Unframed: true, EmptyText: "F5 runs the program."})
		})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Ready", Items: []any{ui.P{"key": "place", "text": "Line 1, column 1"}, "Python"}})
	})
}

// WHY not ui.SetTheme(): it would draw the page again, and the text and the output would be made again too.
func toggleTheme(bar *ui.Object) {
	isDark = !isDark
	theme, ground := "light", lightGround
	if isDark {
		theme, ground = "dark", darkGround
	}
	page.Set("color", ground)
	for _, obj := range []*ui.Object{menu, tools, split, editor, output, status} {
		obj.Call("setTheme", theme)
	}
	bar.Call("setItemChecked", "dark", isDark)
}

func toggleLineNumbers(bar *ui.Object) {
	lineNumbers = !lineNumbers
	editor.Set("showLineNumbers", lineNumbers)
	bar.Call("setItemChecked", "numbers", lineNumbers)
}

func showCursor(code *ui.Object) {
	status.Call("setItemText", "place", fmt.Sprintf("Line %d, column %d", code.Int("cursorLine"), code.Int("cursorColumn")))
}

func clearOutput() {
	output.Call("clearLines")
}

// The text as it is on the page now: the editor does not send it back by itself.
func programText() string {
	return fmt.Sprint(editor.Ask("getValue"))
}

func openProgram() {
	chosen := ui.OpenFile(ui.FileOptions{Types: []string{"py"}})
	if chosen == "" {
		return
	}
	content, err := os.ReadFile(chosen)
	if err != nil {
		status.Call("setMessageText", err.Error(), "error")
		return
	}
	editor.SetValue(string(content))
	path = chosen
	status.Set("messageText", path)
}

func saveProgram() {
	if path == "" {
		path = ui.SaveFile(ui.FileOptions{FileName: "program.py", Types: []string{"py"}})
	}
	if path == "" {
		return
	}
	if err := os.WriteFile(path, []byte(programText()), 0o644); err != nil {
		status.Call("setMessageText", err.Error(), "error")
		return
	}
	status.Call("setMessageText", "Saved", "success", 3000)
}

func runProgram() {
	code := programText()
	output.Call("clearLines")
	status.Set("busy", true)
	status.Set("messageText", "Running...")
	go runInPython(code) // the window keeps working while the program runs
}

func runInPython(code string) {
	command := exec.Command(pythonCommand(), "-u", "-c", code)
	pipe, err := command.StdoutPipe()
	command.Stderr = command.Stdout
	if err == nil {
		err = command.Start()
	}
	if err != nil {
		ui.Do(func() { endRun(err.Error()) })
		return
	}
	lines := bufio.NewScanner(pipe)
	for lines.Scan() {
		line := lines.Text()
		ui.Do(func() { output.Call("addLine", line) })
	}
	err = command.Wait()
	ui.Do(func() {
		if exit, ok := err.(*exec.ExitError); ok {
			endRun(fmt.Sprintf("The program ended with the code %d", exit.ExitCode()))
			return
		}
		if err != nil {
			endRun(err.Error())
			return
		}
		output.Call("success", "The program ended")
		endRun("")
	})
}

func endRun(problem string) {
	if problem != "" {
		output.Call("error", problem)
	}
	status.Set("busy", false)
	status.Set("messageText", "Ready")
}

// python3 on macOS and Linux, python on Windows.
func pythonCommand() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	return "python"
}
