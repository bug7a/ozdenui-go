/* Bismillah */

/*

app.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The window, the tasks of the program and the commands that go to the page.

Goroutines:
- The main goroutine runs the window (the system needs its first thread).
- The app goroutine runs the code of the program: start(page), every event handler, every timer and Do() task, one
  at a time. The code of the program never needs a lock.
- The server (server.go) takes the messages of the page on goroutines of its own: they only put a task in the queue.

The commands of a task are sent to the page together, when the task ends (flush): one message for a whole start(),
not one for each object.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Version is the version of the module.
const Version = "0.2.1"

const (
	readyTimeout    = 15 * time.Second
	darkPageColor   = "#1C1C1B"
	earlyPageColor  = "ozdenui-page-color:" // window.name of a page loaded again for another theme (bridge.js)
	reloadDelayMs   = 50                    // the page of SetTheme() loads itself again this much after the script that asks it
	minRepeat       = time.Millisecond      // WHY: Every(0) would spin a whole CPU core
	accentStyleName = "minimalAccent"
)

// Window is how the window opens: ozdenui.Run(start, ozdenui.Window{Title: "Counter", Width: 400, Height: 300})
type Window struct {
	Title  string // "OzdenUI" when empty
	Width  int    // 800 when 0
	Height int    // 600 when 0
	// Theme: "light" (default), "dark" or "auto" (what the system asks). In a dark window every component takes its
	// dark look unless it is given another one, and the plain objects and the page are dark too.
	Theme string
	// Style: "" (the usual looks) or "minimal": small and made of lines, for the components that have it.
	Style string
	// Accent: a color of the minimal look ("#2A78D6"): an on toggle, a checked box, the selected tab and the main
	// button (SmartButton with Tone: "accent") take it. AccentDark: its color in a dark window ("": made from Accent).
	Accent     string
	AccentDark string
	// Language: the language of the texts of the window ("en" when empty; "tr", "de"...). The page writes capital
	// letters, breaks words and sorts by it: a table title in Turkish is "SİPARİŞLER" only with "tr".
	Language string
	// Debug: the developer tools of the web view (right click > Inspect).
	Debug bool
}

type application struct {
	window      Window
	start       func(page *Object)
	server      *server
	control     windowControl
	wantedTheme string
	theme       string // "light" or "dark"

	tasks    taskQueue
	ready    chan struct{} // closed when the page said it is ready for the first time
	readyMu  sync.Mutex
	isReady  bool
	closed   chan struct{} // closed when the window is gone
	drawn    chan struct{} // closed when the first drawing of start() is on the page
	isClosed bool

	appGoroutine uint64
	commands     []map[string]any
	containers   []string // the open boxes of the task that runs: the last one takes the new objects
	objects      map[string]*Object
	leaving      map[string]*Object // objects the page removed, until the events they sent before are handled
	callbacks    map[string]*callback
	failedIDs    map[string]bool // the objects the page could not make (an answer that belongs to one is nil)
	answers      map[string]*answer
	timers       map[*timer]bool
	timersMu     sync.Mutex

	nextObject        int
	nextCallback      int
	prefersDark       bool
	pageWidth         float64
	pageHeight        float64
	isLoadingForTheme bool
	hasStarted        bool
	page              *Object
	basicVersion      string
	warned            map[string]bool
}

type callback struct {
	handler any
	ownerID string
}

// An answer the code waits for (a dialog): from the page, or nil when it can not come. Only the first one counts.
type answer struct {
	value   any
	isGiven bool
}

func (a *answer) give(value any) {
	if a.isGiven {
		return
	}
	a.value = value
	a.isGiven = true
}

var current *application

func init() {
	// WHY: the window of macOS works only on the first thread of the program; Run() is called from main().
	runtime.LockOSThread()
}

// Run opens the window, calls start(page) when its page is ready and returns when the window is closed.
// Call it once, at the end of main().
func Run(start func(page *Object), window ...Window) {
	options := Window{}
	if len(window) > 0 {
		options = window[0]
	}
	if options.Title == "" {
		options.Title = "OzdenUI"
	}
	if options.Width == 0 {
		options.Width = 800
	}
	if options.Height == 0 {
		options.Height = 600
	}
	if options.Theme == "" {
		options.Theme = "light"
	}
	if theme := os.Getenv("OZDENUI_THEME"); theme != "" {
		options.Theme = theme // a picture of the dark window without a change of the program (OZDENUI_LOOK)
	}
	checkWindow(options)
	if current != nil {
		panic("ozdenui: Run() is called once in a program")
	}
	app := &application{
		window:      options,
		start:       start,
		wantedTheme: options.Theme,
		theme:       "light",
		ready:       make(chan struct{}),
		closed:      make(chan struct{}),
		drawn:       make(chan struct{}),
		objects:     map[string]*Object{},
		leaving:     map[string]*Object{},
		callbacks:   map[string]*callback{},
		failedIDs:   map[string]bool{},
		answers:     map[string]*answer{},
		timers:      map[*timer]bool{},
		warned:      map[string]bool{},
	}
	if options.Theme == "dark" {
		app.theme = "dark"
	}
	app.tasks.init()
	server, err := startServer(app)
	if err != nil {
		panic("ozdenui: the page could not be served: " + err.Error())
	}
	app.server = server
	current = app
	go app.runAppGoroutine()
	openWindow(app, server.address()) // returns when the window is closed
	app.onWindowClosed()
	server.stop()
}

func checkWindow(options Window) {
	switch options.Theme {
	case "light", "dark", "auto":
	default:
		panic(fmt.Sprintf(`ozdenui: Window{Theme: %q}: it is "light", "dark" or "auto"`, options.Theme))
	}
	if options.Style != "" && options.Style != "minimal" {
		panic(fmt.Sprintf(`ozdenui: Window{Style: %q}: it is "" or "minimal"`, options.Style))
	}
	if options.AccentDark != "" && options.Accent == "" {
		panic("ozdenui: Window{AccentDark: ...}: it goes with Accent")
	}
	if options.Accent != "" && options.Style != "minimal" {
		panic(`ozdenui: Window{Accent: ...}: the color of the minimal look; give Style: "minimal" too`)
	}
}

func currentApp() *application {
	if current == nil {
		panic("ozdenui: create the objects in start(page): ozdenui.Run(start)")
	}
	return current
}

// *** THE WINDOW, THE TIME AND THE TASKS, from any function of the program.

// Close closes the window (from any goroutine); Run() returns.
func Close() {
	app := currentApp()
	if app.control != nil {
		app.control.close()
	}
}

// SetTitle changes the title of the window.
func SetTitle(title string) {
	app := currentApp()
	app.window.Title = title
	if app.control != nil {
		app.control.setTitle(title)
	}
	app.setPageTitle(title)
}

// Theme is the theme of the window: "light" or "dark".
func Theme() string { return currentApp().theme }

// IsDark is true in a dark window: color := map[bool]string{false: "#5A6270", true: "#9AA3B5"}[ozdenui.IsDark()]
func IsDark() bool { return currentApp().theme == "dark" }

// SetTheme draws the window again in another theme: "light", "dark" or "auto" (what the system asks). The page is
// made again and start() runs again, as when the window opened: its objects start again (a text typed in a field is
// gone), what the program keeps in its own variables stays. Any goroutine may call it; the change comes after the task
// that runs now.
func SetTheme(theme string) {
	switch theme {
	case "light", "dark", "auto":
	default:
		panic(fmt.Sprintf(`ozdenui: SetTheme(%q): it is "light", "dark" or "auto"`, theme))
	}
	app := currentApp()
	app.tasks.put(func() { app.changeTheme(theme) })
}

// Do runs task() on the app goroutine, after the task that runs now. Any goroutine may call it: it is the way a
// goroutine of the program (a download, a long job) changes the objects.
//
//	go func() {
//	    text := download()                       // the window does not wait for it
//	    ozdenui.Do(func() { label.SetText(text) })
//	}()
func Do(task func()) {
	currentApp().tasks.put(task)
}

// After runs task() once, after a time (milliseconds). It returns a function that cancels it.
func After(milliseconds int, task func()) (cancel func()) {
	return currentApp().newTimer(milliseconds, task, false)
}

// Every runs task() again and again (the loop() of basic.js). It returns a function that stops it.
func Every(milliseconds int, task func()) (stop func()) {
	return currentApp().newTimer(milliseconds, task, true)
}

// Flush sends the waiting commands to the page now: a long task can show its progress.
func Flush() {
	app := currentApp()
	app.check("Flush()")
	app.flush()
}

// Evaluate runs JavaScript on the page and returns its value (after the waiting commands).
func Evaluate(script string) any {
	app := currentApp()
	app.check("Evaluate()")
	app.flush()
	value, err := app.server.request(map[string]any{"kind": "eval", "script": script})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ozdenui: %s: Evaluate(): %v\n", sourceLine(), err)
		return nil
	}
	return app.goValue(value)
}

// *** TIMERS:

type timer struct {
	app      *application
	stopped  chan struct{}
	once     sync.Once
	isQueued bool // WHY: a slow task must not pile up ticks in the queue
	mu       sync.Mutex
}

func (app *application) newTimer(milliseconds int, task func(), repeat bool) func() {
	wait := time.Duration(milliseconds) * time.Millisecond
	if wait < 0 {
		wait = 0
	}
	if repeat && wait < minRepeat {
		wait = minRepeat
	}
	t := &timer{app: app, stopped: make(chan struct{})}
	app.timersMu.Lock()
	app.timers[t] = true
	app.timersMu.Unlock()
	run := func() {
		t.mu.Lock()
		t.isQueued = false
		t.mu.Unlock()
		select {
		case <-t.stopped:
			return
		default:
		}
		task()
	}
	go func() {
		for {
			select {
			case <-t.stopped:
				return
			case <-app.closed:
				return
			case <-time.After(wait):
			}
			t.mu.Lock()
			isQueued := t.isQueued
			t.isQueued = true
			t.mu.Unlock()
			if !isQueued {
				app.tasks.put(run)
			}
			if !repeat {
				app.forgetTimer(t)
				return
			}
		}
	}()
	return t.stop
}

func (t *timer) stop() {
	t.once.Do(func() { close(t.stopped) })
	t.app.forgetTimer(t)
}

func (app *application) forgetTimer(t *timer) {
	app.timersMu.Lock()
	delete(app.timers, t)
	app.timersMu.Unlock()
}

func (app *application) stopTimers() {
	app.timersMu.Lock()
	timers := make([]*timer, 0, len(app.timers))
	for t := range app.timers {
		timers = append(timers, t)
	}
	app.timersMu.Unlock()
	for _, t := range timers {
		t.stop()
	}
}

// *** THE TASK QUEUE:

type taskQueue struct {
	mu       sync.Mutex
	cond     *sync.Cond
	tasks    []func()
	isClosed bool
}

func (q *taskQueue) init() { q.cond = sync.NewCond(&q.mu) }

func (q *taskQueue) put(task func()) {
	q.mu.Lock()
	q.tasks = append(q.tasks, task)
	q.mu.Unlock()
	q.cond.Signal()
}

func (q *taskQueue) close() {
	q.mu.Lock()
	q.isClosed = true
	q.mu.Unlock()
	q.cond.Broadcast()
}

// The next task; ok is false when the window is closed.
func (q *taskQueue) get() (task func(), ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.tasks) == 0 && !q.isClosed {
		q.cond.Wait()
	}
	if q.isClosed {
		return nil, false
	}
	task = q.tasks[0]
	q.tasks = q.tasks[1:]
	return task, true
}

// *** THE APP GOROUTINE:

func (app *application) runAppGoroutine() {
	app.appGoroutine = goroutineID()
	select {
	case <-app.ready:
	case <-app.closed:
		return
	case <-time.After(readyTimeout):
		fmt.Fprintln(os.Stderr, "ozdenui: the page is not ready yet (basic.js or the bridge did not load)")
		// WHY it goes on waiting: a page that comes late on a slow computer must still start the program
		select {
		case <-app.ready:
		case <-app.closed:
			return
		}
	}
	app.runTask(app.startPage)
	close(app.drawn)
	app.runTasksUntil(func() bool { return false })
}

func (app *application) startPage() {
	app.hasStarted = true
	app.makePage()
	app.addAccentPackages()
	app.start(app.page)
}

func (app *application) makePage() {
	app.page = newPage(app)
	if app.theme == "dark" {
		app.page.Set("color", darkPageColor)
	}
}

func (app *application) addAccentPackages() {
	if app.window.Accent == "" {
		return
	}
	app.flush()
	_, err := app.server.request(map[string]any{
		"kind": "accent", "name": accentStyleName, "color": app.window.Accent, "colorDark": app.window.AccentDark,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ozdenui: Window{Accent}: %v\n", err)
	}
}

func (app *application) runTasksUntil(isDone func() bool) {
	for !isDone() {
		task, ok := app.tasks.get()
		if !ok {
			return
		}
		app.runTask(task)
	}
}

func (app *application) runTask(task func()) {
	// WHY its own open boxes: a task that runs while another one waits for a dialog (a click, a timer) would make its
	//      objects in the open box of that one; a task that failed in a box leaves nothing open either.
	outer := app.containers
	app.containers = nil
	func() {
		defer func() {
			if problem := recover(); problem != nil {
				// WHY not the end of the program: one handler that failed must not close the window with the work in it
				fmt.Fprintf(os.Stderr, "ozdenui: a task failed: %v\n%s", problem, userStack())
			}
		}()
		task()
	}()
	app.containers = outer
	app.flush()
}

// Runs the other tasks until isDone(): the code that waits (a dialog) keeps the window working.
func (app *application) waitFor(isDone func() bool) {
	app.flush()
	app.runTasksUntil(func() bool { return isDone() || app.isClosed })
}

func (app *application) check(what string) {
	if goroutineID() != app.appGoroutine {
		panic("ozdenui: " + what + " works only in start(), an event or a timer. From another goroutine: ozdenui.Do(func() { ... })")
	}
}

// *** THE COMMANDS:

func (app *application) send(command map[string]any) {
	app.check(fmt.Sprint(command["op"]))
	app.commands = append(app.commands, command)
}

func (app *application) flush() {
	if len(app.commands) == 0 || app.isClosed || app.isLoadingForTheme {
		return
	}
	commands := app.commands
	app.commands = nil
	value, err := app.server.request(map[string]any{"kind": "apply", "commands": commands})
	if err != nil {
		return // the page is gone (closed, loaded again): its commands go with it
	}
	app.afterApply(value)
}

type applyAnswer struct {
	Errors []struct {
		Op      string `json:"op"`
		ID      any    `json:"id"`
		Src     string `json:"src"`
		Message string `json:"message"`
	} `json:"errors"`
	Removed []string `json:"removed"`
}

func (app *application) afterApply(value json.RawMessage) {
	var result applyAnswer
	if err := json.Unmarshal(value, &result); err != nil {
		return
	}
	for _, problem := range result.Errors {
		id := ""
		if problem.ID != nil {
			id = fmt.Sprint(problem.ID)
		}
		where := strings.TrimSpace(problem.Op + " " + id)
		if problem.Src != "" {
			where = problem.Src + ": " + where
		}
		fmt.Fprintf(os.Stderr, "ozdenui: %s: %s\n", where, problem.Message)
		if problem.Op == "create" && id != "" {
			app.giveUpAnswerOf(id)
		}
	}
	if len(result.Removed) == 0 {
		return
	}
	for _, id := range result.Removed {
		if obj, ok := app.objects[id]; ok {
			delete(app.objects, id)
			app.leaving[id] = obj
		}
	}
	// WHY later: a component sends an event and then leaves the page (a Dialog: onClose, then remove()). The events the
	//     page sent before this answer are already in the queue: the task that forgets comes after them.
	removed := result.Removed
	app.tasks.put(func() { app.forgetRemoved(removed) })
}

// WHY: the functions given to a removed object would be kept for ever (a list made again every second grows).
func (app *application) forgetRemoved(removed []string) {
	gone := map[string]bool{}
	for _, id := range removed {
		gone[id] = true
		delete(app.leaving, id)
	}
	for id, entry := range app.callbacks {
		if gone[entry.ownerID] {
			delete(app.callbacks, id)
		}
	}
}

func (app *application) giveUpAnswerOf(id string) {
	app.failedIDs[id] = true
	if pending, ok := app.answers[id]; ok {
		pending.give(nil)
	}
}

func (app *application) newID() string {
	app.nextObject++
	return fmt.Sprintf("o%d", app.nextObject)
}

func (app *application) registerCallback(handler any, ownerID string) string {
	app.nextCallback++
	id := fmt.Sprintf("c%d", app.nextCallback)
	app.callbacks[id] = &callback{handler: handler, ownerID: ownerID}
	return id
}

func (app *application) containerID() string {
	if len(app.containers) == 0 {
		return "page"
	}
	return app.containers[len(app.containers)-1]
}

// Calls a function of a component (Toast.show, Dialog.confirm). wait: the function gives a Promise; its value is
// returned, and the window keeps working while the code waits for it. keep: the id of the object that stands for what
// the function returns (a bubble).
func (app *application) callStatic(target, method string, args []any, wait bool, keep string) any {
	command := map[string]any{"op": "static", "target": target, "method": method, "args": args, "src": sourceLine()}
	if keep != "" {
		command["keep"] = keep
	}
	if !wait {
		app.send(command)
		return nil
	}
	pending := &answer{}
	id := app.registerCallback(func(value any) { pending.give(value) }, "")
	command["answer"] = id
	app.send(command)
	value := app.waitForAnswer(id, pending)
	delete(app.callbacks, id)
	return value
}

// Waits for an answer (the window keeps working). key: the callback id of the answer, or the id of the object that
// gives it (a Dialog). It is nil when it can not come: the object could not be made, the page was loaded again or
// the window was closed.
func (app *application) waitForAnswer(key string, pending *answer) any {
	app.answers[key] = pending
	defer delete(app.answers, key)
	app.waitFor(func() bool {
		if app.failedIDs[key] {
			pending.give(nil)
		}
		return pending.isGiven
	})
	return pending.value
}

// *** FROM THE PAGE (the goroutines of the server; they only queue tasks):

type readyMessage struct {
	Page         string  `json:"page"`
	BasicVersion string  `json:"basicVersion"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	PrefersDark  bool    `json:"prefersDark"`
}

func (app *application) onPageReady(message readyMessage) {
	app.readyMu.Lock()
	isFirst := !app.isReady
	app.isReady = true
	app.readyMu.Unlock()
	if isFirst {
		app.takeReady(message)
		if fitter, ok := app.control.(interface {
			fitPage(int, int, float64, float64)
		}); ok {
			fitter.fitPage(app.window.Width, app.window.Height, message.Width, message.Height)
		}
		close(app.ready)
		return
	}
	// The page was loaded again (SetTheme, the reload of the developer tools): its objects are gone, start() makes
	// them again.
	app.tasks.put(func() {
		app.takeReady(message)
		app.startAgain()
	})
}

func (app *application) takeReady(message readyMessage) {
	app.basicVersion = message.BasicVersion
	app.pageWidth, app.pageHeight = message.Width, message.Height
	app.prefersDark = message.PrefersDark
	if app.wantedTheme == "auto" {
		app.theme = "light"
		if message.PrefersDark {
			app.theme = "dark"
		}
	}
}

func (app *application) onPageEvent(target, event string, data json.RawMessage) {
	app.tasks.put(func() { app.dispatch(target, event, data) })
}

func (app *application) dispatch(target, event string, raw json.RawMessage) {
	data := map[string]any{}
	_ = json.Unmarshal(raw, &data)
	switch event {
	case "callback":
		app.runCallback(target, data)
		return
	case "__message":
		fmt.Fprintf(os.Stderr, "ozdenui (page %v): %v\n", data["kind"], data["text"])
		return
	}
	obj := app.objects[target]
	if obj == nil {
		return // an event of an object removed in the meantime
	}
	obj.receive(event, data)
}

func (app *application) runCallback(id string, data map[string]any) {
	entry := app.callbacks[id]
	if entry == nil {
		return
	}
	owner := app.objects[entry.ownerID]
	if owner == nil {
		owner = app.leaving[entry.ownerID]
	}
	if state, ok := data["state"].(map[string]any); ok && owner != nil {
		owner.takeState(state)
	}
	args, _ := data["args"].([]any)
	for i, value := range args {
		args[i] = app.goValueOf(value)
	}
	callHandler(entry.handler, args)
}

// A value from the page: a basic.js object becomes its *Object; JSON data that came whole (the wholeArgs of a
// callback) becomes its value.
func (app *application) goValue(raw json.RawMessage) any {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}
	return app.goValueOf(value)
}

func (app *application) goValueOf(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		if id, has := typed["__bpObject"]; has {
			return app.objectValue(id, typed)
		}
		if text, has := typed["__bpJson"].(string); has {
			var whole any
			_ = json.Unmarshal([]byte(text), &whole)
			return whole
		}
		for key, item := range typed {
			typed[key] = app.goValueOf(item)
		}
		return typed
	case []any:
		for i, item := range typed {
			typed[i] = app.goValueOf(item)
		}
		return typed
	}
	return value
}

// A basic.js object from the page: its *Object, or {"id", "data"} for a part of a component (a row of a list), or nil.
func (app *application) objectValue(id any, value map[string]any) any {
	if text, ok := id.(string); ok {
		if obj := app.objects[text]; obj != nil {
			return obj
		}
		if obj := app.leaving[text]; obj != nil {
			return obj
		}
	}
	rest := map[string]any{}
	for key, item := range value {
		if key != "__bpObject" {
			rest[key] = app.goValueOf(item)
		}
	}
	if len(rest) == 0 {
		return nil
	}
	return rest
}

// *** ANOTHER THEME, A PAGE LOADED AGAIN:

// WHY the page is loaded again: a component takes its look when it is made, and a new page keeps nothing of the old
// one (an open menu or bubble, a listener of a component).
// WHY window.name: it stays over the load, so the new page has the ground of its theme from its first drawing.
func (app *application) changeTheme(theme string) {
	if app.isClosed || theme == app.wantedTheme {
		return
	}
	app.wantedTheme = theme
	if theme == "auto" {
		theme = "light"
		if app.prefersDark {
			theme = "dark" // the page tells it again when it is ready
		}
	}
	app.theme = theme
	app.stopTimers() // WHY: a tick would send its commands to a page that is being loaded
	color := ""
	if theme == "dark" {
		color = darkPageColor
	}
	app.flush()
	app.isLoadingForTheme = true
	name, _ := json.Marshal(earlyPageColor + color)
	script := fmt.Sprintf("window.name = %s; setTimeout(function () { location.reload(); }, %d); 0", name, reloadDelayMs)
	_, _ = app.server.request(map[string]any{"kind": "eval", "script": script})
}

func (app *application) startAgain() {
	if app.isLoadingForTheme {
		app.isLoadingForTheme = false // WHY no word: the program asked for the new page (SetTheme)
	} else if app.hasStarted {
		fmt.Fprintln(os.Stderr, "ozdenui: the page was loaded again: start() runs again")
	}
	// WHY: the dialogs of the old page are gone; the code that waits for one of them gets nil
	for _, pending := range app.answers {
		pending.give(nil)
	}
	app.failedIDs = map[string]bool{}
	app.stopTimers()
	app.commands = nil
	app.objects = map[string]*Object{}
	app.leaving = map[string]*Object{}
	app.callbacks = map[string]*callback{}
	app.startPage()
}

func (app *application) onWindowClosed() {
	app.isClosed = true
	close(app.closed)
	app.tasks.close()
}

// *** SMALL HELPERS:

// The number of the goroutine that runs. WHY: Go gives no other way to know that an object is changed from a
// goroutine of the program, which would break the order of the commands without a word.
func goroutineID() uint64 {
	var buffer [64]byte
	text := string(buffer[:runtime.Stack(buffer[:], false)]) // "goroutine 18 [running]:"
	text = strings.TrimPrefix(text, "goroutine ")
	var id uint64
	for _, digit := range text {
		if digit < '0' || digit > '9' {
			break
		}
		id = id*10 + uint64(digit-'0')
	}
	return id
}

// The lines of the program in the stack of a failed task (not the ones of this module and of Go).
func userStack() string {
	lines := strings.Split(string(debug.Stack()), "\n")
	kept := []string{}
	for i := 0; i+1 < len(lines); i++ {
		isFile := strings.HasPrefix(lines[i+1], "\t")
		if !isFile || strings.HasPrefix(lines[i], "\t") {
			continue
		}
		isOurs := strings.HasPrefix(lines[i], modulePath+".") || strings.HasPrefix(lines[i], modulePath+"/internal")
		isGo := strings.HasPrefix(lines[i], "runtime") || strings.HasPrefix(lines[i], "reflect.") || strings.HasPrefix(lines[i], "panic(")
		if isOurs || isGo {
			continue
		}
		kept = append(kept, lines[i], lines[i+1])
	}
	if len(kept) == 0 {
		return ""
	}
	return strings.Join(kept, "\n") + "\n"
}
