/* Bismillah */

// A chat with an AI that is only a sample: for every message it writes a small plan, word by word, as an AI that
// streams its answer does (in a goroutine, so the window keeps working). Stop ends the answer, New chat starts again.
// Your AI goes into writeAnswer(): every piece it gives goes to the message with chat.Call("appendText", key, piece).
package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	lightGround = "#FFFFFF" // the ground of the ChatView in each theme
	darkGround  = "#1C1C1B"
	thinkTime   = 600 * time.Millisecond // before the first words: the dots
	wordTime    = 40 * time.Millisecond  // between two words
)

const answerText = "Here is a small plan for **{topic}**:\n\n" +
	"1. Write down what you want to do\n2. Split it into small steps\n3. Do one step a day\n\n" +
	"| Day | Step |\n|:---|:---|\n| Monday | The list |\n| Tuesday | The first step |\n\n" +
	"```python\nsteps = [\"the list\", \"the first step\"]\nprint(len(steps), \"steps\")\n```\n\n" +
	"Ask again if something is not clear."

var firstMessages = []ui.P{{"role": "system", "text": "A sample: the answers are written by this program, not by an AI"}}

// One word with the spaces before it, as a stream gives the pieces.
var piecePattern = regexp.MustCompile(`\s*\S+|\s+$`)

var (
	answerNumber = 0 // WHY: Stop and New chat start a new number, and the pieces of an older answer are dropped
	chat         *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "A chat", Width: 720, Height: 760})
}

func start(page *ui.Object) {
	ground := lightGround
	if ui.IsDark() {
		ground = darkGround
	}
	page.Set("color", ground)
	ui.VGroup(ui.Props{Gap: 0}, func() {
		ui.HGroup(ui.Props{Width: "100%", Height: "auto", Align: "left center", Gap: 12, Padding: []int{16, 10}}, func() {
			ui.Label(ui.Props{Text: "A chat", FontSize: 17, Bold: true, Fill: true})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "New chat", OnClick: newChat})
		})
		chat = ui.ChatView(ui.ChatViewProps{Width: "100%", Fill: true, Messages: firstMessages, EmptyText: "Ask something",
			OnSend: answer, OnStop: stopAnswer})
		chat.Focus()
	})
}

func answer(view *ui.Object, text string) {
	answerNumber++
	key := fmt.Sprintf("answer%d", answerNumber)
	chat.Call("addMessage", ui.P{"key": key, "role": "assistant"})
	chat.Set("busy", true)
	go writeAnswer(answerNumber, key, text)
}

// In a goroutine: the window keeps working while the answer comes.
func writeAnswer(number int, key, topic string) {
	time.Sleep(thinkTime)
	for _, piece := range piecePattern.FindAllString(strings.ReplaceAll(answerText, "{topic}", topic), -1) {
		ui.Do(func() {
			if number == answerNumber {
				chat.Call("appendText", key, piece)
			}
		})
		time.Sleep(wordTime)
	}
	ui.Do(func() {
		if number == answerNumber { // WHY: after Stop the page set busy to 0 already
			chat.Set("busy", false)
		}
	})
}

func stopAnswer() {
	answerNumber++
}

func newChat() {
	answerNumber++
	chat.Set("busy", false)
	chat.Set("messages", firstMessages)
}
