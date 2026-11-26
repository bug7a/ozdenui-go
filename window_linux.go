//go:build linux && !ozdenui_browser && !(cgo && ozdenui_webview)

/* Bismillah */

/*

window_linux.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
Linux: the web view of the system (WebKitGTK 4.1, or 4.0), in the program itself, without a C compiler: its library
is loaded when the program runs (github.com/ebitengine/purego). Without it, or without a screen: the app window of
Chrome, Chromium or Edge (window.go).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"sync"

	"github.com/ebitengine/purego"
)

const (
	maxTitleBar       = 120 // px: a larger difference is not a title bar (a small screen cut the window)
	gtkWindowToplevel = 0
	gtkWindowCenter   = 1
)

// The libraries of WebKitGTK, the newest first (Ubuntu 24.04 and newer have only 4.1). GTK, GObject and GLib are found
// through them.
var webkitLibraries = []string{"libwebkit2gtk-4.1.so.0", "libwebkit2gtk-4.0.so.37"}

var (
	gtkInitCheck             func(argc, argv uintptr) int32
	gtkWindowNew             func(kind int32) uintptr
	gtkWindowSetTitle        func(window uintptr, title string)
	gtkWindowSetDefaultSize  func(window uintptr, width, height int32)
	gtkWindowResize          func(window uintptr, width, height int32)
	gtkWindowSetPosition     func(window uintptr, position int32)
	gtkContainerAdd          func(container, widget uintptr)
	gtkWidgetShowAll         func(widget uintptr)
	gtkWidgetDestroy         func(widget uintptr)
	gtkMain                  func()
	webkitWebViewNew         func() uintptr
	webkitWebViewLoadURI     func(view uintptr, uri string)
	webkitWebViewGetSettings func(view uintptr) uintptr
	webkitSetDeveloperExtras func(settings uintptr, enabled int32)
	gSignalConnectData       func(instance uintptr, signal string, handler, data, destroyData uintptr, flags int32) uint64
	gIdleAdd                 func(function, data uintptr) uint32
)

var (
	gtkMainQuit uintptr // the C functions given to signals as they are
	gtkTrue     uintptr
	loadOnce    sync.Once
	loaded      bool
)

// Loads WebKitGTK and its functions; false when this computer has none.
func loadWebKitGTK() bool {
	loadOnce.Do(func() {
		library := uintptr(0)
		for _, name := range webkitLibraries {
			if handle, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL); err == nil {
				library = handle
				break
			}
		}
		if library == 0 {
			return
		}
		functions := map[string]any{
			"gtk_init_check":                              &gtkInitCheck,
			"gtk_window_new":                              &gtkWindowNew,
			"gtk_window_set_title":                        &gtkWindowSetTitle,
			"gtk_window_set_default_size":                 &gtkWindowSetDefaultSize,
			"gtk_window_resize":                           &gtkWindowResize,
			"gtk_window_set_position":                     &gtkWindowSetPosition,
			"gtk_container_add":                           &gtkContainerAdd,
			"gtk_widget_show_all":                         &gtkWidgetShowAll,
			"gtk_widget_destroy":                          &gtkWidgetDestroy,
			"gtk_main":                                    &gtkMain,
			"webkit_web_view_new":                         &webkitWebViewNew,
			"webkit_web_view_load_uri":                    &webkitWebViewLoadURI,
			"webkit_web_view_get_settings":                &webkitWebViewGetSettings,
			"webkit_settings_set_enable_developer_extras": &webkitSetDeveloperExtras,
			"g_signal_connect_data":                       &gSignalConnectData,
			"g_idle_add":                                  &gIdleAdd,
		}
		for name, function := range functions {
			symbol, err := purego.Dlsym(library, name)
			if err != nil {
				return // WHY not a panic: an old or broken library is only the browser window
			}
			purego.RegisterFunc(function, symbol)
		}
		var errQuit, errTrue error
		gtkMainQuit, errQuit = purego.Dlsym(library, "gtk_main_quit")
		gtkTrue, errTrue = purego.Dlsym(library, "gtk_true")
		loaded = errQuit == nil && errTrue == nil
	})
	return loaded
}

// *** THE TASKS OF THE WINDOW: GTK answers only on its own thread (the main thread of the program).

var (
	tasksMu      sync.Mutex
	windowTasks  []func()
	idleCallback uintptr
)

func runWindowTasks(_ uintptr) uintptr {
	tasksMu.Lock()
	tasks := windowTasks
	windowTasks = nil
	tasksMu.Unlock()
	for _, task := range tasks {
		task()
	}
	return 0 // G_SOURCE_REMOVE: one call for each g_idle_add
}

// Runs task on the thread of the window; Close() and SetTitle() come from the tasks of the program.
func onWindowThread(task func()) {
	tasksMu.Lock()
	windowTasks = append(windowTasks, task)
	tasksMu.Unlock()
	gIdleAdd(idleCallback, 0)
}

type gtkWindow struct {
	window uintptr
}

func (w *gtkWindow) setTitle(title string) {
	onWindowThread(func() { gtkWindowSetTitle(w.window, title) })
}

func (w *gtkWindow) close() {
	onWindowThread(func() { gtkWidgetDestroy(w.window) })
}

// WHY: the size given to the window may hold its title bar (on some desktops); the size of a program is the size of
// its page (Window{Width, Height}), as its objects are placed in it.
func (w *gtkWindow) fitPage(wantedWidth, wantedHeight int, pageWidth, pageHeight float64) {
	extraWidth, extraHeight := wantedWidth-int(pageWidth+0.5), wantedHeight-int(pageHeight+0.5)
	if extraWidth == 0 && extraHeight == 0 || extraWidth < 0 || extraHeight < 0 || extraHeight > maxTitleBar {
		return
	}
	onWindowThread(func() {
		gtkWindowResize(w.window, int32(wantedWidth+extraWidth), int32(wantedHeight+extraHeight))
	})
}

func openNativeWindow(app *application, address string) bool {
	if !loadWebKitGTK() || gtkInitCheck(0, 0) == 0 {
		return false // no WebKitGTK, or no screen: the browser window
	}
	if idleCallback == 0 {
		idleCallback = purego.NewCallback(runWindowTasks)
	}
	window := gtkWindowNew(gtkWindowToplevel)
	gtkWindowSetTitle(window, app.window.Title)
	gtkWindowSetDefaultSize(window, int32(app.window.Width), int32(app.window.Height))
	gtkWindowSetPosition(window, gtkWindowCenter)
	view := webkitWebViewNew()
	if app.window.Debug {
		webkitSetDeveloperExtras(webkitWebViewGetSettings(view), 1) // right click > Inspect Element
	} else {
		// WHY: the menu of the web view (Back, Reload) is not a part of a program; gtk_true says "handled"
		gSignalConnectData(view, "context-menu", gtkTrue, 0, 0, 0)
	}
	gtkContainerAdd(window, view)
	gSignalConnectData(window, "destroy", gtkMainQuit, 0, 0, 0)
	app.control = &gtkWindow{window: window}
	webkitWebViewLoadURI(view, address)
	gtkWidgetShowAll(window)
	gtkMain()
	return true
}
