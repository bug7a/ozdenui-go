/* Bismillah */

/*

dialogs.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The dialogs and the toasts of the kit, and the file dialogs of the system.

A dialog returns its answer, like a question in a terminal: the code waits for it, and the window keeps working
meanwhile (other events and timers run).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// DialogOptions are the texts of a dialog; what is not given keeps the text of the kit.
type DialogOptions struct {
	Title       string
	ConfirmText string // the button that says yes ("Confirm", "OK")
	CancelText  string
	Danger      bool   // Confirm(): the yes button is red (a delete)
	Value       string // Prompt(): the text in the field at the start
	Placeholder string // Prompt()
	FieldText   string // Prompt(): the title of the field
	Required    bool   // Prompt(): an empty answer is not taken
}

func dialogOptions(options []DialogOptions) DialogOptions {
	if len(options) > 0 {
		return options[0]
	}
	return DialogOptions{}
}

func showDialog(method string, params P) any {
	app := currentApp()
	app.check(method + " dialog")
	spec := kit["Dialog"]
	values := P{}
	takeWindowLook(app, spec, values)
	if name, has := values["styleName"]; has {
		params["styleName"] = name
	}
	for key, value := range params {
		if text, ok := value.(string); ok && text == "" && key != "messageText" && key != "value" {
			delete(params, key) // the text of the kit
		}
	}
	return app.callStatic("Dialog", method, []any{params}, true, "")
}

// Alert shows a message and waits until it is closed.
func Alert(message string, options ...DialogOptions) {
	given := dialogOptions(options)
	showDialog("alert", P{"titleText": given.Title, "messageText": message, "buttonText": given.ConfirmText})
}

// Confirm asks a yes / no question: if ozdenui.Confirm("Delete it?", ozdenui.DialogOptions{Danger: true}) { ... }
func Confirm(message string, options ...DialogOptions) bool {
	given := dialogOptions(options)
	if given.Title == "" {
		given.Title = "Are you sure?"
	}
	return isOn(showDialog("confirm", P{
		"titleText": given.Title, "messageText": message, "confirmText": given.ConfirmText,
		"cancelText": given.CancelText, "danger": flagAsNumber(given.Danger),
	}))
}

// Prompt asks for a text. ok is false when the dialog was cancelled.
//
//	name, ok := ozdenui.Prompt("The name of the list:")
func Prompt(message string, options ...DialogOptions) (text string, ok bool) {
	given := dialogOptions(options)
	answer := showDialog("prompt", P{
		"titleText": given.Title, "messageText": message, "value": given.Value, "placeholder": given.Placeholder,
		"fieldText": given.FieldText, "required": flagAsNumber(given.Required), "confirmText": given.ConfirmText,
		"cancelText": given.CancelText,
	})
	if answer == nil {
		return "", false
	}
	return textOf(answer), true
}

// Ask shows a message with buttons of the program and returns the text of the pressed one ("": closed with Escape or
// a click outside). The first button is the main one.
//
//	switch ozdenui.Ask("Save the changes?", "Save", "Don't save", "Cancel") { ... }
func Ask(message string, buttons ...string) string { return AskWithTitle("", message, buttons...) }

// AskWithTitle is Ask with a title over the message.
func AskWithTitle(title, message string, buttons ...string) string {
	app := currentApp()
	app.check("Ask()")
	if len(buttons) == 0 {
		buttons = []string{"OK"}
	}
	list := []P{}
	for i, text := range buttons {
		kind := "secondary"
		if i == 0 {
			kind = "primary"
		}
		list = append(list, P{"text": text, "key": text, "type": kind})
	}
	pending := &answer{}
	params := P{"messageText": message, "buttons": list}
	if title != "" {
		params["titleText"] = title
	}
	params["onClose"] = func(_ *Object, reason string, key string) {
		if reason != "button" {
			key = ""
		}
		pending.give(key)
	}
	dialog := Component("Dialog", params)
	// WHY by the id of the dialog: when the page can not make it, its answer is nil (the code would wait for ever)
	return textOf(app.waitForAnswer(dialog.id, pending))
}

// *** TOASTS: small messages that never wait.

func showToast(kind, message string, title []string) {
	app := currentApp()
	theme := app.theme
	if app.window.Style == "minimal" {
		theme = "minimal"
		if app.window.Accent != "" {
			theme = accentStyleName
		}
		if app.theme == "dark" {
			theme += "Dark"
		}
	}
	params := P{"messageText": message, "type": kind, "theme": theme}
	if len(title) > 0 {
		params["titleText"] = title[0]
	}
	app.callStatic("Toast", "show", []any{params}, false, "")
}

// ToastSuccess shows a small green message for a while: ozdenui.ToastSuccess("Saved."). A title may follow.
func ToastSuccess(message string, title ...string) { showToast("success", message, title) }

// ToastError: see ToastSuccess.
func ToastError(message string, title ...string) { showToast("error", message, title) }

// ToastWarning: see ToastSuccess.
func ToastWarning(message string, title ...string) { showToast("warning", message, title) }

// ToastInfo: see ToastSuccess.
func ToastInfo(message string, title ...string) { showToast("info", message, title) }

// *** THE FILE DIALOGS OF THE SYSTEM. They return a real path, or "" when cancelled.

// FileOptions are what a file dialog starts with.
type FileOptions struct {
	Title     string
	Directory string   // the folder it opens in
	FileName  string   // SaveFile(): the name it offers
	Types     []string // the endings it shows: []string{"png", "jpg"}
}

func fileOptions(options []FileOptions) FileOptions {
	if len(options) > 0 {
		return options[0]
	}
	return FileOptions{}
}

// OpenFile is the path of the file the user chose ("": cancelled).
func OpenFile(options ...FileOptions) string {
	paths := fileDialog("open", fileOptions(options))
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// OpenFiles is as OpenFile, with more than one file.
func OpenFiles(options ...FileOptions) []string { return fileDialog("open-many", fileOptions(options)) }

// SaveFile is the path to save to ("": cancelled). The file is not made: the program writes it.
func SaveFile(options ...FileOptions) string {
	paths := fileDialog("save", fileOptions(options))
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// ChooseFolder is the path of the folder the user chose ("": cancelled).
func ChooseFolder(options ...FileOptions) string {
	paths := fileDialog("folder", fileOptions(options))
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func fileDialog(kind string, options FileOptions) []string {
	if current != nil && goroutineID() == current.appGoroutine {
		current.flush() // what the program drew before the dialog is on the page
	}
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("osascript", "-e", appleScript(kind, options))
	case "windows":
		command = exec.Command("powershell", "-NoProfile", "-STA", "-Command", powerShell(kind, options))
	default:
		command = zenity(kind, options)
	}
	output, err := command.Output()
	if err != nil {
		if _, isExit := err.(*exec.ExitError); !isExit {
			fmt.Fprintf(os.Stderr, "ozdenui: the file dialog of the system could not open: %v\n", err)
		}
		return nil // cancelled
	}
	paths := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

func quoted(text string) string {
	return `"` + strings.ReplaceAll(strings.ReplaceAll(text, `\`, `\\`), `"`, `\"`) + `"`
}

func appleScript(kind string, options FileOptions) string {
	extra := ""
	if options.Title != "" {
		extra += " with prompt " + quoted(options.Title)
	}
	if options.Directory != "" {
		if folder, err := filepath.Abs(options.Directory); err == nil {
			extra += " default location POSIX file " + quoted(folder)
		}
	}
	types := ""
	if len(options.Types) > 0 {
		list := []string{}
		for _, ending := range options.Types {
			list = append(list, quoted(strings.TrimPrefix(ending, ".")))
		}
		types = " of type {" + strings.Join(list, ", ") + "}"
	}
	choose := ""
	switch kind {
	case "open":
		choose = "set chosen to {choose file" + types + extra + "}"
	case "open-many":
		choose = "set chosen to choose file" + types + extra + " with multiple selections allowed"
	case "save":
		if options.FileName != "" {
			extra += " default name " + quoted(options.FileName)
		}
		choose = "set chosen to {choose file name" + extra + "}"
	case "folder":
		choose = "set chosen to {choose folder" + extra + "}"
	}
	// WHY the front application: the dialog belongs to the window of the program, so it opens in front of it
	return "tell application (path to frontmost application as text)\n" + choose + "\nend tell\n" +
		"set out to \"\"\nrepeat with one in chosen\nset out to out & POSIX path of one & linefeed\nend repeat\nout"
}

func powerShell(kind string, options FileOptions) string {
	single := func(text string) string { return "'" + strings.ReplaceAll(text, "'", "''") + "'" }
	script := "Add-Type -AssemblyName System.Windows.Forms; [Console]::OutputEncoding = [Text.Encoding]::UTF8; "
	if kind == "folder" {
		script += "$d = New-Object System.Windows.Forms.FolderBrowserDialog; "
		if options.Title != "" {
			script += "$d.Description = " + single(options.Title) + "; "
		}
		if options.Directory != "" {
			script += "$d.SelectedPath = " + single(options.Directory) + "; "
		}
		return script + "if ($d.ShowDialog() -eq 'OK') { $d.SelectedPath }"
	}
	class := "OpenFileDialog"
	if kind == "save" {
		class = "SaveFileDialog"
	}
	script += "$d = New-Object System.Windows.Forms." + class + "; "
	if options.Title != "" {
		script += "$d.Title = " + single(options.Title) + "; "
	}
	if options.Directory != "" {
		script += "$d.InitialDirectory = " + single(options.Directory) + "; "
	}
	if options.FileName != "" {
		script += "$d.FileName = " + single(options.FileName) + "; "
	}
	if len(options.Types) > 0 {
		patterns := []string{}
		for _, ending := range options.Types {
			patterns = append(patterns, "*."+strings.TrimPrefix(ending, "."))
		}
		list := strings.Join(patterns, ";")
		script += "$d.Filter = " + single("Files ("+list+")|"+list+"|All files (*.*)|*.*") + "; "
	}
	if kind == "open-many" {
		script += "$d.Multiselect = $true; "
	}
	return script + "if ($d.ShowDialog() -eq 'OK') { $d.FileNames }"
}

func zenity(kind string, options FileOptions) *exec.Cmd {
	args := []string{"--file-selection"}
	if options.Title != "" {
		args = append(args, "--title="+options.Title)
	}
	switch kind {
	case "open-many":
		args = append(args, "--multiple", "--separator=\n")
	case "save":
		args = append(args, "--save", "--confirm-overwrite")
	case "folder":
		args = append(args, "--directory")
	}
	start := options.Directory
	if start != "" && !strings.HasSuffix(start, "/") {
		start += "/"
	}
	if start+options.FileName != "" {
		args = append(args, "--filename="+start+options.FileName)
	}
	if len(options.Types) > 0 {
		patterns := []string{}
		for _, ending := range options.Types {
			patterns = append(patterns, "*."+strings.TrimPrefix(ending, "."))
		}
		args = append(args, "--file-filter="+strings.Join(patterns, " "))
	}
	return exec.Command("zenity", args...)
}
