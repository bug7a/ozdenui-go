/* Bismillah */

// Catch the dots: move the square with the arrow keys and touch the dot before the time is over.
// Objects at a place of their own (Left, Top), the keys of the page and a timer.
package main

import (
	"fmt"
	"math/rand"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	pageColor   = "#14213D"
	fieldColor  = "#1F2E52"
	playerColor = "#FCA311"
	dotColor    = "#2EC4B6"
	textLight   = "#E5E5E5"
	fieldWidth  = 560
	fieldHeight = 360
	playerSize  = 34
	dotSize     = 18
	step        = 22
	gameSeconds = 30
)

var moves = map[string][2]int{
	"ArrowLeft":  {-step, 0},
	"ArrowRight": {step, 0},
	"ArrowUp":    {0, -step},
	"ArrowDown":  {0, step},
}

var (
	score     = 0
	secsLeft  = gameSeconds
	running   = false
	stopTimer func()

	// The places are kept in the program: no question to the page at every key.
	playerLeft, playerTop int
	dotLeft, dotTop       int

	scoreLabel *ui.Object
	timeLabel  *ui.Object
	player     *ui.Object
	dot        *ui.Object
	playButton *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Catch the dots", Width: 620, Height: 520})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 14, Padding: 24}, func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 40}, func() {
			scoreLabel = ui.Label(ui.Props{Text: "Score: 0", FontSize: 20, Bold: true, TextColor: textLight})
			timeLabel = ui.Label(ui.Props{Text: fmt.Sprintf("Time: %d", gameSeconds), FontSize: 20, Bold: true, TextColor: textLight})
		})
		ui.Box(ui.Props{Width: fieldWidth, Height: fieldHeight, Color: fieldColor, Round: 10}, func() {
			dot = ui.Box(ui.Props{Left: 0, Top: 0, Width: dotSize, Height: dotSize, Color: dotColor, Round: dotSize / 2})
			player = ui.Box(ui.Props{Left: 0, Top: 0, Width: playerSize, Height: playerSize, Color: playerColor, Round: 6})
		})
		playButton = ui.SmartButton(ui.SmartButtonProps{LabelText: "Play", StyleName: "ozdenDark", OnClick: play})
	})
	player.Call("setMotion", "left 0.08s, top 0.08s")
	page.On("keydown", func(page *ui.Object, event ui.Event) { move(event.Key) })
	placePlayerInMiddle()
	placeDot()
}

func placePlayer(left, top int) {
	playerLeft, playerTop = left, top
	player.Set("left", left, "top", top)
}

func placePlayerInMiddle() {
	placePlayer((fieldWidth-playerSize)/2, (fieldHeight-playerSize)/2)
}

func placeDot() {
	dotLeft = rand.Intn(fieldWidth - dotSize + 1)
	dotTop = rand.Intn(fieldHeight - dotSize + 1)
	dot.Set("left", dotLeft, "top", dotTop)
}

func play() {
	if running {
		return
	}
	score, secsLeft, running = 0, gameSeconds, true
	scoreLabel.SetText("Score: 0")
	timeLabel.SetText(fmt.Sprintf("Time: %d", gameSeconds))
	playButton.Set("enabled", false)
	placePlayerInMiddle()
	placeDot()
	stopTimer = ui.Every(1000, tick)
}

func tick() {
	secsLeft--
	timeLabel.SetText(fmt.Sprintf("Time: %d", secsLeft))
	if secsLeft > 0 {
		return
	}
	running = false
	stopTimer()
	timeLabel.SetText("Time is over")
	playButton.Set("enabled", true)
}

func move(key string) {
	change, isArrow := moves[key]
	if !running || !isArrow {
		return
	}
	left := min(max(playerLeft+change[0], 0), fieldWidth-playerSize)
	top := min(max(playerTop+change[1], 0), fieldHeight-playerSize)
	placePlayer(left, top)
	if touches() {
		score++
		scoreLabel.SetText(fmt.Sprintf("Score: %d", score))
		placeDot()
	}
}

// touches is true when the player and the dot overlap.
func touches() bool {
	return playerLeft < dotLeft+dotSize && dotLeft < playerLeft+playerSize &&
		playerTop < dotTop+dotSize && dotTop < playerTop+playerSize
}
