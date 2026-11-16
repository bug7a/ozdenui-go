# 3. Layout

A place of its own (`Left: 20, Top: 20`) is good for a game or a drawing. A window of a program is easier with
**groups**: they place their objects, and place them again when the window changes its size.

## VGroup and HGroup

```go
ui.VGroup(ui.Props{Gap: 10, Padding: 20, Align: "left top"}, func() {
	ui.Label(ui.Props{Text: "One"})
	ui.Label(ui.Props{Text: "Two"}) // under "One", 10 px lower
	ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Yes"})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "No"}) // next to "Yes"
	})
})
```

- `VGroup` puts its objects one under the other, `HGroup` side by side.
- A group is as large as its box: `Width: "100%"`, `Height: "100%"`. `Hug: true` makes it only as large as its objects.
- `Gap`: the space between the objects. `Padding`: the space inside its edge: `20`, or `[]int{x, y}` (`[]int{30, 12}`).
- `Align`: where the objects stand: `"center"` (the default), `"left top"`, `"right center"`, `"center bottom"`...
- `Wrap: true`: the objects go on to the next line when there is no room (cards, tags).
- A group is a `Box`: `Color`, `Round`, `Border` work on it too.

In a group an object needs only its size: `ui.Label(ui.Props{Text: "One", Width: 200})`. Without a size it is as
large as its content.

**A group inside a group** almost always needs one of these: `Hug: true` (as large as its objects), a height
(`Height: "auto"` or a number), or `Fill: true` (below). Two groups of 100% in one `VGroup` squeeze each other.

**More than fits:** a group with more than its box can hold squeezes its objects, and one can cover another (a title
under a card). Give the window more room, or put the content in a `ScrollArea` (below).

## Taking the room that is left: Fill

A typical window: a bar on top, a part in the middle that takes all the room, a line at the bottom.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

func main() {
	ui.Run(start, ui.Window{Title: "Notes", Width: 640, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 0}, func() {
		ui.HGroup(ui.Props{Height: 52, Padding: []int{16, 0}, Gap: 10, Align: "left center", Color: "white"}, func() {
			ui.Label(ui.Props{Text: "Notes", FontSize: 18, Bold: true})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "New", Tone: "cyan"})
		})
		ui.HGroup(ui.Props{Gap: 0, Fill: true}, func() {
			ui.VGroup(ui.Props{Width: 180, Height: "100%", Padding: 16, Gap: 8, Align: "left top", Color: "#E9EBEF"}, func() {
				for _, folder := range []string{"All", "Work", "Home"} {
					ui.Label(ui.Props{Text: folder, FontSize: 15})
				}
			})
			ui.ScrollArea(ui.ScrollAreaProps{Padding: 16, Gap: 12}, func() {
				for number := 1; number <= 40; number++ {
					ui.Label(ui.Props{Text: fmt.Sprintf("Note %d", number), FontSize: 15})
				}
			})
		})
		ui.Label(ui.Props{Text: "40 notes", FontSize: 13, TextColor: "#5A6270", Height: 28})
	})
}
```

- `Fill: true` gives an object the room the others leave: **the height** that is left in a `VGroup` (the middle row
  between the bar and the bottom line), **the width** that is left in an `HGroup`. Any object and any component
  takes it: `ui.HGroup(ui.Props{Fill: true})`, `ui.LogViewProps{Fill: true}`, `ui.SplitterProps{Fill: true}`.
- `Fill` is given when the object is made; it can not be changed with `Set()` later.
- `ui.ScrollArea` takes the room the others leave by itself (its `Width` and `Height` are `"100%"` of what is left)
  and scrolls what is in it. As an object of the page (not in a group) it is the whole window.
- Its objects go one under the other, as in a `VGroup`, from the top left: give it `Align`, `Gap` and `Padding`
  (the padding scrolls with them). `ScrollX: true` scrolls to the side too. It has the scroll bar of the kit, light or
  dark with the window.

## GridGroup: columns and rows

```go
package main

import ui "github.com/bug7a/ozdenui-go"

func main() {
	ui.Run(start, ui.Window{Title: "Keys", Width: 320, Height: 300})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Padding: 16, Align: "left top"}, func() {
		ui.GridGroup(ui.GridGroupProps{Columns: 4, Gap: 8}, func() {
			for _, key := range "789/456*123-0.=+" {
				ui.SmartButton(ui.SmartButtonProps{LabelText: string(key), Width: "100%"})
			}
			clearKey := ui.SmartButton(ui.SmartButtonProps{LabelText: "Clear", Width: "100%", Tone: "orange"})
			ui.Current().Place(clearKey, ui.P{"columnSpan": 4})
		})
	})
}
```

- The objects go into the cells in order, left to right, row by row.
- `Columns: 4`: equal columns. `Columns: []any{120, "1fr", "auto"}`: 120 px, the rest, as wide as its content.
- An object fills its cell; `Align: "left center"` on the grid puts each one in its cell instead.
- `grid.Place(obj, ui.P{"columnSpan": 2})`, `grid.Place(obj, ui.P{"column": 1, "row": 0})`: a cell of its own
  (counted from 0).
- **Inside the function of a group, its own variable is not set yet** (`grid := ui.GridGroup(..., func() { grid... })`
  is `nil` there): the group is `ui.Current()`.
- Unlike the other groups, a grid is as high as its rows (`Height: "auto"`). `Padding`, `Color` and `Round` are
  fields of its own, as in a group. `examples/calculator` is a whole calculator.

## Things that go out of a group

A group does not cut what goes out of its edge: a button's focus ring, a shadow or a button that grows under the mouse
shows whole, also in a small group:

```go
ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
	ui.SmartButton(ui.SmartButtonProps{LabelText: "Save", Tone: "cyan"})
	ui.SmartButton(ui.SmartButtonProps{LabelText: "Cancel"})
})
```

`ClipContent: true` cuts, where cutting is the point (a picture in a frame). A group that scrolls (`ScrollY: true`)
and a `ScrollArea` cut by themselves.

Two defaults are not the ones of basic.js, so that a program says less:

- A `VGroup` / `HGroup` of basic.js cuts its content; in Go it does not.
- `MenuBar`, `ToolBar`, `StatusBar`, `SearchInput` and `FilePicker` do not shrink in a group (`Shrink: 0`): they keep
  their height while a `ScrollArea` or a `Fill: true` object takes the rest.

## A place of its own

Outside a group, in a `Box` or on the page, an object stands where you put it: `Left`, `Top`, `Width`, `Height`. It
is the way for a game, a drawing or an overlay:

```go
ui.Box(ui.Props{Width: 400, Height: 300, Color: "#1E2230"}, func() {
	player = ui.Box(ui.Props{Left: 180, Top: 130, Width: 40, Height: 40, Color: "#5BBAC1", Round: 8})
})
player.Set("left", 200) // moves it to 200 px from the left of its box
```

Keep the place in a variable of your own (`playerLeft += 20`, then `player.Set("left", playerLeft)`): reading it
back is a question to the page. `examples/game` moves a box with the arrow keys this way. A size may be a text:
`Width: "100%"` fills its box.

## When the window changes its size

Groups place their objects again by themselves. For your own sums, `page.Width()` and `page.Height()` are always
current, and `page.On("resize", fn)` tells you:

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var size *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "The size", Width: 420, Height: 240})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	size = ui.Label(ui.Props{Left: 20, Top: 60, FontSize: 28, Bold: true})
	ui.Label(ui.Props{Left: 20, Top: 20, Text: "Change the size of the window.", FontSize: 15, TextColor: "#5A6270"})
	showSize(page)
	page.On("resize", showSize)
}

func showSize(page *ui.Object) {
	size.SetText(fmt.Sprintf("%.0f x %.0f", page.Width(), page.Height()))
}
```

`page.On("resize", fn)` takes the same kinds of function as every event (chapter 4): here one that takes the page.

## Next

[Chapter 4: Events and time](04-events-and-time.md): clicks, keys, timers and long jobs.
