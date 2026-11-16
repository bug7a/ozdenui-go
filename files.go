/* Bismillah */

/*

files.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The files of a FilePicker with their paths on this computer.

The dialog of the page gives a file without its path, so a click on a FilePicker opens the file dialog of the system
(unless the program gave an OnOpenDialog of its own), and the chosen files go into its list with their real paths.
A file dropped on it still comes from the page: its name, not its path.

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxPreviews     = 12              // more files than this: no pictures in the list (each one is read)
	maxPreviewBytes = 2 * 1024 * 1024 // a larger picture is not read for the list
)

// The endings behind a kind of files in the accept of a FilePicker ("image/*").
var endingsOfKind = map[string][]string{
	"image/*": {"png", "jpg", "jpeg", "gif", "webp", "svg", "bmp"},
	"audio/*": {"mp3", "wav", "ogg", "m4a", "flac"},
	"video/*": {"mp4", "mov", "webm", "mkv", "avi"},
	"text/*":  {"txt", "md", "csv", "json", "html", "css", "js"},
}

// AddFiles adds files of this computer to the list of a FilePicker, with their real paths. Its Accept, MaxSize and
// MaxFiles look at them too, and its OnChange runs. A file that can not be read is told and left out.
//
//	picker.AddFiles("/Users/ada/report.pdf", "notes.txt")
func (obj *Object) AddFiles(paths ...string) *Object {
	obj.mustBe("FilePicker", "AddFiles")
	showPreview, isGiven := obj.lookup("showPreview")
	withPreview := len(paths) <= maxPreviews && (!isGiven || isOn(showPreview))
	files := make([]P, 0, len(paths))
	for _, path := range paths {
		if file := fileForPage(path, withPreview); file != nil {
			files = append(files, file)
		} else {
			obj.app.warn("FilePicker.AddFiles: %s could not be read", path)
		}
	}
	return obj.Call("addFiles", files)
}

// Paths are the files of a FilePicker that have a path on this computer: the ones chosen in the file dialog or
// given by AddFiles. A file dropped on the window has only its name (picker.Get("files") lists every file).
func (obj *Object) Paths() []string {
	obj.mustBe("FilePicker", "Paths")
	paths := []string{}
	files, _ := obj.lookup("files")
	list, _ := files.([]any)
	for _, item := range list {
		file, _ := item.(map[string]any)
		if path := textOf(file["path"]); filepath.IsAbs(path) {
			paths = append(paths, path)
		}
	}
	return paths
}

// What the page needs of a file: a plain object in the place of the File object of a browser. nil: not a file.
func fileForPage(path string, withPreview bool) P {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return nil
	}
	kind := mime.TypeByExtension(strings.ToLower(filepath.Ext(path)))
	if cut := strings.Index(kind, ";"); cut >= 0 {
		kind = kind[:cut] // "text/plain; charset=utf-8"
	}
	file := P{"name": filepath.Base(path), "path": path, "size": info.Size(), "type": kind, "lastModified": info.ModTime().UnixMilli()}
	if withPreview && strings.HasPrefix(kind, "image/") && info.Size() <= maxPreviewBytes {
		file["previewAddress"] = pictureOfPath(path)
	}
	return file
}

// The OnOpenDialog of a FilePicker that was given none: the file dialog of the system.
func openDialogOfPicker(picker *Object) {
	if directory, _ := picker.lookup("directory"); isOn(directory) {
		if folder := ChooseFolder(); folder != "" {
			picker.addFolder(folder)
		}
		return
	}
	accept, _ := picker.lookup("accept")
	options := FileOptions{Types: fileTypesOf(textOf(accept))}
	if multiple, _ := picker.lookup("multiple"); isOn(multiple) {
		if paths := OpenFiles(options); len(paths) > 0 {
			picker.AddFiles(paths...)
		}
		return
	}
	if path := OpenFile(options); path != "" {
		picker.AddFiles(path)
	}
}

// Every file of a folder, as the dialog of the page gives a folder: a row says where its file is.
func (obj *Object) addFolder(folder string) {
	files := []P{}
	_ = filepath.WalkDir(folder, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		file := fileForPage(path, false)
		if file == nil {
			return nil
		}
		if relative, err := filepath.Rel(filepath.Dir(folder), path); err == nil {
			file["webkitRelativePath"] = filepath.ToSlash(relative)
		}
		files = append(files, file)
		return nil
	})
	obj.Call("addFiles", files)
}

// The endings of the accept of a FilePicker (".pdf,.docx", "image/*", "image/png"), for the dialog of the system.
// nil: every file.
func fileTypesOf(accept string) []string {
	var endings []string
	for _, part := range strings.Split(accept, ",") {
		part = strings.ToLower(strings.TrimSpace(part))
		switch {
		case part == "":
		case strings.HasPrefix(part, "."):
			endings = append(endings, part[1:])
		case endingsOfKind[part] != nil:
			endings = append(endings, endingsOfKind[part]...)
		default:
			found, _ := mime.ExtensionsByType(part)
			if len(found) == 0 {
				return nil // WHY: a kind nobody knows the endings of: the dialog must not hide its files
			}
			for _, ending := range found {
				endings = append(endings, strings.TrimPrefix(ending, "."))
			}
		}
	}
	return endings
}
