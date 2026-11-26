//go:build !(cgo && (darwin || ozdenui_webview) && !ozdenui_browser) && !(windows && !ozdenui_browser) && !(linux && !ozdenui_browser)

/* Bismillah */

/*

window_nowebview.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
A build without the web view of the system (macOS without cgo, the tag ozdenui_browser, a system that is not macOS,
Windows or Linux): the window is the app window of a browser (window.go). Windows: window_windows.go, Linux:
window_linux.go.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

func openNativeWindow(app *application, address string) bool { return false }
