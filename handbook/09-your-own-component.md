# 9. Your own component

A part you need more than once (a card, a row of a form, a rating) becomes a type of yours: a struct with a
constructor that makes its objects, keeps the ones it changes, and has methods. No JavaScript is needed, and nothing
to register: it is plain Go.

## A card

```go
package main

import ui "github.com/bug7a/ozdenui-go"

// Card is a component of this program: a title, a text and a button.
type Card struct {
	Box   *ui.Object // the group that holds the card
	title *ui.Object
}

// NewCard makes a card in the group it is called in. onOpen runs when its button is clicked.
func NewCard(title, text string, onOpen func(card *Card)) *Card {
	card := &Card{}
	card.Box = ui.VGroup(ui.Props{Width: 220, Height: "auto", Align: "left top",
		Padding: 16, Gap: 8, Color: "white", Round: 12}, func() {
		card.title = ui.Label(ui.Props{Text: title, FontSize: 18, Bold: true})
		ui.Label(ui.Props{Text: text, FontSize: 15, TextColor: "#5A6270"})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Open", OnClick: func() {
			if onOpen != nil {
				onOpen(card)
			}
		}})
	})
	return card
}

// SetTitle changes its title.
func (card *Card) SetTitle(title string) {
	card.title.SetText(title)
}

func main() {
	ui.Run(start, ui.Window{Title: "Cards", Width: 760, Height: 240})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.HGroup(ui.Props{Gap: 16, Padding: 24, Wrap: true}, func() {
		for _, name := range []string{"Sales", "Stock", "Orders"} {
			NewCard(name, "Updated today.", func(card *Card) { card.SetTitle("Opened") })
		}
	})
}
```

- **The constructor makes the objects where it is called**, as a component of the kit does: inside the function of a
  group, the card goes into that group.
- **The struct keeps its objects** (`Box`, `title`). An exported field (`Box`) is for the program that uses the card;
  a small one (`title`) stays inside.
- **Methods change it** (`SetTitle()`), as the methods of the components do.
- **Its objects are ordinary objects:** `card.Box.Remove()` removes the card and everything in it;
  `card.Box.On("mouseenter", fn)` and `card.title.Measure("width")` work as on any object.
- A struct may embed the object instead (`type Card struct { *ui.Object; title *ui.Object }`, then
  `card.Object = ui.VGroup(...)`): the card has every method of an object, `card.Hide()`, `card.Remove()`.

A part that keeps nothing needs no type: a function that makes its objects and returns the `*ui.Object` is enough.

```go
// caption is the small title over a part of a page.
func caption(text string) *ui.Object {
	return ui.Label(ui.Props{Text: strings.ToUpper(text), FontSize: 11, Bold: true, TextColor: "#5A6270"})
}
```

## State and a callback of its own

A component of yours keeps its state in Go and tells the program when the user changes it:

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	starOn  = "#E0A038"
	starOff = "#C9CDD4"
)

// Rating is five stars: Value (0 to 5); OnChange runs when the user clicks a star.
type Rating struct {
	Value    int
	OnChange func(rating *Rating)
	Box      *ui.Object
	stars    []*ui.Object
}

func NewRating(value int, onChange func(rating *Rating)) *Rating {
	rating := &Rating{Value: value, OnChange: onChange}
	rating.Box = ui.HGroup(ui.Props{Hug: true, Gap: 4}, func() {
		for number := 1; number <= 5; number++ {
			star := ui.Label(ui.Props{Text: "★", FontSize: 28, OnClick: func() { rating.click(number) }})
			rating.stars = append(rating.stars, star)
		}
	})
	rating.show()
	return rating
}

// SetValue is for the code: it does not call OnChange, as the components do.
func (rating *Rating) SetValue(value int) {
	rating.Value = value
	rating.show()
}

func (rating *Rating) click(number int) {
	rating.SetValue(number)
	if rating.OnChange != nil {
		rating.OnChange(rating)
	}
}

func (rating *Rating) show() {
	for index, star := range rating.stars {
		color := starOff
		if index < rating.Value {
			color = starOn
		}
		star.Set("textColor", color)
	}
}

func main() {
	ui.Run(start, ui.Window{Title: "Rating", Width: 320, Height: 200})
}

func start(page *ui.Object) {
	var said *ui.Object
	ui.VGroup(ui.Props{Gap: 12}, func() {
		said = ui.Label(ui.Props{Text: "Rate it", FontSize: 16})
		NewRating(3, func(rating *Rating) { said.SetText(fmt.Sprintf("%d stars", rating.Value)) })
	})
}
```

`func() { rating.click(number) }` gives each star its own number: since Go 1.22 every turn of a `for` loop has its
own variable.

## A component that holds objects

Take the objects as a function, as `ui.VGroup` does, and return the box that holds them, so the program can add more
later with `Add()`:

```go
// Section is a titled part of a form. content makes its objects; section.Add(...) puts more in later.
func Section(title string, content func()) *ui.Object {
	var body *ui.Object
	ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 8, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: strings.ToUpper(title), FontSize: 11, Bold: true, TextColor: "#5A6270"})
		body = ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 10, Align: "left top"}, content)
	})
	return body
}

sound := Section("Sound", func() {
	ui.CheckBox(ui.CheckBoxProps{LabelText: "Play a sound"})
	ui.Slider(ui.SliderProps{Value: 7, MinValue: 0, MaxValue: 10})
})
sound.Add(func() { ui.Toggle(ui.ToggleProps{Value: true}) })
```

## Sharing it: a package

A component used by more than one file goes into a package of its own: a folder of your module.

```
myapp/
  go.mod          module myapp
  main.go         import "myapp/widgets"  ...  widgets.NewCard("Sales", "Updated today.", nil)
  widgets/
    card.go       package widgets  (imports ui "github.com/bug7a/ozdenui-go")
```

- The names other packages use start with a capital letter: `NewCard`, `Card.Box`, `SetTitle`.
- For many programs, give the package a module of its own (`go mod init github.com/you/widgets`), and
  `go get github.com/you/widgets` in each program.
- A `go doc` comment over each exported name tells its user what it takes, as the functions of the kit do.

## What belongs to Go, what to the page

**"A click, then a change" belongs to Go.** Your function runs, and the page shows the change when the function ends.

**"Every frame, every mouse move" belongs to the page.** A motion made of many steps sent from Go would stutter.
Tell the page how to move, then give the end:

```go
ball.Call("setMotion", "left 0.3s, opacity 0.2s") // these values now change smoothly
ball.Set("left", 300)                             // the page moves it there in 0.3 s
```

A `Canvas` does so by itself: the page draws the free lines and moves the dragged shapes, Go hears only their ends.

### A callback that must answer at once

Some callbacks of the components must return an answer at once (a function that draws a row of a list, a check that
stops a click). The answer of Go comes a little later, so such a callback can not be a Go function: format the values
in Go before they go to the page instead, and use the params the components give (chapter 10 lists them).

## Next

[Chapter 10: Going further](10-going-further.md).
