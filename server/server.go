package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Message struct {
	receiverID     string
	senderID       string
	messageContent string
}

var (
	connections = make(map[string]net.Conn)
	count       = 0
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
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
		count++
	}
}

// partially from https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
func handleConnection(c net.Conn) {
	for {
		text, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		textParsed := parseLine(text)
		temp := strings.TrimSpace(string(text))
		if temp == "STOP" {
			break
		}

		if len(textParsed) == 1 && strings.Contains(text, "/") {
			// fmt.Println("Contains '/' ")
			username := strings.Trim(text, "/")
			username = strings.Trim(username, " \r\n")
			connections[username] = c
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
				fmt.Fprintf(c, "Invalid input! Please type in the form of {To:user} {From:user} {message} "+"\n")
				c.Write([]byte(text))
			}
			broadcastMessage(m)
		}
	}
	c.Close()
}

func parseLine(line string) []string {
	return strings.Split(line, " ")
}

func broadcastMessage(m Message) {
	fmt.Println("this is connection# ", count)

	// Loop through all the open connections and send messages to a specific user
	for item := range connections {
		if item == m.receiverID {
			connections[item].Write([]byte(m.messageContent))
		}
	}
}
