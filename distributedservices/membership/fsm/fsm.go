package fsm

import (
	"app/membership"
	"app/membership/communication"
	"app/membership/utilities"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

/*This module is responsible for managing heartbeats
 */
const HeartBeatInterval = 2 * time.Second

type Fsm struct {
	State int
	Mserv membership.Membership
}

func Init(initialState int) *Fsm {
	instance := &Fsm{
		State: initialState,
		Mserv: membership.Membership{
			ChanOut: make(chan utilities.Packet),
			ChanIn:  make(chan utilities.Packet),
		},
	}

	return instance
}

func (fsm *Fsm) ProcessFsm() error {
	ChanRm, ChanSm := fsm.Mserv.KeepMembershipUpdated()
	switch {
	case fsm.State == 1:
		//Listen for "ADD" requests from peers
		addreq := ReceiveAddRequest()

		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error getting ip\n")
		}
		for {
			select {
			case addR := <-addreq:
				//Send Ack back to the peer
				ackres := SendAcknowledgement()
				ackres <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   addR.FromIp,
					Seq:    rand.Int63(),
					Req:    2,
				}
				fmt.Printf("Received ADD request\n")
				ChanSm <- addR
			case recvM := <-ChanRm:
				fmt.Printf("Membership to fsm:%v\n", recvM)
			}
		}
	case fsm.State == 2:
		/*State 2 is a state for non-leader node
		send ADD request to Leader and Wait for Ack*/
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
					// ask SendAddReqToLeader() generator to stop sending
					// ADD req.
					//channel.ControlC <- true
					//Fix it later: SendAddReqToLeader is still sending ADD req
					//stop <- true
					break LoopState2
				}
			}
		}

		/*Regular heartbeat begin here*/

		fmt.Printf("Regular heartbeats begin\n")
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			return errors.New("Can't get IP address\n")
		}
		channelS := communication.GetComm()("send", 50001)

		go func() {
			for {
				time.Sleep(HeartBeatInterval)
				channelS.DataC <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   net.ParseIP("172.16.238.4"),
					Seq:    rand.Int63(),
					Req:    3,
				}
			}
		}()

		go func() {
			for {
				select {
				case hbR := <-channel.DataC:
					fmt.Printf("Received %v\n", hbR)
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}
	return errors.New("Shouldn't have returned from here\n")
}
