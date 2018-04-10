package fsm

import (
	"app/membership"
	"app/membership/communication"
	"app/membership/utilities"
	"fmt"
	"math/rand"
	"net"
	"time"
)

/*This module is responsible for managing heartbeats
 */

type Fsm struct {
	State int
	Mserv membership.Membership
}

func Init(initialState int) *Fsm {
	instance := &Fsm{
		State: initialState,
		Mserv: membership.Membership{},
	}
	return instance
}

func (fsm *Fsm) ProcessFsm() (error, int) {

	switch {
	case fsm.State == 1:
		//Listen for "ADD" requests from peers
		//Forward the request to Membership service
		//Send Ack back to the peer
		addreq := ReceiveAddRequest()
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error getting ip\n")
		}
		for {
			select {
			case addR := <-addreq:
				//fmt.Printf("Received %v\n", addR)
				ackres := SendAcknowledgement()
				ackres <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   addR.FromIp,
					Seq:    rand.Int63(),
					Req:    2,
				}
			}
		}
	case fsm.State == 2:
		/*State 2 is transient state to send ADD request to Leader and Wait for Ack*/
		//Keep sending "ADD" request to leader
		SendAddReqToLeader()
		channel := communication.GetComm()("receive", 50001)
	LoopState2:
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case ack := <-channel.DataC:
				if ack.Req == 2 {
					fmt.Printf("ack received:%v\n", ack)
					//channel.ControlC <- true
					break LoopState2
				}
			}
		}
		return nil, 3

	case fsm.State == 3:
		fmt.Printf("Moved to state 3\n")
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error getting ip\n")
			return nil, fsm.State
		}
		channelS := communication.GetComm()("send", 50002)
		channelR := communication.GetComm()("receive", 50002)
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				channelS.DataC <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   net.ParseIP("172.16.238.6"),
					Seq:    rand.Int63(),
					Req:    3,
				}
			}
		}()

		go func() {
			for {
				select {
				case hbR := <-channelR.DataC:
					fmt.Printf("Received %v\n", hbR)
				}
			}
		}()
	}
	return nil, fsm.State
}
