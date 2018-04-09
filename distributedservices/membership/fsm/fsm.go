package fsm

import (
	"app/log"
	"app/membership"
	"app/membership/communication"
	"app/membership/utilities"
	"context"
	"errors"
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
	/*
		Find out the ip address of the leader if this machine
		is not the leader.

		If you are a leader publish your ip address over multicast
	*/
	return instance
}

func SendAcknowledgement() chan utilities.Packet {
	//send what we receive on this channel
	ctx := context.Background()
	sendFromThisChan := make(chan utilities.Packet)
	speakChannel := communication.CommSend(ctx, 50001)
	go func() {
		for {
			select {
			case toBeSent := <-sendFromThisChan:
				fmt.Printf("Sending Ack to %v, req:%d\n", toBeSent.ToIp, toBeSent.Req)
				speakChannel <- toBeSent
			}
		}
	}()

	return sendFromThisChan
}

func ReceiveAddRequest() chan utilities.Packet {
	addRequestChannel := make(chan utilities.Packet)
	ctx := context.Background()
	listenChannel := communication.CommReceive(ctx, 50000)
	//speakChannel := communication.CommSend(ctx, 50000)
	go func() {
	TheForLoopState1:
		for {
			select {
			case receivedHbPacket := <-listenChannel:
				if receivedHbPacket.Req == 1 {
					log.Log(ctx, fmt.Sprintf("Received ADD request from IP:%s\n", receivedHbPacket.FromIp.String()))
					ips := utilities.MyIpAddress()
					if len(ips) <= 0 {
						fmt.Printf("Error getting ip\n")
					}
					addRequestChannel <- receivedHbPacket
					/* speakChannel <- utilities.Packet{
						FromIp: ips[0],
						ToIp:   receivedHbPacket.FromIp,
						Seq:    rand.Int63(),
						Req:    2,
					} */
					continue TheForLoopState1
				}
			case <-time.After(100 * time.Millisecond):
				continue TheForLoopState1
			}
		}
	}()

	return addRequestChannel

}

func ReceiveAddAcknowledgement() chan utilities.Packet {
	addAckChannel := make(chan utilities.Packet)
	ctx := context.Background()
	listenChannel := communication.CommReceive(ctx, 50001)
	go func() {
	TheForLoopState1:
		for {
			select {
			case receivedHbPacket := <-listenChannel:
				if receivedHbPacket.Req == 2 {
					log.Log(ctx, fmt.Sprintf("Received ACK from IP:%s\n", receivedHbPacket.FromIp.String()))
					ips := utilities.MyIpAddress()
					if len(ips) <= 0 {
						fmt.Printf("Error getting ip\n")
					}
					addAckChannel <- receivedHbPacket
					/* speakChannel <- utilities.Packet{
						FromIp: ips[0],
						ToIp:   receivedHbPacket.FromIp,
						Seq:    rand.Int63(),
						Req:    2,
					} */
					continue TheForLoopState1
				}
			case <-time.After(100 * time.Millisecond):
				continue TheForLoopState1
			}
		}
	}()

	return addAckChannel

}

func (fsm *Fsm) ProcessFsm() (error, int) {

	switch {
	case fsm.State == 1:
		//Listen for "ADD" requests from peers
		//Forward the request to Membership service
		//Send Ack back to the peer
		addreq := ReceiveAddRequest()
		ackres := SendAcknowledgement()
		for {
			select {
			case addR := <-addreq:
				fmt.Printf("Received %v\n", addR)
				ips := utilities.MyIpAddress()
				if len(ips) <= 0 {
					fmt.Printf("Error getting ip\n")
				}
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
		ctx := context.Background()
		ackChannel := ReceiveAddAcknowledgement()
		speakChannel := communication.CommSend(ctx, 50000)
		for {
			select {
			case ackPacket := <-ackChannel:
				fmt.Printf("Received ACK in STATE 2%v\n", ackPacket)
				/*Send ADD event to Membership.go*/
				/*Move to state 3*/
				return nil, 3
			default:
				ips := utilities.MyIpAddress()
				if len(ips) <= 0 {
					return errors.New("Error accessing local ip address"), 2
				}
				speakChannel <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   net.ParseIP("172.16.238.2"),
					Seq:    rand.Int63(),
					Req:    1,
				}
			}
		}

	case fsm.State == 3:
		fmt.Printf("Moved to state 3\n")
		for {
			time.Sleep(1 * time.Second)
		}
	}
	return nil, fsm.State
}
