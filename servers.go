package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func startserver(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()
	server := NewServer()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accepted new connection", conn.RemoteAddr())
		server.connects = append(server.connects, conn)
		go handle(conn, &server)
	}
}

func handle(conn net.Conn, server *Server) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	fmt.Println("Client connected:", conn.RemoteAddr())

	writeonlytoclient("Welcome! What is your name?\r\n", conn)

	nameLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name := textRefactor(nameLine)
	if name == "" {
		name = "anonymous"
	}

	client := Client{
		Name: name,
		conn: conn,
	}
	writeonlytoclient("\rWelcome to my server: "+client.Name+"\r\n", conn)

	writeother(server, "\rNew user "+client.Name+" joined\r\n", conn)

	fmt.Printf("New user {%s} joined\r\n", client.Name)

	for {
		writeonlytoclient("["+time.Now().Format("2006-01-02 15:04:05")+"] "+client.Name+": ", conn)

		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			writeother(server, "Client disconnected: "+client.Name+"\r\n", nil)
			return
		}

		msg := textRefactor(message)
		if msg == "" {
			continue
		}

		handled, quit := checkCommand(&client, msg)
		if handled {
			if quit {
				return
			}
			continue
		}

		mes := Message{Name: client.Name, Time: time.Now(), Message: msg}
		text := mes.NewMessage()
		writeother(server, text, conn)
		fmt.Printf(text)

		HistoryMassive = append(HistoryMassive, text)

	}
}
