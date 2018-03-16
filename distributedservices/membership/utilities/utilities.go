package utilities

import (
	"sync/atomic"
	"net"
	"os"
	"fmt"
)

/*
These utilities are meant for the use by membership service.
All the utilities must be tested for accuracy and speed
*/

/*
An internal counter for requests to the network
To maintain a vector clock in association with hostname

For every event that the node sends or receives as a membership
service the Sequence will be incremented atomicly.

This counter must be managed by the membership service, 
We will follow the philosophy of single point of maintainance.
http://www.ifsq.org/single-point-of-maintenance.html

*/

/************************************************************************************
Sequence counter
************************************************************************************/
type count64 int64

type vectorClock struct{
	Sequence int64
	Hostname string
}

func (c *count64) increment() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *count64) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

/************************************************************************************
A Node in the chord;s extended ring
************************************************************************************/
type Node struct {
	Ip string
	Alive bool
}

/************************************************************************************
A heartbeat packet
************************************************************************************/
type Packet struct {
	FromIp net.IP //the fromip 
	ToIp net.IP //and toip required - udp packets
	//Sequence number is required for ordering
	//If we need to log and distinguish heartbeats
	Seq int64
	/*
		stores the extended ring with the 
		heartbeat status
	*/
	Req int
	Data []Node
}

func MyIpAddress() []net.IP{
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		os.Exit(1)
	}
	ips := make([]net.IP, 0)
	for _, i := range ifaces {
    	addrs, err := i.Addrs()
		// handle err
		if err != nil {
			fmt.Printf("Error:%v\n", err)
			os.Exit(1)
		}
		var ip net.IP
    	for _, addr := range addrs {
        		switch v := addr.(type) {
        		case *net.IPNet:
                	ip = v.IP
        		case *net.IPAddr:
            	    ip = v.IP
			}
			if (ip.String() != "127.0.0.1") {
				ips = append(ips, ip)
			}
    	}
	}
	return ips
}