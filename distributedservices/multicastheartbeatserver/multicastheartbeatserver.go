package multicastheartbeatserver

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

/*
Specification

Find out if I am not the leader node.
If I am not the leader node I take a different course of action
(not implemented yet)

If I am the leader then I listen to pings and ask my local
http server to update it's internal representation of cluster.
*/
func CatchMultiCastDatagramsAndBounce(iListenOnIp, iListenOnPort string, c chan utilities.HeartBeatUpperStack) {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	for {
		//addrstr := string("224.0.0.1")
		port, errconv := strconv.Atoi(iListenOnPort)
		if errconv != nil {
			fmt.Println("Port is not a numeric string")
		}

		conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
			IP:   net.ParseIP(iListenOnIp),
			Port: port,
		})
		if err != nil {
			fmt.Println("Fault 4: ", err)
		}

		defer conn.Close()

		buf := make([]byte, 100)
		n, udpAddr, err2 := conn.ReadFromUDP(buf)
		if err2 != nil {
			fmt.Printf("Error Reading From UDP:%v\n", err2.Error())
		}
		buf = buf[:n]
		var Result utilities.HeartBeat

		errUnmarshal := json.Unmarshal(buf, &Result)
		if errUnmarshal != nil {
			fmt.Printf("Error Unmarshalling:%v\n", errUnmarshal.Error())
		}
		//Decode the data
		//Read JSON from the peer udp datagrams

		//send the information up the stack for processing
		//fmt.Printf("Received data:%v\n", Result)
		//fmt.Printf("Request Number:%v", Result.ReqNumber)

		var hbu utilities.HeartBeatUpperStack
		hbu.Hb = Result
		hbu.Ip = udpAddr.IP.String()
		c <- hbu
	}
}

func CatchUniCastDatagramsAndBounce(iListenOnPort string, c chan utilities.HeartBeatUpperStack) {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	for {
		//addrstr := string("224.0.0.1")
		port, errconv := strconv.Atoi(iListenOnPort)
		if errconv != nil {
			fmt.Println("Port is not a numeric string")
		}
		conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
		if err != nil {
			fmt.Println("Fault 4: ", err)
		}

		defer conn.Close()

		buf := make([]byte, 100)
		n, udpAddr, err2 := conn.ReadFromUDP(buf)
		if err2 != nil {
			fmt.Printf("Error Reading From UDP:%v\n", err2.Error())
		}
		buf = buf[:n]
		var Result utilities.HeartBeat

		errUnmarshal := json.Unmarshal(buf, &Result)
		if errUnmarshal != nil {
			fmt.Printf("Error Unmarshalling:%v\n", errUnmarshal.Error())
		}
		//Decode the data
		//Read JSON from the peer udp datagrams

		//send the information up the stack for processing
		//fmt.Printf("Received data:%v\n", Result)
		//fmt.Printf("Request Number:%v", Result.ReqNumber)

		var hbu utilities.HeartBeatUpperStack
		hbu.Hb = Result
		hbu.Ip = udpAddr.IP.String()
		c <- hbu
	}
}
