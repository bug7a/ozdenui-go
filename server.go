/* Bismillah */

/*

server.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
The page of the window and its messages, served on this computer only (127.0.0.1, a port of its own).

The files of the page (web/) are inside the program (go:embed): a built program is one file.
The program sends its messages over an event stream ("commands"); the page answers and tells its events with POST
"emit", one request at a time, so they come in the order they happened (bridge.js).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"net"
	"net/http"
	"sync"
)

//go:embed web
var webFiles embed.FS

var errPageGone = errors.New("the page is gone")

type server struct {
	app      *application
	listener net.Listener
	token    string // WHY: another program or page of this computer must not reach the page of the window

	mu        sync.Mutex
	pageID    string                    // the page that listens now (a page loaded again is a new one)
	outgoing  chan []byte               // the messages for it
	waiting   map[int]chan replyMessage // the answers the program waits for
	nextReply int
	emitMu    sync.Mutex        // one batch of the page at a time, in order
	onStream  func(isOpen bool) // the window of a browser tab is known to be closed by its stream (window_browser.go)
}

type replyMessage struct {
	value json.RawMessage
	err   error
}

type pageMessage struct {
	Ready  *readyMessage   `json:"ready"`
	Reply  *int            `json:"reply"`
	Value  json.RawMessage `json:"value"`
	Error  *string         `json:"error"`
	Target string          `json:"target"`
	Event  string          `json:"event"`
	Data   json.RawMessage `json:"data"`
}

func startServer(app *application) (*server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return nil, err
	}
	s := &server{
		app:      app,
		listener: listener,
		token:    hex.EncodeToString(random),
		outgoing: make(chan []byte, 256),
		waiting:  map[int]chan replyMessage{},
	}
	web, _ := fs.Sub(webFiles, "web")
	files := http.FileServer(http.FS(web))
	prefix := "/" + s.token + "/"
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"commands", s.serveCommands)
	mux.HandleFunc(prefix+"emit", s.serveEmit)
	mux.HandleFunc(prefix+"index.html", s.serveIndex)
	mux.Handle(prefix, http.StripPrefix(prefix, noCache(files)))
	go func() { _ = http.Serve(listener, mux) }()
	return s, nil
}

// WHY: a web view keeps a file of an older version of the module under the same address of another run
func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (s *server) address() string {
	return fmt.Sprintf("http://%s/%s/index.html", s.listener.Addr().String(), s.token)
}

func (s *server) stop() {
	s.failWaiting()
	_ = s.listener.Close()
}

// The page, with the title of the window (a browser window shows the title of its page).
func (s *server) serveIndex(w http.ResponseWriter, r *http.Request) {
	page, err := webFiles.ReadFile("web/index.html")
	if err != nil {
		http.Error(w, "no page", http.StatusInternalServerError)
		return
	}
	title := "<title>" + html.EscapeString(s.app.window.Title) + "</title>"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	// The size of the page the program asked for: a browser window fits itself to it before the program starts (bridge.js)
	size := fmt.Sprintf(`<meta name="ozdenui-size" content="%dx%d">`, s.app.window.Width, s.app.window.Height)
	page = bytes.Replace(page, []byte("<title>ozdenui</title>"), []byte(title+"\n        "+size), 1)
	if language := s.app.window.Language; language != "" {
		page = bytes.Replace(page, []byte(`<html lang="en">`), []byte(`<html lang="`+html.EscapeString(language)+`">`), 1)
	}
	_, _ = w.Write(page)
}

// The event stream of the page: every message of the program is one "data:" line.
func (s *server) serveCommands(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "no stream", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	s.mu.Lock()
	outgoing := s.channelOf(r.URL.Query().Get("page"))
	onStream := s.onStream
	s.mu.Unlock()
	if onStream != nil {
		onStream(true)
		defer onStream(false)
	}
	for {
		select {
		case <-r.Context().Done():
			return
		case message := <-outgoing:
			if _, err := fmt.Fprintf(w, "data: %s\n\n", message); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (s *server) serveEmit(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || r.Method != http.MethodPost {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var batch []pageMessage
	if err := json.Unmarshal(body, &batch); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.emitMu.Lock()
	for _, message := range batch {
		s.take(message)
	}
	s.emitMu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) take(message pageMessage) {
	switch {
	case message.Ready != nil:
		// A new page: what waited for the old one gets no answer, and its messages are not for the new one.
		s.mu.Lock()
		s.channelOf(message.Ready.Page)
		s.mu.Unlock()
		s.failWaiting()
		s.app.onPageReady(*message.Ready)
	case message.Reply != nil:
		s.mu.Lock()
		waiter := s.waiting[*message.Reply]
		delete(s.waiting, *message.Reply)
		s.mu.Unlock()
		if waiter == nil {
			return
		}
		if message.Error != nil {
			waiter <- replyMessage{err: errors.New(*message.Error)}
			return
		}
		waiter <- replyMessage{value: message.Value}
	default:
		s.app.onPageEvent(message.Target, message.Event, message.Data)
	}
}

// The messages of a page, by its id (the lock is held). WHY by its id: the stream of a page may open before or after
// the server took its "ready"; a channel made new at "ready" was sometimes not the one its stream listened to, and
// the program waited for ever for an answer.
func (s *server) channelOf(pageID string) chan []byte {
	if pageID != s.pageID {
		s.pageID = pageID
		s.outgoing = make(chan []byte, 256)
	}
	return s.outgoing
}

func (s *server) failWaiting() {
	s.mu.Lock()
	waiting := s.waiting
	s.waiting = map[int]chan replyMessage{}
	s.mu.Unlock()
	for _, waiter := range waiting {
		waiter <- replyMessage{err: errPageGone}
	}
}

// Sends a message to the page and waits for its answer.
func (s *server) request(message map[string]any) (json.RawMessage, error) {
	waiter := make(chan replyMessage, 1)
	s.mu.Lock()
	s.nextReply++
	number := s.nextReply
	s.waiting[number] = waiter
	outgoing := s.outgoing
	s.mu.Unlock()
	message["n"] = number
	text, err := json.Marshal(message)
	if err != nil {
		s.mu.Lock()
		delete(s.waiting, number)
		s.mu.Unlock()
		return nil, fmt.Errorf("a value can not be sent to the page: %w", err)
	}
	select {
	case outgoing <- text:
	case <-s.app.closed:
		return nil, errPageGone
	}
	select {
	case reply := <-waiter:
		return reply.value, reply.err
	case <-s.app.closed:
		return nil, errPageGone
	}
}
