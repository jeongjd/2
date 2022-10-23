package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	t           = time.Now()
	hostAddress = " "
	port        = " "
)

func main() {
	fmt.Print("Enter a host address: ")
	fmt.Scanln(&hostAddress)
	fmt.Print("Enter a port number: ")
	fmt.Scanln(&port)
	createTCPClient()
}

func createTCPClient() {
	c, err := net.Dial("tcp", hostAddress+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Enter your username in the format /name: ")
	go read(c)
	write(c)
}

func read(c net.Conn) {
	for {
		//reader := bufio.NewReader(c)
		//message, err1 := reader.ReadString('\n')
		var message string
		dec := gob.NewDecoder(c)
		err := dec.Decode(&message)
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		if err == io.EOF {
			fmt.Println("Connection closed.")
			c.Close()
			os.Exit(0)
		}
		//newMessage = strings.TrimSpace(message)
		//printMessage := fmt.Sprintf("[%s] %s\n", t.Format(time.Kitchen), message)
		fmt.Println(message)
		fmt.Print(">> ")
	}
}

func write(c net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(message, "EXIT") {
			fmt.Println("Exiting the client...")
			return
		}
		// message = fmt.Sprintf("%s: %s\n", username, strings.Trim(message, "\n"))
		enc := gob.NewEncoder(c)
		if err := enc.Encode(message); err != nil {
			log.Fatal(err)
		}
		//c.Write([]byte(message))
		fmt.Print(">> ")
	}
}
