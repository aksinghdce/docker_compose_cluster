package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Exiting with Error: ", err)
		os.Exit(1)
	}
}

/*
Multicast address: 172.20.0.1 ?
*/
func main() {
	ifsArr, err := net.Interfaces()
	CheckError(err)

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
			for {
				for index, addr := range multicastaddresses {
					fmt.Println("Network:", addr.Network())
					addstr := addr.String()
					addstr += ":10001"
					udpaddr, err := net.ResolveUDPAddr("udp", addstr)
					if err != nil {
						fmt.Println("Fault 2: ", err)
						continue
					}

					fmt.Printf("multicast address %d : %s\n", index, addr.String())
					//func ListenMulticastUDP(network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)
					interf, err := net.InterfaceByIndex(ifs.Index)
					if err != nil {
						fmt.Println("Fault 3: ", err)
						continue
					}
					conn, err := net.ListenMulticastUDP("udp", interf, udpaddr)
					if err != nil {
						fmt.Println("Fault 4: ", err)
						continue
					}
					defer conn.Close()

					buf := make([]byte, 256)
					_, udpaddr, err2 := conn.ReadFromUDP(buf)
					if err2 != nil {
						fmt.Println("Fault 5: ", err2)
						continue
					}
					fmt.Printf("From %s Data received: %s\n", udpaddr.String(), string(buf))
					break
				}
			}

		}

	}

}
