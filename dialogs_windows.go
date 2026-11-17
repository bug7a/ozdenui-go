//go:build windows

/* Bismillah */

/*

dialogs_windows.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
Windows: the PowerShell of a file dialog starts without a console window.
WHY: a program built without a console (-ldflags -H=windowsgui, as a program given to others is) got a new console
window behind the file dialog, an empty black terminal.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"os/exec"
	"syscall"
)

const createNoWindow = 0x08000000 // CREATE_NO_WINDOW

func withoutConsole(command *exec.Cmd) *exec.Cmd {
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
	return command
}
