//go:build cgo && (darwin || ozdenui_webview) && !ozdenui_browser

/* Bismillah */

/*

window_webview.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The window as the web view of the system (github.com/webview/webview_go): WKWebView on macOS.
On Windows (WebView2) and Linux (WebKitGTK) it needs a C compiler and their headers: go build -tags ozdenui_webview

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	webview "github.com/webview/webview_go"
)

const maxTitleBar = 120 // px: a larger difference is not a title bar (a small screen cut the window)

type webviewWindow struct {
	view webview.WebView
}

func (w *webviewWindow) setTitle(title string) {
	w.view.Dispatch(func() { w.view.SetTitle(title) })
}

func (w *webviewWindow) close() { w.view.Terminate() }

// WHY: the size given to the web view is the size of the window with its title bar; the size of a program is the
// size of its page (Window{Width, Height}), as its objects are placed in it.
func (w *webviewWindow) fitPage(wantedWidth, wantedHeight int, pageWidth, pageHeight float64) {
	extraWidth, extraHeight := wantedWidth-int(pageWidth), wantedHeight-int(pageHeight)
	if extraWidth == 0 && extraHeight == 0 || extraWidth < 0 || extraHeight < 0 || extraHeight > maxTitleBar {
		return
	}
	w.view.Dispatch(func() { w.view.SetSize(wantedWidth+extraWidth, wantedHeight+extraHeight, webview.HintNone) })
}

func openNativeWindow(app *application, address string) bool {
	view := webview.New(app.window.Debug)
	if view == nil {
		return false // no web view on this computer: the browser window
	}
	defer view.Destroy()
	app.control = &webviewWindow{view: view}
	view.SetTitle(app.window.Title)
	view.SetSize(app.window.Width, app.window.Height, webview.HintNone)
	setPageSize(app.window.Width, app.window.Height)
	// WHY before the menu: the page loads while the menu is made and the program comes to the front
	view.Navigate(address)
	prepareNativeWindow(app.window.Title)
	view.Run()
	return true
}
