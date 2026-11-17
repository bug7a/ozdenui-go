//go:build !windows

/* Bismillah */

/*

dialogs_other.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
macOS and Linux: a dialog command needs nothing more (dialogs_windows.go hides the console of Windows).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import "os/exec"

func withoutConsole(command *exec.Cmd) *exec.Cmd { return command }
