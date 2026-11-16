/* Bismillah */

// A CSV file in a table: open one (a sample is shown first), filter it, click a row to see it. A SmartTable makes
// only the rows that are seen: "100000 rows" shows that a big file stays fast.
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor   = "#F4F5F7"
	textSoft    = "#5A6270"
	tableWidth  = 820
	tableHeight = 400
	columnWidth = 170
	manyRows    = 100000
)

var sampleColumns = []string{"city", "country", "people"}

var sample = []ui.P{
	{"city": "İstanbul", "country": "Türkiye", "people": "15,655,924"},
	{"city": "Ankara", "country": "Türkiye", "people": "5,803,482"},
	{"city": "Berlin", "country": "Germany", "people": "3,878,100"},
	{"city": "Madrid", "country": "Spain", "people": "3,332,035"},
	{"city": "Rome", "country": "Italy", "people": "2,748,109"},
	{"city": "Paris", "country": "France", "people": "2,102,650"},
	{"city": "Vienna", "country": "Austria", "people": "2,005,760"},
	{"city": "Lisbon", "country": "Portugal", "people": "567,131"},
}

var (
	columns []string // the names of the columns in their order: a map of a row has no order
	file    *ui.Object
	holder  *ui.Object
	table   *ui.Object
	detail  *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "CSV viewer", Width: 900, Height: 640})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 14, Padding: 30, Align: "left top"}, func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 14, Align: "left center"}, func() {
			ui.Label(ui.Props{Text: "CSV viewer", FontSize: 24, Bold: true})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Open a CSV", Tone: "cyan", OnClick: openCsv})
			ui.SmartButton(ui.SmartButtonProps{LabelText: fmt.Sprintf("%d rows", manyRows), OnClick: showManyRows})
		})
		file = ui.Label(ui.Props{Text: "A sample: some cities", FontSize: 14, TextColor: textSoft})
		holder = ui.Box(ui.Props{Width: tableWidth, Height: tableHeight, Color: "transparent"})
		detail = ui.Label(ui.Props{Text: "Click a row.", FontSize: 15})
	})
	showRows(sampleColumns, sample)
}

// showRows makes a new table for the rows (a new file has other columns).
func showRows(names []string, records []ui.P) {
	if table != nil {
		table.Remove()
	}
	columns = names
	titles := []ui.P{}
	for _, name := range names {
		titles = append(titles, ui.P{"name": name, "dataTitle": name, "dataType": "string", "width": columnWidth, "sortable": 1})
	}
	holder.Add(func() {
		table = ui.SmartTable(ui.SmartTableProps{TitleDataList: titles, ItemDataList: records, FitLineCount: 1,
			Width: "100%", Height: "100%", OnSelect: showRow})
	})
}

func showRow(table *ui.Object, row ui.P) {
	parts := []string{}
	for _, name := range columns {
		parts = append(parts, fmt.Sprintf("%s: %v", name, row[name]))
	}
	detail.SetText(strings.Join(parts, "   "))
}

func showManyRows() {
	file.SetText(fmt.Sprintf("%d rows made in Go", manyRows))
	records := make([]ui.P, manyRows)
	for index := range records {
		records[index] = ui.P{"number": index, "square": index * index, "half": float64(index) / 2}
	}
	showRows([]string{"number", "square", "half"}, records)
}

func openCsv() {
	path := ui.OpenFile(ui.FileOptions{Types: []string{"csv"}})
	if path == "" {
		return
	}
	names, records, err := readCsv(path)
	if err != nil {
		ui.ToastError(err.Error(), "Could not read it")
		return
	}
	file.SetText(path)
	showRows(names, records)
}

// readCsv reads a CSV file with a line of titles: a row is a map by those titles.
func readCsv(path string) ([]string, []ui.P, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	text := strings.TrimPrefix(string(content), "\xef\xbb\xbf") // the mark Excel puts at the start of a UTF-8 file
	reader := csv.NewReader(strings.NewReader(text))
	reader.FieldsPerRecord = -1 // a short or a long line is read as it is
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(lines) == 0 {
		return nil, []ui.P{}, nil
	}
	names := lines[0]
	records := []ui.P{}
	for _, line := range lines[1:] {
		record := ui.P{}
		for index, name := range names {
			if index < len(line) {
				record[name] = line[index]
			}
		}
		records = append(records, record)
	}
	return names, records, nil
}
