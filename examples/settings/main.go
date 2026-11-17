/* Bismillah */

// Settings on three tabs. Save shows them (also in the terminal).
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor     = "#F4F5F7"
	cardColor     = "#FFFFFF"
	cardColorDark = "#232322"
	textSoft      = "#5A6270"
	formWidth     = 460
	panelHeight   = 250
	firstName     = "Ada"
)

var qualities = []string{"Low", "Normal", "High"}

var (
	name, language, email, sms *ui.Object
	sounds, volume, saved      *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Settings", Width: 580, Height: 560})
}

// cardGround is the color of the card in the theme of the window.
func cardGround() string {
	if ui.IsDark() {
		return cardColorDark
	}
	return cardColor
}

func start(page *ui.Object) {
	if !ui.IsDark() { // a dark window has its own ground
		page.Set("color", pageColor)
	}
	ui.VGroup(ui.Props{Gap: 16, Padding: 30, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "Settings", FontSize: 24, Bold: true})
		ui.VGroup(ui.Props{Width: formWidth + 40, Height: "auto", Align: "left top", Gap: 16, Padding: 20, Color: cardGround(), Round: 12}, func() {
			tabs := ui.Tabs(ui.TabsProps{Tabs: []string{"General", "Notifications", "Sound"}, PanelHeight: panelHeight})
			tabs.Panel("General", generalPanel)
			tabs.Panel("Notifications", notificationsPanel)
			tabs.Panel("Sound", soundPanel)
			ui.HGroup(ui.Props{Hug: true, Gap: 10}, func() {
				ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan", OnClick: save})
				ui.SmartButton(ui.SmartButtonProps{LabelText: "Reset", OnClick: reset})
			})
		})
		saved = ui.Label(ui.Props{FontSize: 13, TextColor: textSoft})
	})
}

func panelGroup(gap int, content func()) {
	ui.VGroup(ui.Props{Width: "100%", Height: "auto", Align: "left top", Gap: gap, Padding: []int{0, 16}}, content)
}

func generalPanel() {
	panelGroup(14, func() {
		name = ui.TextField(ui.TextFieldProps{TitleText: "Name", Value: firstName, Width: formWidth})
		ui.Label(ui.Props{Text: "Language", FontSize: 13, TextColor: textSoft})
		language = ui.SelectBox(ui.SelectBoxProps{Options: []string{"English", "Türkçe", "Deutsch"}, Value: "English", Width: formWidth})
	})
}

func notificationsPanel() {
	panelGroup(16, func() {
		email = ui.CheckBox(ui.CheckBoxProps{LabelText: "E-mail", Checked: true})
		sms = ui.CheckBox(ui.CheckBoxProps{LabelText: "SMS"})
		ui.HGroup(ui.Props{Hug: true, Gap: 12, Align: "left center"}, func() {
			sounds = ui.Toggle(ui.ToggleProps{Value: true})
			ui.Label(ui.Props{Text: "Play a sound", FontSize: 15})
		})
	})
}

func soundPanel() {
	panelGroup(16, func() {
		ui.Label(ui.Props{Text: "Volume", FontSize: 13, TextColor: textSoft})
		volume = ui.Slider(ui.SliderProps{Value: 7, MinValue: 0, MaxValue: 10, Width: formWidth, StyleName: "round"})
		ui.HGroup(ui.Props{Hug: true, Gap: 18}, func() {
			for _, quality := range qualities {
				ui.RadioButton(ui.RadioButtonProps{Group: "quality", Value: quality, LabelText: quality, Checked: quality == "Normal"})
			}
		})
	})
}

func values() string {
	return fmt.Sprintf("name: %s, language: %v, email: %t, sms: %t, sounds: %t, volume: %v, quality: %v",
		name.Text(), language.Value(), email.Checked(), sms.Checked(), sounds.Checked(), volume.Value(), ui.RadioValue("quality"))
}

func save() {
	settings := values()
	fmt.Println("saved:", settings)
	saved.SetText(settings)
	ui.ToastSuccess("The settings are saved.")
}

func reset() {
	name.SetText(firstName)
	language.SetValue("English")
	email.Set("checked", true)
	sms.Set("checked", false)
	sounds.SetValue(true)
	volume.SetValue(7)
	saved.SetText("")
}
