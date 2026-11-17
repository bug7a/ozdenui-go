//go:build windows && !ozdenui_browser && !(cgo && ozdenui_webview)

/* Bismillah */

/*

window_windows.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
Windows: the web view of the system (WebView2, in every Windows 11 and in most Windows 10), in the program itself,
without a C compiler (github.com/jchv/go-webview2). Without WebView2: the app window of Edge or Chrome (window.go).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"path/filepath"
	"syscall"

	webview2 "github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/webviewloader"
)

const (
	baseDpi = 96 // the dpi of a scale of 100 %
	// WHY system aware: the window and its text are sharp on a screen of 125 % or 150 % (a window that does not know
	//     the scale is drawn small and stretched, so it is blurred). Per monitor would need the window to follow a
	//     move to another screen, which the web view package does not do.
	dpiAwarenessSystemAware = ^uintptr(1) // DPI_AWARENESS_CONTEXT_SYSTEM_AWARE (-2)
)

var (
	user32                        = syscall.NewLazyDLL("user32.dll")
	setProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	getDpiForWindow               = user32.NewProc("GetDpiForWindow")
)

type webview2Window struct {
	view  webview2.WebView
	scale float64
}

// WHY Dispatch: the window answers only on its own thread; Close() and SetTitle() come from the tasks of the program.
func (w *webview2Window) setTitle(title string) {
	w.view.Dispatch(func() { w.view.SetTitle(title) })
}

func (w *webview2Window) close() {
	w.view.Dispatch(func() { w.view.Terminate() })
}

// The window is given the size of its page (SetSize gives the inside of the window); a page that tells another size
// (a scale the window did not know) is fitted once more.
func (w *webview2Window) fitPage(wantedWidth, wantedHeight int, pageWidth, pageHeight float64) {
	if int(pageWidth+0.5) == wantedWidth && int(pageHeight+0.5) == wantedHeight || pageWidth <= 0 || pageHeight <= 0 {
		return
	}
	w.view.Dispatch(func() {
		w.view.SetSize(int(float64(wantedWidth)*w.scale+0.5), int(float64(wantedHeight)*w.scale+0.5), webview2.HintNone)
	})
}

func openNativeWindow(app *application, address string) bool {
	if version, err := webviewloader.GetInstalledVersion(); err != nil || version == "" {
		return false // no WebView2 on this computer: the browser window
	}
	if setProcessDpiAwarenessContext.Find() == nil {
		_, _, _ = setProcessDpiAwarenessContext.Call(dpiAwarenessSystemAware)
	}
	view := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug: app.window.Debug,
		// WHY a folder of its own: by default the web view keeps its data in the roaming AppData, next to the name
		//     of the program, which is a temporary file under go run
		DataPath:      filepath.Join(cacheFolder(), "webview2-"+programName()),
		AutoFocus:     true,
		WindowOptions: webview2.WindowOptions{Title: app.window.Title, Width: uint(app.window.Width), Height: uint(app.window.Height), Center: true},
	})
	if view == nil {
		return false
	}
	defer view.Destroy()
	window := &webview2Window{view: view, scale: windowScale(view)}
	app.control = window
	view.SetSize(int(float64(app.window.Width)*window.scale+0.5), int(float64(app.window.Height)*window.scale+0.5), webview2.HintNone)
	view.Navigate(address)
	view.Run()
	return true
}

// The scale of the screen of the window (1.5 for 150 %); 1 on a Windows older than 10 (1607).
func windowScale(view webview2.WebView) float64 {
	if getDpiForWindow.Find() != nil {
		return 1
	}
	dpi, _, _ := getDpiForWindow.Call(uintptr(view.Window()))
	if dpi == 0 {
		return 1
	}
	return float64(dpi) / baseDpi
}
