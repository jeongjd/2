package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	port = " "
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Message struct {
	receiverID     string
	senderID       string
	messageContent string
}

var things map[string](chan int)
var (
	openConnections   = make(map[net.Conn]bool)
	newConnection     = make(chan net.Conn)
	clientConnections = make(map[string]bool)
	clientIDs         = make(chan string)
	now               = time.Now()
	clientID          = make(map[string]chan string)

	// deadConnection  = make(chan net.Conn)
	receiver = " "
	sender   = " "
	msg      = " "
	count    = 0
	username = " "
	// m        = Message{receiver, sender, msg}
)

func main() {
	fmt.Print("Enter a port number: ")
	fmt.Scanln(&port)
	port = ":" + port
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
			count++
			fmt.Println("Number of connected clients = ", count)
			fmt.Print(">> ")
		}
	}()
	for {
		c := <-newConnection
		go receive(c)
		/*
			username := <-clientIDs
			fmt.Println(username)
			check(c, username)
		*/
	}
}

// change this function to receive(c net.Conn)
// put the broadcast part into another function broadcastMessage with destination
func receive(c net.Conn) {
	for {
		reader := bufio.NewReader(c)
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		textParsed := parseLine(text)
		if len(textParsed) == 1 && strings.Contains(text, "/") {
			// fmt.Println("Contains '/' ")
			username = strings.Trim(text, "/")
			// fmt.Print("username = ", username)

			// channel username into clientID map
			go func() {
				clientConnections[username] = true
				clientIDs <- username
			}()
		} else if len(textParsed) >= 3 {
			receiver = textParsed[0]
			sender = textParsed[1]
			textTrimmed := strings.Join(textParsed, " ")
			needsTrim := receiver + " " + sender
			textTrimmed = strings.TrimPrefix(textTrimmed, needsTrim)
			msg = textTrimmed
		} else {
			fmt.Fprintf(c, "Invalid input! Please type in the form of {To:user} {From:user} {message} "+"\n")
			c.Write([]byte(text))
		}
		m := Message{receiver, sender, msg}

		// name := <-clientIDs
		// fmt.Println(name)
		// check(c, name)
		check(username)
		broadcastMessage(c, m)

	}
}

func broadcastMessage(c net.Conn, m Message) {
	// check which client sent the message
	// check who the client is sending the message to
	// send message to that client

	// loop through all the open connections and send messages to these connections
	// except the connection that sent the message
	for item := range openConnections {
		// fmt.Println(item)
		if item != c {
			item.Write([]byte(m.messageContent))
		}
	}
}

func check(username string) {
	for j := range clientConnections {
		fmt.Println(j)
		/*
			if j != username {
				fmt.Println(j)
				// item.Write([]byte(invalidUser))
			}
			else {
				invalidUser := fmt.Sprintf("user does not exist!")
				fmt.Println(invalidUser)
				// item.Write([]byte(invalidUser))
			}

		*/
	}
	// broadcastMessage(c, m)
}

func parseLine(line string) []string {
	return strings.Split(line, " ")
}
