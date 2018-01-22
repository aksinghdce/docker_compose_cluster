package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type command struct {
	Name    string
	Options string
	Pattern string
}

func main() {
	grepCommand := command{"grep", "-nr", "amit"}
	servAddr := "grepservice:3000"
	var marshalled []byte
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}

	marshalled, err = json.Marshal(grepCommand)
	if err != nil {
		fmt.Println("Marshal call was UNsuccessful")
	}

	//testing on the client side wheather the unmarshalling works
	var cd command
	json.Unmarshal(marshalled, &cd)
	fmt.Println("CD:", cd.Name, cd.Options, cd.Pattern, cd)

	_, err = conn.Write(marshalled)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("reply from server=", string(reply))

	conn.Close()
}
