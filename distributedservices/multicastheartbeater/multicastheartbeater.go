package main

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

/*It's a multicast ip address on which leader listens
to ADD requests.*/
const Leaderaddress = "224.0.0.1:10001"

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", Leaderaddress)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", ":10002")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	i := 0
	for {
		msg := strconv.Itoa(i)
		_ = msg
		i++

		hb := utilities.HeartBeat{
			ReqNumber: 12345,
			ReqCode:   5,
		}
		//encode json data
		fmt.Printf("Data to be Sent:%v\n", hb)
		jsonData, err := json.Marshal(hb)
		fmt.Printf("Marshalled Data:%v\n", string(jsonData))
		n, err := Conn.Write(jsonData)
		if err != nil {
			fmt.Println(msg, err)
		}
		fmt.Printf("Wrote %d bytes\n", n)
		time.Sleep(time.Second * 1)
	}
}