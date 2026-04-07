package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

//func writeother(server *Server, text string, exclude net.Conn) {
//	for i := 0; i < len(server); i++ {
//		c := server.connects[i]
//		if exclude != nil && c == exclude {
//			continue
//		}
//		//writerc := bufio.NewWriter(c)
//		//
//		//if _, err := writerc.WriteString(text); err != nil {
//		//	continue
//		//}
//		//if err := writerc.Flush(); err != nil {
//		//	continue
//		//}
//		writeonlytoclient(text, c)
//	}
//}

func writeonlytoclient(text string, conn net.Conn) {
	writerC := bufio.NewWriter(conn)
	if _, err := writerC.WriteString(text); err != nil {
		return
	}
	if err := writerC.Flush(); err != nil {
		return
	}
}

func textRefactor(text string) string {
	return strings.TrimSpace(strings.TrimRight(text, "\r\n"))
}

func checkCommand(client *Client, msg string) (handled bool, quit bool) {
	fields := strings.Fields(msg)
	if len(fields) == 0 {
		return false, false
	}

	switch fields[0] {
	case "/help":
		writeonlytoclient("Commands:\r\n/help\r\n/name name\r\n/quit\r\n", client.conn)
		return true, false
	case "/name":
		if len(fields) < 2 {
			writeonlytoclient("Use: /name [name]\r\n", client.conn)
			return true, false
		}
		client.ChangeName(fields[1])
		writeonlytoclient("OK\r\n", client.conn)
		fmt.Printf("Client changed his name to %s\r\n", client.Name)
		return true, false
	case "/quit":
		writeonlytoclient("Bye bye\r\n", client.conn)
		fmt.Printf("Client %s left the server\r\n", client.Name)
		_ = client.conn.Close()
		return true, true
	default:
		return false, false
	}
}
