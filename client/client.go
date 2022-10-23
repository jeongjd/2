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
	go read(c)
	write(c)
}

func read(c net.Conn) {
	for {
		reader := bufio.NewReader(c)
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Connection closed.")
			c.Close()
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
			return
		}
		// message = fmt.Sprintf("%s: %s\n", username, strings.Trim(message, "\n"))
		c.Write([]byte(message))
		fmt.Print(">> ")
	}
}
