package membership

import (
	"container/list"
	"app/membership/utilities"
	"fmt"
	"net"
	"math/rand"
)

/*
1. This structure will just have one state that it has to maintain
consistent. 

2. Heartbeat can be calculated from a consistent state
only.

3. This service will expose a method that the users of this package
can call, this method will exclusively be present in this file for
maintainability.
*/
type Membership struct{
	chanOut chan utilities.Packet
	chanIn chan utilities.Packet
	ring list.List
}

/*
Membership Service is invoked by fsm 
KeepMembershipUpdated will expect events
coming on an incoming channel from fsm engine

If the state of ring changes, then membership
service will send it on outgoing channel.
*/
func (m *Membership) KeepMembershipUpdated() (chan utilities.Packet, chan utilities.Packet) {
	packet := utilities.Packet{
		FromIp: net.ParseIP("127.0.0.1"),
		ToIp: net.ParseIP("127.0.0.2"),
		Seq: rand.Int63(),
	}
	go func() {
		/*This go routine will listen for incoming packet.
		If the incoming packet is an "ADD" request, it will update the extended ring
		If the incoming packet is a "REMOVE" request, it will update the extended ring
		If the incoming packet denotes heartbeat miss, then it will update the extended ring
		
		After updating the ring, it will send heartbeat messages to a subset of nodes in the
		extended ring*/
		for {
			m.chanOut <- packet
			fmt.Printf("Received:%v\n", <-m.chanIn)
		}
	}()
	return m.chanOut, m.chanIn
}

func (m *Membership) Insert(key string, value interface{}) bool {
	return false
}

func (m *Membership) Get(key string) interface{} {
	return false
}