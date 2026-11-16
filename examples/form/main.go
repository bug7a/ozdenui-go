/* Bismillah */

// A form: fields of text, a date, a time, a color, a number. Send checks them and shows their values.
package main

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor  = "#F4F5F7"
	cardColor  = "#FFFFFF"
	textSoft   = "#5A6270"
	formWidth  = 520
	firstColor = "#3871E0"
	dateLayout = "2006-01-02" // the dates of the kit: "YYYY-MM-DD"
)

var (
	name, email, password, budget *ui.Object
	date, clock, color            *ui.Object
	guests, noise, result         *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A form", Width: 620, Height: 720})
}

// titled puts a small title over a field that has none of its own: titled("Date", func() { ui.DatePicker(...) }).
func titled(title string, content func()) *ui.Object {
	return ui.VGroup(ui.Props{Hug: true, Gap: 6, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: title, FontSize: 13, TextColor: textSoft})
		content()
	})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	today := time.Now()
	nextWeek := today.AddDate(0, 0, 7)
	ui.ScrollArea(ui.ScrollAreaProps{Align: "center top", Gap: 16, Padding: 30}, func() {
		ui.Label(ui.Props{Text: "A reservation", FontSize: 24, Bold: true})
		ui.VGroup(ui.Props{Width: formWidth, Height: "auto", Align: "left top", Gap: 14, Padding: 20, Color: cardColor, Round: 12}, func() {
			name = ui.TextField(ui.TextFieldProps{TitleText: "Name", Placeholder: "Ada", Required: true, Width: "100%", OnEnter: send})
			email = ui.EmailField(ui.EmailFieldProps{TitleText: "E-mail", Required: true, Width: "100%", OnEnter: send})
			password = ui.PasswordField(ui.PasswordFieldProps{TitleText: "Password", Required: false, MinLength: 6, Width: "100%"})
			budget = ui.NumberField(ui.NumberFieldProps{TitleText: "Budget", Placeholder: ui.Empty, Value: 250, Width: "100%"})
			ui.HGroup(ui.Props{Hug: true, Gap: 12, Wrap: true, Align: "left top"}, func() {
				titled("Date", func() {
					date = ui.DatePicker(ui.DatePickerProps{Date: nextWeek.Format(dateLayout), MinDate: today.Format(dateLayout), Width: 150})
				})
				titled("Time", func() {
					clock = ui.TimePicker(ui.TimePickerProps{Time: "19:30", MinuteStep: 15, Width: 120})
				})
				titled("Table color", func() {
					color = ui.ColorPicker(ui.ColorPickerProps{SelectedColor: firstColor, Width: 150})
				})
			})
			titled("Guests", func() {
				guests = ui.Stepper(ui.StepperProps{Value: 2, MinValue: 1, MaxValue: 12})
			})
			titled("How quiet", func() {
				noise = ui.Slider(ui.SliderProps{Value: 70, Width: formWidth - 40, StyleName: "round", OnChangeEnd: showNoise})
			})
			ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
				ui.SmartButton(ui.SmartButtonProps{LabelText: "Send", Tone: "cyan", OnClick: send})
				ui.SmartButton(ui.SmartButtonProps{LabelText: "Clear", OnClick: clear})
			})
		})
		result = ui.Label(ui.Props{FontSize: 13, TextColor: textSoft, Width: formWidth})
	})
}

func values() string {
	parts := []string{
		"name: " + name.Text(),
		"email: " + email.Text(),
		"budget: " + budget.Text(),
		// WHY Measure: Value() of a picker is "" until the user changes it, also when it was made with a date
		fmt.Sprint("date: ", date.Measure("value")),
		fmt.Sprint("time: ", clock.Measure("value")),
		fmt.Sprint("color: ", color.Value()),
		fmt.Sprint("guests: ", guests.Value()),
		fmt.Sprint("quiet: ", noise.Value()),
	}
	return strings.Join(parts, ", ")
}

func showNoise(slider *ui.Object) {
	result.SetText(fmt.Sprint("How quiet: ", slider.Value()))
}

// isRight shows the warning of an empty required field or a wrong value, and tells whether there is none.
func isRight(field *ui.Object) bool {
	field.Call("checkInputValue")
	return field.Measure("status") == 0.0 // 0: no warning, 1: empty but required, 2: not valid
}

func send() {
	// WHY every field and not "&&": each one shows its warning, not only the first one
	answers := []bool{isRight(name), isRight(email), isRight(password)}
	for _, right := range answers {
		if !right {
			ui.ToastError("Look at the marked fields.")
			return
		}
	}
	result.SetText(values())
	ui.ToastSuccess("Sent.")
}

func clear() {
	for _, field := range []*ui.Object{name, email, password, budget} {
		field.SetText("")
	}
	date.Set("date", "")
	result.SetText("")
}
