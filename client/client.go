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
)

var (
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
		var message string
		dec := gob.NewDecoder(c)
		err := dec.Decode(&message)
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		// Close client connection if the server has shut down
		if err == io.EOF {
			fmt.Println("Server has shut down. Client is now closing... ")
			c.Close()
			os.Exit(0)
		}
		fmt.Println(message)
		fmt.Print(">> ")
	}
}

func write(c net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if strings.Contains(message, "EXIT") {
			fmt.Println("Exiting the client...")
			return
		}
		enc := gob.NewEncoder(c)
		if err := enc.Encode(message); err != nil {
			log.Fatal(err)
		}
	}
}
