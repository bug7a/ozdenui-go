//go:build cgo && ozdenui_webview && !darwin && !ozdenui_browser

/* Bismillah */

/*

window_other.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
Windows and Linux with the web view of the system: nothing to prepare (the keys of the text fields work).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

func prepareNativeWindow(title string) {}

// The page tells its size when it is ready, and the window is fitted then (fitPage).
func setPageSize(width, height int) {}
