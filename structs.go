package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
	Name string
	out  chan string
}

func (m *Client) ChangeName(name string) {
	m.Name = name
}

type Message struct {
	Name    string
	Time    time.Time
	Message string
}

func (m Message) NewMessage() string {
	mes := strings.TrimRight(m.Message, "\r\n")
	dat := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%v] %s: %v\r\n", dat, m.Name, mes)
}

type Server struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan MessageEvent
	clients    map[*Client]struct{}
}

type MessageEvent struct {
	from *Client
	text string
}

func NewServer() *Server {
	return &Server{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan MessageEvent, 64),
		clients:    make(map[*Client]struct{}),
	}
}

func (m *Server) Register(client *Client) {
	m.register <- client
}

func (m *Server) Unregister(client *Client) {
	m.unregister <- client
}
func (m *Server) Broadcast(message MessageEvent) {
	m.broadcast <- message
}
func (m *Server) Run() {
	for {
		select {
		case c := <-m.register:
			m.clients[c] = struct{}{}
		case c := <-m.unregister:
			if _, ok := m.clients[c]; ok {
				delete(m.clients, c)
				// Stop the per-client writer goroutine.
				close(c.out)
			}
		case message := <-m.broadcast:
			for s := range m.clients {
				// Do not send the message back to the author.
				if message.from != nil && message.from.conn == s.conn {
					continue
				}
				// Best-effort: don't let a slow client block the whole server.
				select {
				case s.out <- mesToText(message):
				default:
				}
			}
		}
	}
}
