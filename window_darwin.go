//go:build cgo && darwin && !ozdenui_browser

/* Bismillah */

/*

window_darwin.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
macOS: the menu of the program.
WHY: a program without a menu has no Cmd+C / Cmd+V / Cmd+A in its text fields and no Cmd+Q: the system sends these
keys through the Edit menu and the menu of the application.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#include <stdlib.h>

static void ozdenuiAddMenu(const char *title) {
    NSString *name = [NSString stringWithUTF8String:title];
    NSMenu *bar = [[NSMenu alloc] init];

    NSMenuItem *appItem = [[NSMenuItem alloc] init];
    NSMenu *appMenu = [[NSMenu alloc] initWithTitle:name];
    [appMenu addItemWithTitle:[@"Hide " stringByAppendingString:name] action:@selector(hide:) keyEquivalent:@"h"];
    [appMenu addItem:[NSMenuItem separatorItem]];
    [appMenu addItemWithTitle:[@"Quit " stringByAppendingString:name] action:@selector(terminate:) keyEquivalent:@"q"];
    [appItem setSubmenu:appMenu];
    [bar addItem:appItem];

    NSMenuItem *editItem = [[NSMenuItem alloc] init];
    NSMenu *editMenu = [[NSMenu alloc] initWithTitle:@"Edit"];
    [editMenu addItemWithTitle:@"Undo" action:@selector(undo:) keyEquivalent:@"z"];
    [editMenu addItemWithTitle:@"Redo" action:@selector(redo:) keyEquivalent:@"Z"];
    [editMenu addItem:[NSMenuItem separatorItem]];
    [editMenu addItemWithTitle:@"Cut" action:@selector(cut:) keyEquivalent:@"x"];
    [editMenu addItemWithTitle:@"Copy" action:@selector(copy:) keyEquivalent:@"c"];
    [editMenu addItemWithTitle:@"Paste" action:@selector(paste:) keyEquivalent:@"v"];
    [editMenu addItemWithTitle:@"Select All" action:@selector(selectAll:) keyEquivalent:@"a"];
    [editItem setSubmenu:editMenu];
    [bar addItem:editItem];

    NSMenuItem *windowItem = [[NSMenuItem alloc] init];
    NSMenu *windowMenu = [[NSMenu alloc] initWithTitle:@"Window"];
    [windowMenu addItemWithTitle:@"Minimize" action:@selector(performMiniaturize:) keyEquivalent:@"m"];
    [windowMenu addItemWithTitle:@"Close" action:@selector(performClose:) keyEquivalent:@"w"];
    [windowItem setSubmenu:windowMenu];
    [bar addItem:windowItem];

    [NSApp setMainMenu:bar];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    [NSApp activateIgnoringOtherApps:YES];
}

// WHY: the size given to the web view is the size of the window with its title bar; the size of a program is the
// size of its page (Window{Width, Height}), as its objects are placed in it.
static void ozdenuiSetPageSize(int width, int height) {
    NSWindow *window = [[NSApp windows] firstObject];
    if (window == nil) return;
    [window setContentSize:NSMakeSize(width, height)];
    [window center];
}
*/
import "C"

import "unsafe"

func setPageSize(width, height int) { C.ozdenuiSetPageSize(C.int(width), C.int(height)) }

func prepareNativeWindow(title string) {
	text := C.CString(title)
	defer C.free(unsafe.Pointer(text))
	C.ozdenuiAddMenu(text)
}
