/* Bismillah */

/*

doc.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The page of the module in `go doc`.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

// Package ozdenui makes desktop windows in plain Go: the objects of basic.js and the components of the OzdenUI
// kit (buttons, fields, lists, tables, charts, dialogs...), with no HTML, CSS or JavaScript to write.
//
//	package main
//
//	import ui "github.com/bug7a/ozdenui-go"
//
//	func main() {
//	    ui.Run(start, ui.Window{Title: "Hello", Width: 400, Height: 300})
//	}
//
//	func start(page *ui.Object) {
//	    ui.VGroup(ui.Props{Gap: 16}, func() {
//	        label := ui.Label(ui.Props{Text: "Hello", FontSize: 32})
//	        ui.SmartButton(ui.SmartButtonProps{LabelText: "Press", Tone: "cyan", OnClick: func() {
//	            label.SetText("Pressed")
//	        }})
//	    })
//	}
//
// # How it works
//
// The window shows a page of basic.js that is inside the program (the built program is one file). An object in Go
// is a proxy: it keeps an id and the values it was given, and sends commands to the page, where the real object
// lives; the events of the page come back as calls of the functions of the program. Everything of the program runs
// on one goroutine, one task at a time (start, an event, a timer), and the commands of a task go to the page
// together.
//
// # The parts
//
//   - Run opens the window and calls start(page); Window says how it opens (its size, Theme, Style).
//   - Box, Label, Button, Input, Icon and the groups VGroup, HGroup take Props.
//   - Every component of the kit is a function with its own Props type: CheckBox(CheckBoxProps{...}),
//     SmartTable(SmartTableProps{...})... Their fields are the params of the kit, each with its default in its
//     comment. Component() makes one by its name.
//   - An Object is changed with Set (by the names of the kit), read with Get, Text, Value, Checked, and called with
//     Call and Ask. Add puts objects into a box later.
//   - Alert, Confirm, Prompt and Ask are dialogs that return their answer; ToastSuccess and its kin never wait;
//     OpenFile, SaveFile and ChooseFolder are the file dialogs of the system.
//   - After, Every and Do are the timers and the way back from a goroutine of the program.
//
// # A picture of the window
//
// OZDENUI_LOOK=app.png go run . runs the program without a window and saves a picture of it (it needs Chrome, Edge
// or Chromium): the way to see a window from a terminal, for a check or for an AI assistant.
//
// The rules for an AI assistant that writes a program with this module:
// go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents
package ozdenui
