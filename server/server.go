package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

type Connection struct {
	connection net.Conn
	sender     string
}

var (
	connections = make(map[string]net.Conn)
	receiver    = " "
	sender      = " "
	msg         = " "
	count       = 0
	username    = " "
)

//	func main() {
//		fmt.Print("Enter a port number: ")
//		fmt.Scanln(&port)
//		port = ":" + port
//		fmt.Println("Launching a TCP Chatroom Server...")
//		go createTCPServer(port)
//		reader := bufio.NewReader(os.Stdin)
//		fmt.Print(">> ")
//		text, _ := reader.ReadString('\n')
//		if strings.Contains(text, "EXIT") {
//			fmt.Println("Exiting the server...")
//			os.Exit(0)
//		}
//	}
//
// partially from https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
func main() {
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
			username = strings.Trim(text, "/")
			username = strings.Trim(username, " \r\n")
			connections[username] = c

		} else {
			if len(textParsed) >= 3 {
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
			//check(username)
			broadcastMessage(c, m)
		}

	}
	c.Close()

}

func parseLine(line string) []string {
	return strings.Split(line, " ")
}

func broadcastMessage(c net.Conn, m Message) {
	fmt.Println("this is connection# ", count)
	// check which client sent the message
	// check who the client is sending the message to
	// send message to that client

	// loop through all the open connections and send messages to these connections
	// except the connection that sent the message
	for item := range connections {
		if item == m.receiverID {
			connections[item].Write([]byte(m.messageContent))
		}
		//if item.sender == m.receiverID {
		//	fmt.Println("ReceiverID ", m.receiverID)
		//	fmt.Println("Entered the boolean ")
		//	item.connection.Write([]byte(m.messageContent))
		//}
	}
}
