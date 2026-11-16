/* Bismillah */

// A place to look at a picture: open a file, or draw a plot again and again: the zoom and the place stay.
// The plot is an image.Image of the standard library: ui.Picture() takes it as it is (also bytes of a PNG or a JPEG).
package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	plotWidth  = 640
	plotHeight = 400
	pageLight  = "#FFFFFF"
	pageDark   = "#1E1E1E"
	curveWidth = 3
	phaseStep  = 0.5
	gridStep   = 40
)

var (
	ground = color.RGBA{255, 255, 255, 255}
	grid   = color.RGBA{225, 228, 232, 255}
	curve  = color.RGBA{43, 140, 150, 255}
)

var (
	phase  = 0.0
	viewer *ui.Object
	status *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Image viewer", Width: 760, Height: 560})
}

// plot draws a sine pixel by pixel: as a program of yours would draw one.
func plot(phase float64) image.Image {
	picture := image.NewRGBA(image.Rect(0, 0, plotWidth, plotHeight))
	for y := 0; y < plotHeight; y++ {
		for x := 0; x < plotWidth; x++ {
			picture.Set(x, y, ground)
			if x%gridStep == 0 || y%gridStep == 0 {
				picture.Set(x, y, grid)
			}
		}
	}
	for x := 0; x < plotWidth; x++ {
		y := plotHeight/2 - int(math.Sin(float64(x)/50+phase)*plotHeight*0.35)
		for thickness := 0; thickness < curveWidth; thickness++ {
			picture.Set(x, y+thickness, curve) // Set leaves out a point outside the picture
		}
	}
	return picture
}

func start(page *ui.Object) {
	page.Set("color", map[bool]string{false: pageLight, true: pageDark}[ui.IsDark()])
	ui.VGroup(ui.Props{Gap: 0, Align: "left top"}, func() {
		ui.HGroup(ui.Props{Width: "100%", Height: "auto", Gap: 8, Padding: 12, Align: "left center", Shrink: 0}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Open a picture", Tone: "cyan", OnClick: openPicture})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Draw the plot again", OnClick: drawAgain})
		})
		viewer = ui.ImageView(ui.ImageViewProps{Fill: true, Unframed: true, ImageFile: ui.Picture(plot(phase)),
			AltText: "A plot", OnZoomChange: showSize,
			OnError: func() { status.Call("setMessageText", "It can not be opened", "error") }})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "A plot", Items: []ui.P{{"key": "size", "text": "-"}}})
	})
}

func showSize(view *ui.Object) {
	status.Call("setItemText", "size", fmt.Sprintf("%d x %d at %d%%", view.Int("naturalWidth"), view.Int("naturalHeight"),
		int(math.Round(view.Num("zoomLevel")*100))))
}

func drawAgain() {
	phase += phaseStep
	viewer.Call("setImage", ui.Picture(plot(phase)))
	status.Call("setMessageText", "Drawn again: the zoom and the place stay", "success", 2000)
}

func openPicture() {
	path := ui.OpenFile(ui.FileOptions{Types: []string{"png", "jpg", "jpeg", "gif", "webp", "svg"}})
	if path == "" {
		return
	}
	viewer.Call("setImage", ui.Picture(path))
	status.Call("setMessageText", path)
}
