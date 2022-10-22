package main

import (
	"bufio"
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
	myTime      = t.Format(time.RFC3339) + "\n"
	hostAddress = " "
	port        = " "
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Print("Enter a host address: ")
	fmt.Scanln(&hostAddress)
	fmt.Print("Enter a port number: ")
	fmt.Scanln(&port)
	createTCPClient()
}

func createTCPClient() {
	c, err := net.Dial("tcp", hostAddress+":"+port)
	logFatal(err)

	fmt.Print("enter your username : ")
	/*
		reader := bufio.NewReader(os.Stdin)
		username, err := reader.ReadString('\n')
		logFatal(err)
		username = strings.Trim(username, " \r\n")
		fmt.Printf("Welcome user %s! Send messages to other users.\n", username)
		fmt.Print(">> ")

	*/

	// read
	go read(c)

	// write
	write(c)
}

func read(c net.Conn) {
	for {
		reader := bufio.NewReader(c)
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			c.Close()
			fmt.Println("Connection closed.")
			os.Exit(0)
		}
		message = strings.TrimSpace(message)
		printMessage := fmt.Sprintf("[%s] %s\n", t.Format(time.Kitchen), message)
		fmt.Println(printMessage)
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
			os.Exit(0)
		}
		// message = fmt.Sprintf("%s: %s\n", username, strings.Trim(message, "\n"))
		c.Write([]byte(message))
		fmt.Print(">> ")
	}
}
