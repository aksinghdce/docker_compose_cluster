package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/*
Check generic error.
*/
func CheckError(err error) {
	if err != nil {
		fmt.Println("Exiting with Error: ", err)
		os.Exit(1)
	}
}

/*
Leader wants to listen on this port for heartbeats
And listen to serious data and control requests on
TCP protocol 8080
*/
const LEADER_MULTICAST_UDP_PORT_STRING = ":10001"

/*
Multicast address: 172.20.0.1 ?
*/
func main() {
	/*What's my hostname?
	 */
	hostname, err := os.Hostname()
	CheckError(err)
	fmt.Printf("My hostname:%s\n", hostname)

	// If my hostname is not leader.assignment2
	// I should just be doing the service and not lead
	// I will ask the leader to add me to his group of
	// workers until he doesn't add me. If he says no, I
	// should just die.
	if hostname != "leader.assignment2" {
		/*
			We can branch out of here for non-leaders to send ADD request to the leader
			Cases:
				1. If the leader doesn't respond to ping datagrams: Leader is dead
				2. If the leader does respond to ping datagrams but denies ADD request
				3. Leader is alive, does ADD me and tells me my responsibilities of pinging
				a few nodes and reporting about their death.
					If I report about the death of a node for whom I am responsible
						then the leader gives me more nodes to look after
		*/
		fmt.Printf("I am not a leader. I am too old to serve. I will just die")
		os.Exit(0)
	}

	/*Get all the interfaces on the machine*/
	ifsArr, err := net.Interfaces()
	CheckError(err)

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
			CatchDatagramsAndBounce(multicastaddresses, ifs)
		}

	}

}

/*
Specification

Find out if I am not the leader node.
If I am not the leader node I take a different course of action
(not implemented yet)

If I am the leader then I listen to pings and ask my local
http server to update it's internal representation of cluster.
*/
func CatchDatagramsAndBounce(multicastaddresses []net.Addr, ifs net.Interface) {
	/*I am the leader I will keep listening to events
	from my group.*/
	for {
		for index, addr := range multicastaddresses {
			fmt.Println("Network:", addr.Network())
			addstr := addr.String()
			addstr += LEADER_MULTICAST_UDP_PORT_STRING
			udpaddr, err := net.ResolveUDPAddr("udp", addstr)
			if err != nil {
				fmt.Println("Fault 2: ", err)
				/*We will continue until we find a multicast UDP
				port that works*/
				continue
			}

			fmt.Printf("multicast address %d : %s\n", index, addr.String())
			//func ListenMulticastUDP(network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)
			interf, err := net.InterfaceByIndex(ifs.Index)
			if err != nil {
				fmt.Println("Fault 3: ", err)
				/*We will continue until we find a multicast UDP
				port that works*/
				continue
			}
			conn, err := net.ListenMulticastUDP("udp", interf, udpaddr)
			if err != nil {
				fmt.Println("Fault 4: ", err)
				/*We will continue until we find a multicast UDP
				port that works*/
				continue
			}
			defer conn.Close()

			buf := make([]byte, 256)
			_, udpaddr, err2 := conn.ReadFromUDP(buf)
			if err2 != nil {
				fmt.Println("Fault 5: ", err2)
				continue
			}

			/*
				Got something interesting: someone wants to join the group.
				How to deal with this new request?

				Send it to the http servemux
			*/
			reportPing := fmt.Sprintf("From %s Data received: %s\n", udpaddr.String(), string(buf))
			//strings.NewReader(cmd.Encode())
			cmd := url.Values{}
			cmd.Add("add", reportPing)
			req, err := http.NewRequest("POST", "http://leader.assignment2:8080/membership/add", strings.NewReader(cmd.Encode()))
			ctx := context.Background()
			// Don't wait for more than a second to get the grep result from remote server
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			req2 := req.WithContext(ctx)
			resp, err := http.DefaultClient.Do(req2)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				log.Println("ERROR: sending request to leader.assignment2")
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading response from remote")
			}

			fmt.Printf("Body:%v\n", string(body))
			/*
				Handled the request now will will just break out of it and loop
				throught he whole exercise to listen to the next ADD request or
				other such request.
			*/
			break
		}
	}
}
