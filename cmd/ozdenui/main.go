/* Bismillah */

/*

ozdenui - ozdenui v0.1 (OzdenUI: basic.js for Go)
The command of the module. It writes the rules of OzdenUI for an AI assistant (AGENTS.md, next to this file) into the
folder of a project, where the assistants read them by themselves.

  go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents              Into AGENTS.md of this folder (most assistants read it).
  go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents CLAUDE.md    Into another file (Claude Code reads CLAUDE.md).
  go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents --print      Only prints them.

A file that is there keeps its own lines: the rules go between two marks at its end, and the next run changes only them.
(A picture of the window of a program needs no command: OZDENUI_LOOK=app.png go run .)

WHY it does not import the module: the command must build in a second, without the web view and its C compiler.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed AGENTS.md
var rules string

const (
	version       = "0.2.1" // the Version of the module (its test compares them)
	defaultTarget = "AGENTS.md"
	command       = "go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest"
	startMark     = "<!-- ozdenui-go: start (written by the agents command of ozdenui-go; the next run writes this part again) -->"
	endMark       = "<!-- ozdenui-go: end -->"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(words []string) int {
	if len(words) == 0 || words[0] != "agents" || len(words) > 2 {
		fmt.Printf("ozdenui-go %s\n\n%s agents              The rules of OzdenUI for an AI assistant, into AGENTS.md of this folder.\n"+
			"%s agents CLAUDE.md    Into another file.\n%s agents --print      Only prints them.\n\n"+
			"OZDENUI_LOOK=app.png go run .    Runs your program, saves a picture of its window as app.png and ends.\n",
			version, command, command, command)
		if len(words) == 0 {
			return 0
		}
		return 2
	}
	if len(words) == 2 && words[1] == "--print" {
		fmt.Println(strings.TrimSpace(rules))
		return 0
	}
	target := defaultTarget
	if len(words) == 2 {
		target = words[1]
	}
	done, err := writeRules(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	messages := map[string]string{
		"made":     "%s: made, with the rules of OzdenUI for Go %s.\n",
		"added":    "%s: the rules of OzdenUI for Go %s are added at its end; your lines are as they were.\n",
		"updated":  "%s: the rules of OzdenUI for Go are updated to %s; your lines are as they were.\n",
		"the same": "%s: already has the rules of OzdenUI for Go %s.\n",
	}
	fmt.Printf(messages[done], target, version)
	return 0
}

// The text of the file with the rules between the marks: in the place of the old ones, or at its end.
func withRules(oldText string) string {
	block := startMark + "\n" + strings.TrimSpace(rules) + "\n" + endMark + "\n"
	start := strings.Index(oldText, startMark)
	end := strings.Index(oldText, endMark)
	if start != -1 && end > start {
		return oldText[:start] + block + strings.TrimLeft(oldText[end+len(endMark):], "\n")
	}
	if strings.TrimSpace(oldText) == "" {
		return block
	}
	return strings.TrimRight(oldText, "\n") + "\n\n" + block
}

// Writes the rules into the file; it returns what it did: "made", "added", "updated" or "the same".
func writeRules(target string) (string, error) {
	content, err := os.ReadFile(target)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	oldText := string(content)
	newText := withRules(oldText)
	if newText == oldText {
		return "the same", nil
	}
	if err := os.WriteFile(target, []byte(newText), 0o644); err != nil {
		return "", err
	}
	switch {
	case strings.TrimSpace(oldText) == "":
		return "made", nil
	case strings.Contains(oldText, startMark):
		return "updated", nil
	}
	return "added", nil
}
