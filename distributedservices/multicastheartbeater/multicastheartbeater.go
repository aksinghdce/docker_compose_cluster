package multicastheartbeater

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
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

func SendMultiCastAddRequest(toAddress, toPort string) chan utilities.HeartBeat {
	heartbeatChannelIn := make(chan utilities.HeartBeat)
	go func() {
		var LeaderAddress string
		LeaderAddress = toAddress
		LeaderAddress += ":"
		LeaderAddress += toPort
		ServerAddr, err := net.ResolveUDPAddr("udp", LeaderAddress)

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
}

/*
Specification:
Returns a channel of utilities.HeartBeat
The caller can read heartbeats on this channel at the speed that UDP
provides; with a time lag associated with go channels
*/
func SendHeartBeatMessages(toAddress, toPort string, fromPort string) chan utilities.HeartBeat {
	heartbeatChannelIn := make(chan utilities.HeartBeat)
	go func() {
		for {
			select {
			case hbin := <-heartbeatChannelIn:
				ifsArr, err := net.Interfaces()
				if err != nil {
					fmt.Print(err.Error())
				}
				for _, ifs := range ifsArr {
					//FlagPointToPoint
					flag := ifs.Flags.String()
					if strings.Contains(flag, "up") && !strings.Contains(flag, "loopback") {
						unicastaddresses, err := ifs.Addrs()
						if err != nil {
							continue
						}

						for _, uniaddr := range unicastaddresses {
							if !strings.Contains(toAddress, ":") {
								toAddress += ":"
								toAddress += toPort
							}
							toAddr, err := net.ResolveUDPAddr("udp", toAddress)
							if err != nil {
								break
							}

							SenderPort := ":" + fromPort
							uniaddrStr := uniaddr.String()

							indexOfSlash := strings.Index(uniaddrStr, "/")
							if indexOfSlash > 0 {
								arr := strings.Split(uniaddrStr, "/")
								uniaddrStr = arr[0]
							}

							uniaddrStr += SenderPort
							fmt.Print("GOT STOP ADD MESSAGE TO SEND LOWER LEVEL\n")
							fromAddr, err := net.ResolveUDPAddr("udp", uniaddrStr)
							if err != nil {
								fmt.Printf("ERRRRRRRORRRRR:%s\n", err.Error())
								break
							}

							Conn, err := net.DialUDP("udp", fromAddr, toAddr)
							if err != nil {
								fmt.Print("3333333333333333", err.Error())
								break
							}
							defer Conn.Close()

							//encode json data
							//fmt.Printf("Data to be Sent:%v\n", hb)
							jsonData, err := json.Marshal(hbin)
							fmt.Printf("Marshalled Data:%v\n", string(jsonData))
							_, err = Conn.Write(jsonData)
							if err != nil {
								fmt.Println(err.Error())
							}
							//fmt.Printf("Wrote %d bytes\n", n)
							//time.Sleep(time.Second * 1)

						}
					}

				}
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}

	}()
	return heartbeatChannelIn
}
