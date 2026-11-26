/* Bismillah */

/*

window.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
How the window opens: the web view of the system, or the window of a browser of this computer.

- macOS: the web view of the system (WKWebView), in the program itself (cgo: the command line tools of Xcode).
- Windows: the web view of the system (WebView2, window_windows.go), without cgo; without WebView2 the app window below.
- Linux: the web view of the system (WebKitGTK, window_linux.go), loaded when the program runs, without cgo; without it
  the app window below.
- macOS without cgo, and the computers above without their web view: an app window of Chrome, Chromium or Edge (no
  address bar, no tabs). With no such browser: a tab of the default browser.
- OZDENUI_LOOK=picture.png: no window; the page is drawn by a hidden Chrome and saved as a picture (look.go).

Environment: OZDENUI_WINDOW=browser (the browser window also where the web view works), =tab (the default browser),
=headless (a hidden Chrome: a check without a screen; the program closes it with Close()).

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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type windowControl interface {
	setTitle(title string)
	close()
}

const streamGoneAfter = 1500 * time.Millisecond // a page that does not listen this long is a closed window

// WHY: Edge signed the profile of the program in with the Microsoft account of Windows on its own, and synced the
// passwords, the history and the tabs of the user into it, with a dialog over the window. (One flag of
// --disable-features: a second one would replace it.)
var profileFlags = []string{"--disable-sync", "--disable-features=msImplicitSignin"}

// Opens the window and returns when it is closed.
func openWindow(app *application, address string) {
	mode := os.Getenv("OZDENUI_WINDOW")
	if os.Getenv("OZDENUI_LOOK") != "" || mode == "headless" {
		if err := openHeadless(app, address); err != nil {
			fmt.Fprintf(os.Stderr, "ozdenui: %v\n", err)
		}
		return
	}
	if mode == "" && openNativeWindow(app, address) {
		return
	}
	openBrowserWindow(app, address, mode == "tab")
}

// *** A WINDOW OF A BROWSER:

type browserWindow struct {
	app     *application
	process *exec.Cmd
	done    chan struct{}
	once    sync.Once
}

func (w *browserWindow) setTitle(title string) {} // the page changes its own title (SetTitle)

func (w *browserWindow) close() {
	w.once.Do(func() {
		// WHY the page closes itself: the process of the browser may be one that was already open
		go func() { _, _ = w.app.server.request(map[string]any{"kind": "eval", "script": "window.close(); 0"}) }()
		time.AfterFunc(300*time.Millisecond, func() { close(w.done) })
	})
}

func openBrowserWindow(app *application, address string, isTab bool) {
	window := &browserWindow{app: app, done: make(chan struct{})}
	app.control = window
	watchStream(app, func() { window.once.Do(func() { close(window.done) }) })
	browser := ""
	if !isTab {
		browser = findBrowser()
	}
	if browser != "" {
		// WHY a profile of the program: the browser is a process of its own then, which ends with the window (a browser
		//      that was already open would take the window into its process, with its tabs and its size)
		profile := browserProfile(browser)
		_ = os.MkdirAll(profile, 0o755)
		arguments := append([]string{
			"--app=" + address, fmt.Sprintf("--window-size=%d,%d", app.window.Width, app.window.Height),
			"--user-data-dir=" + profile, "--no-first-run", "--no-default-browser-check"}, profileFlags...)
		window.process = exec.Command(browser, arguments...)
		if err := window.process.Start(); err != nil {
			browser = ""
		} else {
			go func() { _ = window.process.Wait() }()
		}
	}
	if browser == "" {
		if err := openInDefaultBrowser(address); err != nil {
			fmt.Fprintf(os.Stderr, "ozdenui: no browser could be opened (%v). The page of the program: %s\n", err, address)
		}
	}
	<-window.done
	if window.process != nil && window.process.Process != nil {
		_ = window.process.Process.Kill() // WHY: on macOS a browser without a window keeps running
	}
}

func programName() string {
	name := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
	safe := strings.Map(func(letter rune) rune {
		if letter >= 'a' && letter <= 'z' || letter >= 'A' && letter <= 'Z' || letter >= '0' && letter <= '9' || letter == '-' {
			return letter
		}
		return '_'
	}, name)
	return safe
}

// Calls onGone when the page stopped listening and did not come back (a closed tab or window; a page that is loaded
// again comes back at once).
func watchStream(app *application, onGone func()) {
	var mu sync.Mutex
	var timer *time.Timer
	open := 0
	app.server.mu.Lock()
	app.server.onStream = func(isOpen bool) {
		mu.Lock()
		defer mu.Unlock()
		if isOpen {
			open++
			if timer != nil {
				timer.Stop()
				timer = nil
			}
			return
		}
		open--
		if open == 0 {
			timer = time.AfterFunc(streamGoneAfter, onGone)
		}
	}
	app.server.mu.Unlock()
}

func cacheFolder() string {
	folder, err := os.UserCacheDir()
	if err != nil {
		folder = os.TempDir()
	}
	return filepath.Join(folder, "ozdenui")
}

// The profile of the program in a browser.
func browserProfile(browser string) string {
	// WHY the folder of a snap: Ubuntu's Chromium is a snap, which can not write a hidden folder of the home (~/.cache):
	//      it showed "Failed To Create Data Directory" and the page never came. Its own ~/snap/<name>/common it can.
	if name := snapName(browser); name != "" {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "snap", name, "common", "ozdenui", "browser-profile-"+programName())
		}
	}
	return filepath.Join(cacheFolder(), "browser-profile-"+programName())
}

// The snap a browser command starts: /snap/bin/chromium, or a script that runs one (Ubuntu's
// /usr/bin/chromium-browser); "" for any other browser.
func snapName(browser string) string {
	if runtime.GOOS != "linux" {
		return ""
	}
	text := browser
	if !strings.HasPrefix(browser, "/snap/bin/") {
		file, err := os.Open(browser)
		if err != nil {
			return ""
		}
		defer file.Close()
		head := make([]byte, 4096) // a script is small; a browser itself is a program, not read further
		count, _ := file.Read(head)
		if count < 2 || string(head[:2]) != "#!" {
			return ""
		}
		text = string(head[:count])
	}
	at := strings.Index(text, "/snap/bin/")
	if at < 0 {
		return ""
	}
	name := text[at+len("/snap/bin/"):]
	end := strings.IndexFunc(name, func(letter rune) bool {
		return !(letter >= 'a' && letter <= 'z' || letter >= '0' && letter <= '9' || letter == '-')
	})
	if end >= 0 {
		name = name[:end]
	}
	return name
}

// A browser of this computer that opens an app window (Chromium kind): Edge is on every Windows.
func findBrowser() string {
	if given := os.Getenv("OZDENUI_BROWSER"); given != "" {
		return given
	}
	candidates := []string{}
	switch runtime.GOOS {
	case "darwin":
		for _, name := range []string{"Google Chrome", "Microsoft Edge", "Chromium", "Brave Browser"} {
			candidates = append(candidates, "/Applications/"+name+".app/Contents/MacOS/"+name)
			candidates = append(candidates, filepath.Join(os.Getenv("HOME"), "Applications", name+".app/Contents/MacOS", name))
		}
	case "windows":
		for _, root := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"), os.Getenv("LocalAppData")} {
			if root == "" {
				continue
			}
			candidates = append(candidates, filepath.Join(root, `Microsoft\Edge\Application\msedge.exe`))
			candidates = append(candidates, filepath.Join(root, `Google\Chrome\Application\chrome.exe`))
		}
	default:
		for _, name := range []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "microsoft-edge", "brave-browser"} {
			if path, err := exec.LookPath(name); err == nil {
				return path
			}
		}
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func openInDefaultBrowser(address string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", address).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", address).Start()
	}
	return exec.Command("xdg-open", address).Start()
}

// The title of the page: the title of a browser window, and what the system lists.
func (app *application) setPageTitle(title string) {
	text, _ := json.Marshal(title)
	app.tasks.put(func() {
		_, _ = app.server.request(map[string]any{"kind": "eval", "script": "document.title = " + string(text) + "; 0"})
	})
}
