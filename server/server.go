package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strings"
)

type Message struct {
	receiverID     string
	senderID       string
	messageContent string
}

var (
	// limit hashmap to 5
	openConnections    = make(map[net.Conn]bool)
	newClient          = make(chan net.Conn)
	disconnectedClient = make(chan net.Conn)
	clientConnections  = make(map[string]net.Conn)
	count              = 0
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
	// look into why this works
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			openConnections[c] = true
			newClient <- c
			count++
		}
	}()
	for {
		select {
		case c := <-newClient:
			// Invoke broadcast message (broadcasts to the other connections
			go handleConnection(c)
		case c := <-disconnectedClient:
			// remove/delete the connection
			for item := range openConnections {
				if item == c {
					delete(openConnections, c)
					fmt.Println("removed connection, remaining: ", openConnections)
					fmt.Println("removed client = ", disconnectedClient)
				}
			}
		}
	}
}

// partially from https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
func handleConnection(c net.Conn) {
	for {
		//text, err := bufio.NewReader(c).ReadString('\n')
		var text string
		dec := gob.NewDecoder(c)
		err := dec.Decode(&text)
		if err != nil {
			disconnectedClient <- c
			name := getKey(clientConnections, c)
			delete(clientConnections, name)
			fmt.Printf("User '%s' left the server\n", name)
			fmt.Println("remaining clients: ", clientConnections)
			fmt.Println(err) // prints "EOF" in server
			return
		}
		// m := parseMessage(c, text)
		textParsed := parseLine(text)
		if len(textParsed) == 1 && strings.Contains(text, "/") {
			username := strings.Trim(text, "/")
			username = strings.Trim(username, " \r\n")
			clientConnections[username] = c
		} else {
			var m Message
			if len(textParsed) >= 3 {
				receiver := textParsed[0]
				sender := textParsed[1]
				textTrimmed := strings.Join(textParsed, " ")
				needsTrim := receiver + " " + sender
				textTrimmed = strings.TrimPrefix(textTrimmed, needsTrim)
				msg := sender + ":" + textTrimmed
				m = Message{receiver, sender, msg}
			} else {
				enc := gob.NewEncoder(c)
				newMessage := "Invalid input! Please type in the form of {To:user} {From:user} {message} \n"
				if err := enc.Encode(newMessage); err != nil {
					log.Fatal(err)
				}
				//fmt.Fprintf(c, "Invalid input! Please type in the form of {To:user} {From:user} {message} \n")
			}
			checkClients(c, m)
		}
	}
}

func parseLine(line string) []string {
	return strings.Split(line, " ")
}

func getKey(clientMap map[string]net.Conn, c net.Conn) string {
	for key, value := range clientConnections {
		if c == value {
			return key
		}
	}
	return "Key does not Exist"
}

func checkKey(str string, clientMap map[string]net.Conn) bool {
	for item := range clientConnections {
		if item == str {
			return true
		}
	}
	return false
}

func parseMessage(c net.Conn, text string) {

}

func checkClients(c net.Conn, m Message) {
	// check if both sender and receiver usernames exist
	if checkKey(m.senderID, clientConnections) == true && checkKey(m.receiverID, clientConnections) {
		// Check if senderID matches client username
		if getKey(clientConnections, c) == m.senderID {
			broadcastMessage(c, m)
		} else {
			fmt.Fprintf(c, "You are not %s! \n", m.senderID)
		}
	} else {
		enc := gob.NewEncoder(c)
		errorMessage := "Invalid user! \n"
		if err := enc.Encode(errorMessage); err != nil {
			log.Fatal(err)
		}
		//fmt.Fprintf(c, "Invalid user! \n")
	}
}

func broadcastMessage(c net.Conn, m Message) {
	fmt.Println("this is connection# ", count)

	// Loop through all the connections and send messages to a specific user
	for item := range clientConnections {
		if item == m.receiverID {
			enc := gob.NewEncoder(clientConnections[item])
			enc.Encode(m.messageContent)
			//clientConnections[item].Write([]byte(m.messageContent))
		}
	}
}
