package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
)

// Message Struct
type Message struct {
	receiverID     string
	senderID       string
	messageContent string
}

// Declare and initialize global variables
var (
	clientConnections = make(map[string]net.Conn)
	option            = 0
)

// partially from https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
func main() {
	var port string
	fmt.Print("Enter a port number: ")
	fmt.Scanln(&port)
	port = ":" + port
	fmt.Println("Launching a TCP Chatroom Server...")

	l, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

// partially from https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
func handleConnection(c net.Conn) {
	for {
		var text string
		// COMMENT
		dec := gob.NewDecoder(c)
		err := dec.Decode(&text)
		// COMMENT
		if err != nil {
			name := getKey(c)
			// COMMENT
			delete(clientConnections, name)
			fmt.Printf("User '%s' left the server\n", name)
			fmt.Println("remaining clients: ", clientConnections)
			return
		}
		m := parseMessage(c, text)
		// COMMENT
		if reflect.ValueOf(m).IsZero() == false {
			// COMMENT
			if checkClients(c, m) == true {
				broadcastMessage(m)
			} else {
				// COMMENT
				printErrorMessage(c, m)
			}
		} else {
			printErrorMessage(c, m)
		}
	}
}

// COMMENT
func getKey(c net.Conn) string {
	for key, value := range clientConnections {
		if c == value {
			return key
		}
	}
	return "Key does not Exist"
}

// COMMENT
func parseMessage(c net.Conn, text string) Message {
	textParsed := parseLine(text)
	var m Message
	if len(textParsed) == 1 && strings.Contains(text, "/") {
		username := strings.Trim(text, "/")
		username = strings.Trim(username, " \r\n")
		clientConnections[username] = c
	} else if len(textParsed) >= 3 {
		receiver := textParsed[0]
		sender := textParsed[1]
		textTrimmed := strings.Join(textParsed, " ")
		needsTrim := receiver + " " + sender
		textTrimmed = strings.TrimPrefix(textTrimmed, needsTrim)
		msg := sender + ":" + textTrimmed
		m = Message{receiver, sender, msg}
	} else {
		// If message is not in the right format
		option = 1
	}
	return m
}

// COMMENT
func parseLine(line string) []string {
	return strings.Split(line, " ")
}

// COMMENT
func checkKey(str string) bool {
	for item := range clientConnections {
		if item == str {
			return true
		}
	}
	return false
}

// COMMENT
func checkClients(c net.Conn, m Message) bool {
	// Check if both sender and receiver usernames exist
	if checkKey(m.senderID) == true && checkKey(m.receiverID) {
		// Check if senderID matches client username
		if getKey(c) == m.senderID {
			return true
		} else {
			// If senderID does not match client username
			option = 2
			return false
		}
	} else {
		// If either or both sender and receiver usernames does not exist
		option = 3
		return false
	}
}

// COMMENT
func broadcastMessage(m Message) {
	// Loop through all the connections and send messages to a specific user
	// COMMENT
	for item := range clientConnections {
		if item == m.receiverID {
			enc := gob.NewEncoder(clientConnections[item])
			enc.Encode(m.messageContent)
		}
	}
}

// COMMENT
func printErrorMessage(c net.Conn, m Message) {
	enc := gob.NewEncoder(c)
	var errorMessage string
	// COMMENT
	switch option {
	case 1:
		errorMessage = "Invalid input! Please type in the form of {To:user} {From:user} {message} \n"
	case 2:
		errorMessage = "You are not " + m.senderID + "!"
	case 3:
		errorMessage = "Invalid user!"
	}
	// COMMENT
	if err := enc.Encode(errorMessage); err != nil {
		log.Fatal(err)
	}
}
