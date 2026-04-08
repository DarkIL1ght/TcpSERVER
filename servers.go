package main

import (
	"bufio"
	"fmt"
	"net"
)

func startserver(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()
	server := NewServer()
	go server.Run()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accepted new connection", conn.RemoteAddr())
		go handle(conn, server)
	}
}

func handle(conn net.Conn, server *Server) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	fmt.Println("Client connected:", conn.RemoteAddr())

	client := Client{
		Name: "anonymous",
		conn: conn,
		out:  make(chan string, 128),
	}

	server.Register(&client)

	defer server.Unregister(&client)

	go clientWriter(&client)

	clientSend(&client, "Welcome! What is your name?\r\n")

	nameLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name := textRefactor(nameLine)
	if name == "" {
		name = "anonymous"
	}
	client.Name = normalizeUserText(name)

	clientSend(&client, "\rWelcome to my server: "+client.Name+"\r\n")
	server.Broadcast(MessageEvent{from: &Client{Name: "Server"}, text: "New user " + client.Name + " joined"})

	fmt.Printf("New user {%s} joined\r\n", client.Name)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			server.Broadcast(MessageEvent{from: &Client{Name: "Server"}, text: "Client disconnected: " + client.Name})
			return
		}

		msg := textRefactor(message)
		msg = normalizeUserText(msg)
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

		server.Broadcast(broadcastText(msg, &client))

	}
}
