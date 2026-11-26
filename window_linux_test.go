//go:build linux && !ozdenui_browser && !(cgo && ozdenui_webview)

/* Bismillah */

/*

window_linux_test.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The window of Linux: a computer without WebKitGTK gets the browser window, not an error.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"sync"
	"testing"
)

func TestLinuxWithoutWebKitGTK(t *testing.T) {
	oldLibraries := webkitLibraries
	defer func() {
		webkitLibraries = oldLibraries
		loadOnce, loaded = sync.Once{}, false
	}()
	webkitLibraries = []string{"libwebkit2gtk-none.so.0"}
	loadOnce, loaded = sync.Once{}, false
	if loadWebKitGTK() {
		t.Error("a missing WebKitGTK was loaded")
	}
	if openNativeWindow(&application{}, "http://127.0.0.1:1/") {
		t.Error("the window of the system opened without WebKitGTK: the browser window must open")
	}
}
