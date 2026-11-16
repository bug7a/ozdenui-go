/* Bismillah */

// A to-do list: add a task (Enter or the button), tick the done ones, delete them (with a question).
package main

import (
	"fmt"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor = "#F4F5F7"
	textSoft  = "#5A6270"
	listWidth = 440
)

var firstTasks = []string{"Buy bread", "Call the bank", "Water the plants"}

var (
	nextNumber = 0
	total      = 0 // the tasks in the list: the program counts them, so it asks the page nothing
	field      *ui.Object
	tasks      *ui.Object
	count      *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "To do", Width: 540, Height: 560})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	items := []ui.P{}
	for _, title := range firstTasks {
		items = append(items, newItem(title))
	}
	total = len(items)
	ui.VGroup(ui.Props{Gap: 14, Padding: 36, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "To do", FontSize: 24, Bold: true})
		ui.HGroup(ui.Props{Hug: true, Gap: 8, Align: "left bottom"}, func() {
			field = ui.TextField(ui.TextFieldProps{TitleText: "A new task", Placeholder: "Buy milk", Width: listWidth - 108, OnEnter: addTask})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Add", Width: 100, Tone: "cyan", OnClick: addTask})
		})
		tasks = ui.ItemList(ui.ItemListProps{Items: items, Selection: "multiple", Width: listWidth, Height: 300,
			EmptyText: "Nothing to do.", OnChange: showCount})
		ui.HGroup(ui.Props{Hug: true, Gap: 14, Align: "left center"}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Delete the done ones", Tone: "orange", OnClick: deleteDone})
			count = ui.Label(ui.Props{FontSize: 14, TextColor: textSoft})
		})
	})
	showCount()
}

func newItem(title string) ui.P {
	nextNumber++
	return ui.P{"key": fmt.Sprintf("t%d", nextNumber), "text": title}
}

func addTask() {
	title := strings.TrimSpace(field.Text())
	if title == "" {
		return
	}
	tasks.Call("addItem", newItem(title))
	total++
	field.SetValue("")
	showCount()
}

func doneKeys() []string {
	keys := []string{}
	selected, _ := tasks.Get("selectedKeys").([]any)
	for _, key := range selected {
		keys = append(keys, fmt.Sprint(key))
	}
	return keys
}

func deleteDone() {
	done := doneKeys()
	if len(done) == 0 {
		ui.Alert("Tick the done tasks first.", ui.DialogOptions{Title: "Nothing to delete"})
		return
	}
	question := fmt.Sprintf("%d done task(s) will be deleted.", len(done))
	if !ui.Confirm(question, ui.DialogOptions{Title: "Delete them?", ConfirmText: "Delete", Danger: true}) {
		return
	}
	for _, key := range done {
		tasks.Call("removeItem", key)
	}
	total -= len(done)
	showCount()
}

func showCount() {
	count.SetText(fmt.Sprintf("%d tasks, %d done", total, len(doneKeys())))
}
