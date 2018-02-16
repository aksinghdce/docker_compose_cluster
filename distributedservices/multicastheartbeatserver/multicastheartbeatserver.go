package multicastheartbeatserver

import (
	"fmt"
	"net"
)

/*
Leader wants to listen on this port for heartbeats
And listen to serious data and control requests on
TCP protocol 8080
*/
const LEADER_MULTICAST_UDP_PORT_STRING = ":10001"

/*
Specification

Find out if I am not the leader node.
If I am not the leader node I take a different course of action
(not implemented yet)

If I am the leader then I listen to pings and ask my local
http server to update it's internal representation of cluster.
*/
func CatchDatagramsAndBounce(c chan string) {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	for {
		addrstr := string("224.0.0.1")
		addrstr += LEADER_MULTICAST_UDP_PORT_STRING
		udpaddr, err := net.ResolveUDPAddr("udp", addrstr)
		if err != nil {
			fmt.Println("Fault 2: ", err)
		}

		conn, err := net.ListenMulticastUDP("udp", nil, udpaddr)
		if err != nil {
			fmt.Println("Fault 4: ", err)
		}

		defer conn.Close()
		buf := make([]byte, 256)
		_, udpadd, err2 := conn.ReadFromUDP(buf)
		if err2 != nil {
			fmt.Println("Fault 5: ", err2)
		}
		//fmt.Printf("Received from udp client: %s", string(buf))
		output := string(buf)
		output += udpadd.String()
		c <- output
	}
}
