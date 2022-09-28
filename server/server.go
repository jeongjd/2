package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	port = ":1123"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	openConnections = make(map[net.Conn]bool)
	newConnection   = make(chan net.Conn)
	// deadConnection  = make(chan net.Conn)
)

func main() {
	fmt.Println("Launching a TCP Chatroom Server...")
	go createTCPServer(port)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	if strings.Contains(text, "EXIT") {
		fmt.Println("Exiting the server...")
		os.Exit(0)
	}
}

func createTCPServer(port string) {
	l, err := net.Listen("tcp", port)
	logFatal(err)
	defer l.Close()

	go func() {
		for {
			c, err := l.Accept()
			logFatal(err)
			openConnections[c] = true
			newConnection <- c
		}
	}()
	for {
		c := <-newConnection
		go broadcastMessage(c)
	}
}

func broadcastMessage(c net.Conn) {
	for {
		reader := bufio.NewReader(c)
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// loop through all the open connections and send messages to these connections
		// except the connection that sent the message
		for item := range openConnections {
			if item != c {
				item.Write([]byte(text))
			}
		}
	}
}
