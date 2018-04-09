package fsm

import (
	"app/membership"
	"app/membership/utilities"
	"fmt"
	"math/rand"
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
		ackChannel, _ := ReceiveAddAcknowledgement()
	LoopState2:
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case ack := <-ackChannel:
				fmt.Printf("ack received:%v\n", ack)
				break LoopState2
			}
		}
		return nil, 3

	case fsm.State == 3:
		fmt.Printf("Moved to state 3\n")
	}
	return nil, fsm.State
}
