#!/bin/sh
# Bismillah
#
# handbook.sh - ozdenui v0.1 (OzdenUI: basic.js for Go)
# Every whole program of the handbook (a Go block with "package main" in handbook/*.md and handbook/turkce/*.md) is
# built and drawn in a hidden browser (OZDENUI_LOOK); one that does not build, fails, or tells a wrong name
# ("ozdenui: ..." in its terminal) fails the check. So the handbook stays true when the module changes.
# Usage: checks/handbook.sh [folder of the pictures (default: pictures/handbook)] [a part of a chapter name (Ex: 03, turkce/05)]
# The pictures are named <chapter>-<number of the program in it>.png: look at them after writing a chapter.
#
# Started Date: November 2026
# Webpage: https://bug7a.github.io/ozdenui-go/
#
# Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
# Licensed under the Apache License, Version 2.0
# SPDX-License-Identifier: Apache-2.0

cd "$(dirname "$0")/.." || exit 1
PICTURES="${1:-pictures/handbook}"
ONLY="${2:-}"
# WHY a folder that starts with "_": the go tool leaves it out of ./... (the programs are not a part of the module)
WORK="handbook/_check-$$"  # its own folder: two checks may run at the same time
rm -rf "$WORK"
mkdir -p "$WORK" "$PICTURES"
case "$PICTURES" in /*) ;; *) PICTURES="$(pwd)/$PICTURES" ;; esac

python3 - "$WORK" "$ONLY" <<'PY'
import pathlib, re, sys
work, only = pathlib.Path(sys.argv[1]), sys.argv[2]
root = pathlib.Path("handbook")
for chapter in sorted(list(root.glob("*.md")) + list(root.glob("turkce/*.md"))):
    name = chapter.relative_to(root).with_suffix("").as_posix()
    if only and only not in name:
        continue
    blocks = re.findall(r"```go\n(.*?)```", chapter.read_text(encoding="utf-8"), re.S)
    for number, code in enumerate([block for block in blocks if "package main" in block], start=1):
        folder = work / (name.replace("/", "-") + "-" + str(number))
        folder.mkdir(parents=True)
        (folder / "main.go").write_text(code, encoding="utf-8")
PY

FAILED=0
COUNT=0
for folder in "$WORK"/*/; do
    [ -d "$folder" ] || continue
    name="$(basename "$folder")"
    COUNT=$((COUNT + 1))
    picture="$PICTURES/$name.png"
    rm -f "$picture"
    unformatted="$(gofmt -l "$folder")"
    output="$(OZDENUI_LOOK_WAIT=0.6 OZDENUI_LOOK="$picture" go run "./$folder" 2>&1 | grep -v "the picture of the window" | grep "^ozdenui\|panic\|^exit status\|cannot\|undefined\|declared and not used\|imported and not used\|syntax error\|main.go:")"
    if [ ! -s "$picture" ] || [ -n "$output" ] || [ -n "$unformatted" ]; then
        FAILED=1
        echo "FAIL $name"
        [ -n "$unformatted" ] && echo "     not as gofmt writes it (tabs, spacing)"
        echo "$output" | sed 's/^/     /'
    else
        echo "ok   $name"
    fi
done
rm -rf "$WORK"
[ "$FAILED" = "0" ] && echo "ALL PASS ($COUNT programs)" || echo "SOME FAILED"
exit $FAILED
