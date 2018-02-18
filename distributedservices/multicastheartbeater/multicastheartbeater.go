package multicastheartbeater

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"net"
	"os"
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

/*
Specification:
Returns a channel of utilities.HeartBeat
The caller can read heartbeats on this channel at the speed that UDP
provides; with a time lag associated with go channels
*/
func SendHeartBeatMessages(toAddress, toPort, fromPort string) chan utilities.HeartBeat {
	heartbeatChannelIn := make(chan utilities.HeartBeat)
	go func() {
		var LeaderAddress string
		LeaderAddress = toAddress
		LeaderAddress += ":"
		LeaderAddress += toPort
		SenderPort := ":" + fromPort
		ServerAddr, err := net.ResolveUDPAddr("udp", LeaderAddress)
		CheckError(err)
		LocalAddr, err := net.ResolveUDPAddr("udp", SenderPort)
		CheckError(err)
		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		CheckError(err)
		defer Conn.Close()
		for {
			hb := <-heartbeatChannelIn
			//encode json data
			//fmt.Printf("Data to be Sent:%v\n", hb)
			jsonData, err := json.Marshal(hb)
			//fmt.Printf("Marshalled Data:%v\n", string(jsonData))
			_, err = Conn.Write(jsonData)
			if err != nil {
				fmt.Println(err.Error())
			}
			//fmt.Printf("Wrote %d bytes\n", n)
			//time.Sleep(time.Second * 1)
		}
	}()
	return heartbeatChannelIn
}
