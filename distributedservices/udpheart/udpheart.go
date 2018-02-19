package udpheartbeat

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type HeartBeat struct {
	Cluster   []string
	ReqNumber int64
	ReqCode   int8
}

type HeartBeatUpperStack struct {
	Ip string
	Hb HeartBeat
}

func Init(readPort string, writePort string) (<-chan HeartBeatUpperStack, chan<- HeartBeatUpperStack) {
	receive := make(chan HeartBeatUpperStack, 10)
	send := make(chan HeartBeatUpperStack, 10)
	go listen(receive, readPort)
	go relay(send, writePort)
	return receive, send
}

func listen(rcvChannel chan HeartBeatUpperStack, readPort string) {
	// on rcvChannel the application will receive after UDP Read
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	go func() {
		for {
			port, errconv := strconv.Atoi(readPort)
			if errconv != nil {
				fmt.Println("Port is not a numeric string")
			}

			conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
				IP:   net.ParseIP("124.0.0.1"),
				Port: port,
			})
			if err != nil {
				fmt.Println("Fault 4 in Multicast: ", err)
			}

			defer conn.Close()

			buf := make([]byte, 100)
			n, udpAddr, err2 := conn.ReadFromUDP(buf)
			if err2 != nil {
				fmt.Printf("Error Reading From UDP:%v\n", err2.Error())
			}
			buf = buf[:n]
			var Result HeartBeat

			errUnmarshal := json.Unmarshal(buf, &Result)
			if errUnmarshal != nil {
				fmt.Printf("Error Unmarshalling:%v\n", errUnmarshal.Error())
			}
			//Decode the data
			//Read JSON from the peer udp datagrams

			//send the information up the stack for processing
			//fmt.Printf("Received data:%v\n", Result)
			//fmt.Printf("Request Number:%v", Result.ReqNumber)

			var hbu HeartBeatUpperStack
			hbu.Hb = Result
			hbu.Ip = udpAddr.IP.String()
			fmt.Printf("Received Lower stack:%v\n", hbu)
			rcvChannel <- hbu
		}
	}()
}

/*
heartbeatChannelIn := make(chan utilities.HeartBeat)
	go func() {
		var LeaderAddress string
		LeaderAddress = toAddress

		LeaderAddress += ":"
		LeaderAddress += toPort
		ServerAddr, err := net.ResolveUDPAddr("udp", LeaderAddress)
		CheckError(err)

		Conn, err := net.DialUDP("udp", nil, ServerAddr)
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
*/
func relay(sendChannel chan HeartBeatUpperStack, writePort string) {
	// What we receive on sendChannel we relay to our network
	broadcast_addr := "224.0.0.1"
	destinationAddress, _ := net.ResolveUDPAddr("udp", broadcast_addr+writePort)
	connection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer connection.Close()
	for {

	}
}
