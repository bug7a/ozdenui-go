/* Bismillah */

/*

look.go - ozdenui v0.1 (OzdenUI: basic.js for Go)
A picture of the window of any program, without a window: the page is drawn by a hidden Chrome (or Edge).

  OZDENUI_LOOK=app.png go run .          the program runs as it is; its window is saved as app.png, and it ends
  OZDENUI_LOOK_WAIT=3                    seconds after the first drawing before the picture (default 1.5)
  OZDENUI_LOOK_SCALE=2                   the pixels of a point of the picture (default 1)
  OZDENUI_THEME=dark                     the dark window, without a change of the program
  OZDENUI_WINDOW=headless go run .       no picture: the program runs in the hidden browser until it calls Close()

WHY: a program that runs without an error can still have a part covered, cut or squeezed. An AI assistant that
writes a program can not see a window; it can read a picture.

The browser is driven over its DevTools protocol, with a small WebSocket client of its own (no module needed).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

package ozdenui

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultLookWait  = 1.5
	browserStartWait = 20 * time.Second
)

type headlessWindow struct {
	done chan struct{}
	once sync.Once
}

func (w *headlessWindow) setTitle(title string) {}
func (w *headlessWindow) close()                { w.once.Do(func() { close(w.done) }) }

func openHeadless(app *application, address string) error {
	picture := os.Getenv("OZDENUI_LOOK")
	browser := findBrowser()
	if browser == "" {
		return errors.New("a picture of the window needs Chrome, Edge or Chromium on this computer (OZDENUI_BROWSER: its path)")
	}
	profile, err := os.MkdirTemp("", "ozdenui-look-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(profile)
	process := exec.Command(browser, "--headless=new", "--remote-debugging-port=0", "--user-data-dir="+profile,
		"--no-first-run", "--no-default-browser-check", "--hide-scrollbars", "--disable-gpu",
		fmt.Sprintf("--window-size=%d,%d", app.window.Width, app.window.Height), "about:blank")
	output, err := process.StderrPipe()
	if err != nil {
		return err
	}
	if err := process.Start(); err != nil {
		return err
	}
	defer func() {
		_ = process.Process.Kill()
		_ = process.Wait()
	}()

	window := &headlessWindow{done: make(chan struct{})}
	app.control = window
	devtools, err := devtoolsAddress(output)
	if err != nil {
		return err
	}
	browserLink, err := dialWebSocket(devtools)
	if err != nil {
		return err
	}
	defer browserLink.close()
	scale := envNumber("OZDENUI_LOOK_SCALE", 1)
	target, err := browserLink.call("", "Target.createTarget", P{"url": "about:blank"})
	if err != nil {
		return err
	}
	attached, err := browserLink.call("", "Target.attachToTarget", P{"targetId": target["targetId"], "flatten": true})
	if err != nil {
		return err
	}
	session, _ := attached["sessionId"].(string)
	metrics := P{"width": app.window.Width, "height": app.window.Height, "deviceScaleFactor": scale, "mobile": false}
	if _, err := browserLink.call(session, "Emulation.setDeviceMetricsOverride", metrics); err != nil {
		return err
	}
	if _, err := browserLink.call(session, "Page.navigate", P{"url": address}); err != nil {
		return err
	}
	if picture == "" {
		<-window.done // a check without a screen: the program closes it
		return nil
	}

	select {
	case <-app.drawn:
	case <-window.done:
		return nil
	case <-time.After(readyTimeout + browserStartWait):
		return errors.New("the page did not start: no picture")
	}
	select {
	case <-time.After(time.Duration(envNumber("OZDENUI_LOOK_WAIT", defaultLookWait) * float64(time.Second))):
	case <-window.done:
	}
	shot, err := browserLink.call(session, "Page.captureScreenshot", P{"format": "png"})
	if err != nil {
		return err
	}
	data, _ := shot["data"].(string)
	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(picture, content, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "ozdenui: the picture of the window: %s\n", picture)
	return nil
}

func envNumber(name string, fallback float64) float64 {
	if number, err := strconv.ParseFloat(os.Getenv(name), 64); err == nil && number > 0 {
		return number
	}
	return fallback
}

// The browser writes "DevTools listening on ws://127.0.0.1:port/devtools/browser/..." when it is ready.
func devtoolsAddress(output io.Reader) (string, error) {
	found := make(chan string, 1)
	go func() {
		pattern := regexp.MustCompile(`ws://\S+`)
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			if address := pattern.FindString(scanner.Text()); address != "" {
				found <- address
				break
			}
		}
		_, _ = io.Copy(io.Discard, output) // WHY: a full pipe would stop the browser
	}()
	select {
	case address := <-found:
		return address, nil
	case <-time.After(browserStartWait):
		return "", errors.New("the hidden browser did not start")
	}
}

// *** A SMALL WEBSOCKET CLIENT, for the DevTools protocol of the browser (text messages of JSON).

type webSocket struct {
	connection net.Conn
	reader     *bufio.Reader
	writeMu    sync.Mutex
	mu         sync.Mutex
	nextID     int
	waiting    map[int]chan P
}

func dialWebSocket(address string) (*webSocket, error) {
	parsed, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	connection, err := net.DialTimeout("tcp", parsed.Host, browserStartWait)
	if err != nil {
		return nil, err
	}
	key := make([]byte, 16)
	_, _ = rand.Read(key)
	request := "GET " + parsed.RequestURI() + " HTTP/1.1\r\nHost: " + parsed.Host + "\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + base64.StdEncoding.EncodeToString(key) + "\r\nSec-WebSocket-Version: 13\r\n\r\n"
	if _, err := connection.Write([]byte(request)); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(connection)
	status, err := reader.ReadString('\n')
	if err != nil || !strings.Contains(status, "101") {
		return nil, fmt.Errorf("the hidden browser refused the link: %s", strings.TrimSpace(status))
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}
	socket := &webSocket{connection: connection, reader: reader, waiting: map[int]chan P{}}
	go socket.read()
	return socket, nil
}

func (s *webSocket) close() { _ = s.connection.Close() }

// Calls a method of the DevTools protocol and returns its result.
func (s *webSocket) call(session, method string, params P) (P, error) {
	waiter := make(chan P, 1)
	s.mu.Lock()
	s.nextID++
	id := s.nextID
	s.waiting[id] = waiter
	s.mu.Unlock()
	message := P{"id": id, "method": method, "params": params}
	if session != "" {
		message["sessionId"] = session
	}
	text, _ := json.Marshal(message)
	if err := s.write(text); err != nil {
		return nil, err
	}
	select {
	case answer := <-waiter:
		if answer == nil {
			return nil, errors.New(method + ": the hidden browser is gone")
		}
		if problem, has := answer["error"].(map[string]any); has {
			return nil, fmt.Errorf("%s: %v", method, problem["message"])
		}
		result, _ := answer["result"].(map[string]any)
		return result, nil
	case <-time.After(browserStartWait):
		return nil, errors.New(method + ": no answer of the hidden browser")
	}
}

// One text frame; a client masks what it sends.
func (s *webSocket) write(payload []byte) error {
	header := []byte{0x81}
	switch {
	case len(payload) < 126:
		header = append(header, 0x80|byte(len(payload)))
	case len(payload) < 65536:
		header = append(header, 0x80|126, byte(len(payload)>>8), byte(len(payload)))
	default:
		header = append(header, 0x80|127)
		header = binary.BigEndian.AppendUint64(header, uint64(len(payload)))
	}
	mask := make([]byte, 4)
	_, _ = rand.Read(mask)
	masked := make([]byte, len(payload))
	for i, value := range payload {
		masked[i] = value ^ mask[i%4]
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_, err := s.connection.Write(append(append(header, mask...), masked...))
	return err
}

func (s *webSocket) read() {
	defer func() {
		s.mu.Lock()
		for _, waiter := range s.waiting {
			waiter <- nil
		}
		s.waiting = map[int]chan P{}
		s.mu.Unlock()
	}()
	var message []byte
	for {
		head := make([]byte, 2)
		if _, err := io.ReadFull(s.reader, head); err != nil {
			return
		}
		isLast := head[0]&0x80 != 0
		opcode := head[0] & 0x0F
		length := uint64(head[1] & 0x7F)
		if length == 126 {
			more := make([]byte, 2)
			if _, err := io.ReadFull(s.reader, more); err != nil {
				return
			}
			length = uint64(binary.BigEndian.Uint16(more))
		} else if length == 127 {
			more := make([]byte, 8)
			if _, err := io.ReadFull(s.reader, more); err != nil {
				return
			}
			length = binary.BigEndian.Uint64(more)
		}
		payload := make([]byte, length)
		if _, err := io.ReadFull(s.reader, payload); err != nil {
			return
		}
		if opcode == 0x8 {
			return // closed
		}
		if opcode > 0x2 {
			continue // a ping or a pong: the browser does not wait for an answer on this link
		}
		message = append(message, payload...)
		if !isLast {
			continue
		}
		var answer P
		if err := json.Unmarshal(message, &answer); err == nil {
			if id, has := answer["id"].(float64); has {
				s.mu.Lock()
				waiter := s.waiting[int(id)]
				delete(s.waiting, int(id))
				s.mu.Unlock()
				if waiter != nil {
					waiter <- answer
				}
			}
		}
		message = nil
	}
}
