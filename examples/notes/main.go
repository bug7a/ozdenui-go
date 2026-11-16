/* Bismillah */

// Notes: the program of the rules for an AI assistant (cmd/ozdenui/AGENTS.md).
package main

import (
	"fmt"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const pageColor = "#F4F5F7"

var (
	notes []string   // the data of the program is plain Go
	name  *ui.Object // the objects the program changes later
	list  *ui.Object
	empty *ui.Object
	count *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Notes", Width: 480, Height: 420}) // once, at the end of main()
}

func start(page *ui.Object) { // runs when the window is ready: make the objects here
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 0}, func() { // a group is as large as its box: here the whole window
		ui.HGroup(ui.Props{Height: 60, Padding: []int{16, 0}, Gap: 10, Align: "left center", Color: "white"}, func() {
			name = ui.TextField(ui.TextFieldProps{Placeholder: "A new note", Width: 260, OnEnter: addNote})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Add", Tone: "cyan", OnClick: addNote})
		})
		list = ui.ScrollArea(ui.ScrollAreaProps{Padding: 16, Gap: 10, Align: "left top"}, func() { // takes the room that is left
			empty = ui.Label(ui.Props{Text: "No notes yet", FontSize: 15})
		})
		count = ui.Label(ui.Props{Text: "Notes: 0", FontSize: 13, Height: 28})
	})
}

func addNote() {
	typed := strings.TrimSpace(name.Text()) // the state is in the program: reading asks nothing of the window
	if typed == "" {
		ui.ToastWarning("Write a note first.")
		return
	}
	notes = append(notes, typed)
	list.Add(func() { // objects made later go into a box the same way
		ui.Label(ui.Props{Text: typed, FontSize: 15})
	})
	empty.Hide()
	name.SetValue("")
	count.SetText(fmt.Sprintf("Notes: %d", len(notes)))
}
