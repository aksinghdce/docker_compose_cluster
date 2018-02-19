package multicastheartbeatserver

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

/*
Specification

Find out if I am not the leader node.
If I am not the leader node I take a different course of action
(not implemented yet)

If I am the leader then I listen to pings and ask my local
http server to update it's internal representation of cluster.
*/
func CatchMultiCastDatagramsAndBounce(iListenOnPort string) chan utilities.HeartBeatUpperStack {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	c := make(chan utilities.HeartBeatUpperStack)
	go func() {
		for {
			port, errconv := strconv.Atoi(iListenOnPort)
			if errconv != nil {
				fmt.Println("Port is not a numeric string")
			}

			/*Get all the interfaces on the machine*/
			ifsArr, err := net.Interfaces()
			if err != nil {
				fmt.Print(err.Error())
			}

			/*For the interfaces that support multicast, listen to
			a udp port and wait for x seconds to see if there is a
			node that wants to join as a node in the cluster I am managing.
			*/
			for _, ifs := range ifsArr {
				flag := ifs.Flags.String()
				if strings.Contains(flag, "multicast") {
					fmt.Println("Interface:", ifs.Name)
					fmt.Println("Interface Flag:", ifs.Flags.String())

					multicastaddresses, err := ifs.MulticastAddrs()
					if err != nil {
						fmt.Println("Fault 1: ", err)
						continue
					}
					/*I am the leader I will keep listening to events
					from my group.*/
					for {
						for _, addr := range multicastaddresses {
							//fmt.Printf("index:%d, addr:%v", index, addr.String())
							conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
								IP:   net.ParseIP(addr.String()),
								Port: port,
							})
							if err != nil {
								continue
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
							fmt.Printf("Received from:%v\n", udpAddr)
							c <- hbu
						}
					}
				}
			}

		}
	}()
	return c
}

func CatchUniCastDatagramsAndBounce(iListenOnPort string) chan utilities.HeartBeatUpperStack {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	c := make(chan utilities.HeartBeatUpperStack)
	go func() {
		for {

			ifsArr, err := net.Interfaces()
			if err != nil {
				fmt.Print(err.Error())
			}

			for _, ifs := range ifsArr {
				//FlagPointToPoint
				flag := ifs.Flags.String()
				if strings.Contains(flag, "pointtopoint") {
					fmt.Printf("Interface:%v\n", ifs.Name)
					fmt.Printf("Interface Flag:%v\n", ifs.Flags.String())
					unicastaddresses, err := ifs.Addrs()
					if err != nil {
						fmt.Println("Fault 1: ", err)
						continue
					}

					for _, uniaddr := range unicastaddresses {
						//addrstr := string("224.0.0.1")
						port, errconv := strconv.Atoi(iListenOnPort)
						if errconv != nil {
							fmt.Println("Port is not a numeric string")
						}
						myaddr := &net.UDPAddr{
							IP:   net.ParseIP(uniaddr.String()),
							Port: port,
						}
						conn, err := net.ListenUDP("udp", myaddr)
						if err != nil {
							fmt.Println("Fault 4 in Unicast: ", err)
							continue
						}
						defer conn.Close()

						buf := make([]byte, 100)
						n, udpAddr, err2 := conn.ReadFromUDP(buf)
						if err2 != nil {
							fmt.Printf("Error Reading From UDP:%v\n", err2.Error())
							continue
						}
						buf = buf[:n]
						var Result utilities.HeartBeat

						errUnmarshal := json.Unmarshal(buf, &Result)
						if errUnmarshal != nil {
							fmt.Printf("Error Unmarshalling:%v\n", errUnmarshal.Error())
							continue
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
			}

		}
	}()
	return c
}
