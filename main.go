package main

import (
	"fmt"
	"os"
)

func main() {

	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) != 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	startserver(port)
}
