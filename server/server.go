package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
)

type Message struct {
	receiverID     string
	senderID       string
	messageContent string
}

var (
	// Map - key: (client) username, value: connection
	clientConnections = make(map[string]net.Conn)
	// For switch/cases - printing error messages
	option = 0
	// Read/Write mutex to synchronize the clientConnections hashmap between the threads (instead of a channel)
	clientConnectionsMutex = sync.RWMutex{}
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
	// If using channels
	// quit := make(chan string)
	// go closeServer(quit)
	go closeServer()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

// Close server when "EXIT" is typed in the server side
func closeServer() {
	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if strings.TrimSpace(text) == "EXIT" {
			fmt.Println("Server is shutting down... ")
			os.Exit(0)
		}
	}
}

// partially from https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
// Handle client connections - invoke other functions depending on the messages received
func handleConnection(c net.Conn) {
	for {
		var text string
		// Reads and decodes data from connection
		dec := gob.NewDecoder(c)
		err := dec.Decode(&text)
		if err != nil {
			log.Fatal(err)
		}
		// If a connection is closed delete the username from map (clientConnections)
		if err != nil {
			name := getKey(c)
			clientConnectionsMutex.Lock()
			defer clientConnectionsMutex.Unlock()
			delete(clientConnections, name)
			fmt.Printf("User '%s' disconnected from the server\n", name)
			fmt.Println("remaining clients: ", clientConnections)
			return
		}
		m := parseMessage(c, text)
		// Check if the struct Message is null or not
		if reflect.ValueOf(m).IsZero() == false {
			// Check if the message has proper format (senderID, receiverID)
			if checkClients(c, m) == true {
				broadcastMessage(m)
			} else {
				printErrorMessage(c, m)
			}
		} else {
			printErrorMessage(c, m)
		}
	}
}

// Get key of a map based on a value
func getKey(c net.Conn) string {
	clientConnectionsMutex.RLock()
	defer clientConnectionsMutex.RUnlock()
	for key, value := range clientConnections {
		if c == value {
			return key
		}
	}
	return "Key does not Exist"
}

// Parse user messages, store variables, and return struct Message
func parseMessage(c net.Conn, text string) Message {
	textParsed := parseLine(text)
	var m Message
	// Store client username in the map (clientConnections)
	if len(textParsed) == 1 && strings.Contains(text, "/") {
		username := strings.Trim(text, "/")
		username = strings.Trim(username, " \r\n")
		// Prevents other go routines from editing the clientConnections hashmap in order to synchronize the routines
		clientConnectionsMutex.Lock()
		clientConnections[username] = c
		clientConnectionsMutex.Unlock()
	} else if len(textParsed) >= 3 {
		receiver := textParsed[0]
		sender := textParsed[1]
		textTrimmed := strings.Join(textParsed, " ")
		needsTrim := receiver + " " + sender
		textTrimmed = strings.TrimPrefix(textTrimmed, needsTrim)
		msg := sender + ":" + textTrimmed
		m = Message{receiver, sender, msg}
	} else {
		// If message has invalid format
		option = 1
	}
	return m
}

// Split a string into a string array
func parseLine(line string) []string {
	return strings.Split(line, " ")
}

// Check if certain keys exist in a map
func checkKey(str string) bool {
	// Prevents other go routines from editing the clientConnections hashmap in order to synchronize the routines
	clientConnectionsMutex.RLock()
	defer clientConnectionsMutex.RUnlock()
	for item := range clientConnections {
		if item == str {
			return true
		}
	}
	return false
}

// Check if certain usernames exist in clientConnections map
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
		// If sender and/or receiver usernames do not exist
		option = 3
		return false
	}
}

// Send private message to a specific client using gob
func broadcastMessage(m Message) {
	// Prevents other go routines from editing the clientConnections hashmap in order to synchronize the routines
	clientConnectionsMutex.RLock()
	defer clientConnectionsMutex.RUnlock()
	for item := range clientConnections {
		if item == m.receiverID {
			enc := gob.NewEncoder(clientConnections[item])
			enc.Encode(m.messageContent)
		}
	}
}

// Print error messages depending on option (which error)
func printErrorMessage(c net.Conn, m Message) {
	enc := gob.NewEncoder(c)
	var errorMessage string
	switch option {
	case 1:
		errorMessage = "Invalid input! Please type in the form of {To:user} {From:user} {message} \n"
	case 2:
		errorMessage = "You are not " + m.senderID + "!"
	case 3:
		errorMessage = "Invalid user! That client has not connected"
	}
	// Encodes and sends error message to client
	if err := enc.Encode(errorMessage); err != nil {
		log.Fatal(err)
	}
}
