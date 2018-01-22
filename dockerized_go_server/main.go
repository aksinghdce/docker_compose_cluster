package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	connPort = "3000"
	connType = "tcp"
)

type command struct {
	Name    string
	Options string
	Pattern string
}

func main() {
	// Liten for incoming connections. Don't need to specify an ipaddress here as per golang
	// documentation
	l, err := net.Listen(connType, ":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + ":" + connPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 2000)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Number of bytes read: %d", reqLen)
	var cd command
	errr := json.Unmarshal(buf, &cd)
	if errr != nil {
		fmt.Println("Error:", errr)
	}
	fmt.Println("Data:", cd)
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
