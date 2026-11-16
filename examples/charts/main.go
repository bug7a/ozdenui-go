/* Bismillah */

// Charts with no other library: a live line, bars with a format, a doughnut, two lines with a target. The page draws them.
package main

import (
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor   = "#F4F5F7"
	textSoft    = "#5A6270"
	chartWidth  = 410
	chartHeight = 250
	livePoints  = 24
	tickMs      = 1000
)

var (
	months  = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	ticks   = 0
	isTable = false
	live    *ui.Object
	sales   *ui.Object
	chosen  *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Charts", Width: 920, Height: 720})
}

func liveValue(tick int) float64 {
	return math.Round((20+6*math.Sin(float64(tick)/3))*10) / 10
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	labels := []string{}
	values := []float64{}
	for tick := 0; tick < livePoints; tick++ {
		labels = append(labels, strconv.Itoa(tick))
		values = append(values, liveValue(tick))
	}
	ui.ScrollArea(ui.ScrollAreaProps{Align: "center top", Gap: 16, Padding: 24}, func() {
		ui.Label(ui.Props{Text: "Charts", FontSize: 24, Bold: true})
		ui.HGroup(ui.Props{Width: "100%", Height: "auto", Align: "center top", Gap: 16, Wrap: true}, func() {
			live = ui.GraphBox(ui.GraphBoxProps{Type: "area", TitleText: "Temperature", SubtitleText: "A new value every second",
				Labels: labels, Values: values, Curve: "smooth", ShowPoints: false, MaxCount: livePoints, MinValue: 0,
				MaxLabels: 6, FormatValue: "{:.1f}°", Width: chartWidth, Height: chartHeight})
			sales = ui.GraphBox(ui.GraphBoxProps{Type: "bar", TitleText: "Sales", Labels: months,
				Values: []int{1200, 1900, 700, 1500, 2300, 1800}, FormatValue: "{:,.0f} TL", OnClick: showMonth,
				Width: chartWidth, Height: chartHeight})
			ui.GraphBox(ui.GraphBoxProps{Type: "doughnut", TitleText: "Where they came from",
				Labels: []string{"Search", "A friend", "An ad"}, Values: []int{55, 30, 15}, CenterText: "100",
				CenterSubText: "people", FormatValue: "{}%", Width: chartWidth, Height: chartHeight})
			ui.GraphBox(ui.GraphBoxProps{Type: "line", TitleText: "Two years", Labels: months,
				Datasets: []ui.P{
					{"label": "2025", "data": []int{8, 11, 9, 14, 13, 17}},
					{"label": "2026", "data": []int{10, 14, 12, 18, 21, 23}, "showPoints": true},
				},
				TargetLines: []ui.P{{"value": 20, "label": "Target"}}, Width: chartWidth, Height: chartHeight})
		})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Line", OnClick: func() { sales.Call("setType", "line") }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Bars", OnClick: func() { sales.Call("setType", "bar") }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "As a table", OnClick: toggleTable})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Save a picture", Tone: "cyan", OnClick: savePicture})
		})
		chosen = ui.Label(ui.Props{Text: "Click a bar.", FontSize: 14, TextColor: textSoft})
	})
	ui.Every(tickMs, tick)
}

func tick() {
	ticks++
	number := livePoints + ticks
	live.Call("addValue", strconv.Itoa(number), liveValue(number))
}

func showMonth(graph *ui.Object, item ui.P) {
	value, _ := item["value"].(float64)
	chosen.SetText(fmt.Sprintf("%v: %s TL", item["label"], thousands(int(value))))
}

// thousands writes 1900 as "1,900".
func thousands(number int) string {
	digits := strconv.Itoa(number)
	var out strings.Builder
	for index, digit := range digits {
		if index > 0 && (len(digits)-index)%3 == 0 {
			out.WriteByte(',')
		}
		out.WriteRune(digit)
	}
	return out.String()
}

func toggleTable() {
	isTable = !isTable
	sales.Call("setTableView", isTable)
}

func savePicture() {
	path := ui.SaveFile(ui.FileOptions{FileName: "sales.png", Types: []string{"png"}})
	if path == "" {
		return
	}
	// getImage gives a data address: "data:image/png;base64,...."
	address, _ := sales.Ask("getImage", "image/png", 2).(string)
	_, encoded, found := strings.Cut(address, ",")
	picture, err := base64.StdEncoding.DecodeString(encoded)
	if !found || err != nil {
		ui.ToastError("The page gave no picture.")
		return
	}
	if err := os.WriteFile(path, picture, 0o644); err != nil {
		ui.ToastError(err.Error(), "Could not save it")
		return
	}
	ui.ToastSuccess("Saved.")
}
