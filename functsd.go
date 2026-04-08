package main

import (
	"bufio"
	"strings"
	"time"
)

func textRefactor(text string) string {
	return strings.TrimSpace(strings.TrimRight(text, "\r\n"))
}

func normalizeUserText(s string) string {
	// Apply backspaces (BS/DEL) and strip terminal control / escape sequences so we
	// don't broadcast cursor movement to other clients.
	out := make([]rune, 0, len(s))
	inEsc := false

	for _, r := range s {
		if inEsc {
			// Rough ANSI/VT escape stripping: end on final byte (@..~).
			if r >= '@' && r <= '~' {
				inEsc = false
			}
			continue
		}

		switch r {
		case '\x1b': // ESC
			inEsc = true
			continue
		case '\b', 0x7f: // BS or DEL
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
			continue
		}

		// Drop other control chars.
		if r < 0x20 {
			continue
		}
		out = append(out, r)
	}
	return strings.TrimSpace(string(out))
}

func clientSend(client *Client, msg string) {
	select {
	case client.out <- msg:
	default:
	}
}

func checkCommand(client *Client, msg string) (handled bool, quit bool) {
	fields := strings.Fields(msg)
	if len(fields) == 0 {
		return false, false
	}

	switch fields[0] {
	case "/help":
		clientSend(client, "Commands:\r\n/help\r\n/name name\r\n/quit\r\n")
		return true, false
	case "/name":
		if len(fields) < 2 {
			clientSend(client, "Use: /name [name]\r\n")
			return true, false
		}
		client.ChangeName(fields[1])
		clientSend(client, "OK\r\n")
		return true, false
	case "/quit":
		clientSend(client, "Bye bye\r\n")
		_ = client.conn.Close()
		return true, true
	default:
		return false, false
	}
}

func clientWriter(client *Client) {
	w := bufio.NewWriter(client.conn)
	for msg := range client.out {
		if _, err := w.WriteString(msg); err != nil {
			return
		}
		if err := w.Flush(); err != nil {
			return
		}
	}
}

func mesToText(mes MessageEvent) string {
	// Leading CRLF helps keep messages from gluing into the receiver's current input line.
	return "\r\n[" + time.Now().Format("2006-01-02 15:04:05") + "] " + mes.from.Name + ": " + mes.text + "\r\n"
}

func broadcastText(text string, client *Client) MessageEvent {
	return MessageEvent{
		from: client,
		text: text,
	}
}
