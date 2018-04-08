package fsm

import (
	"app/membership/communication"
	"app/membership"
	"fmt"
	"app/membership/utilities"
	"context"
	"math/rand"
	"errors"
	"app/log"
	"time"
	"net"
)

/*This module is responsible for managing heartbeats
*/

type Fsm struct {
	State int
	Mserv membership.Membership
}

func Init(initialState int) *Fsm {
	instance := &Fsm{
		State : initialState,
		Mserv : membership.Membership{},
	}
	/*
	Find out the ip address of the leader if this machine
	is not the leader.

	If you are a leader publish your ip address over multicast
	*/
	return instance
}


func (fsm *Fsm) ProcessFsm() (error, int){
	chout, chin := fsm.Mserv.KeepMembershipUpdated()
	switch {
	case fsm.State == 1:
		//Listen for "ADD" requests from peers
		//Forward the request to Membership service
		//Send Ack back to the peer
		
		ctx := context.Background()
		listenChannel, speakChannel := communication.Comm(ctx, 50001, 50002)
		go func() error{
			for {
				fmt.Printf("State 1 Listening \n")
				select {
				case receivedHbPacket := <-listenChannel:
					fmt.Printf("Recceived %v\n", receivedHbPacket)
					if receivedHbPacket.Req == 1 {
						log.Log(ctx, fmt.Sprintf("Received ADD request from IP:%s\n", receivedHbPacket.FromIp.String()))
						ips := utilities.MyIpAddress()
						if len(ips) <= 0 {
							return errors.New("Error accessing local ip address")
						}
						fmt.Printf("State 1: Ip: %v\n", ips[0])
						speakChannel <- utilities.Packet{
							FromIp: ips[0],
							ToIp: receivedHbPacket.FromIp,
							Seq: rand.Int63(),
							Req: 2,
						}
						chin <- receivedHbPacket	
					}
				case membershipOutgoing := <-chout:
					fmt.Printf("Membership service wants to send:%v\n", membershipOutgoing)
				}
			}
			return nil
		}()

		 /*This go routine sends and receives Heartbeat packets*/
		 /*
		listenChannel2, speakChannel2 := communication.Comm(ctx, 50001, 50002)
		go func() error{
			for {
				select {
				case receivedHbPacket := <-listenChannel2:
					if receivedHbPacket.Req == 3 {
						log.Log(ctx, fmt.Sprintf("Received HB from:%s\n", receivedHbPacket.FromIp.String()))
						ips := utilities.MyIpAddress()
						if len(ips) <= 0 {
							return errors.New("Error accessing local ip address")
						}
						fmt.Printf("State 1': Ip: %v\n", ips[0])
						
						speakChannel2 <- utilities.Packet{
							FromIp: ips[0],
							ToIp: receivedHbPacket.FromIp,
							Seq: rand.Int63(),
							Req: 3,
						}	
					}
				}
			}
			return nil
		}() */
	case fsm.State == 2:
		/*State 2 is transient state to send ADD request to Leader and Wait for Ack*/
		ctx := context.Background()
		listenChannel, speakChannel := communication.Comm(ctx, 50002, 50001)
		for{
			select {
			case receivedHbPacket := <-listenChannel :
				if receivedHbPacket.Req == 2 {
					fmt.Printf("Received ACK in STATE 2%v\n", receivedHbPacket)
					/*Send ADD event to Membership.go*/
					/*Move to state 3*/
					return nil, 3
				}
			case <-time.After(1 * time.Second):
				fmt.Printf("Send ADD request from State 2\n")
				ips := utilities.MyIpAddress()
				if len(ips) <= 0 {
					return errors.New("Error accessing local ip address"), 2
				}
				fmt.Printf("State 2: Ip: %v\n", ips[0])
				speakChannel <- utilities.Packet{
					FromIp: ips[0],
					ToIp: net.ParseIP("172.16.238.2"),
					Seq: rand.Int63(),
					Req: 1,
				}
			}
			
		}
		
	case fsm.State == 3:
		fmt.Printf("Moved to state 3\n")
		/* ctx := context.Background()
		listenChannel2, _ := communication.Comm(ctx, 50002, 50001)
		go func() error{
			for {
				select {
				case receivedHbPacket := <-listenChannel2:
					if receivedHbPacket.Req == 3 {
						fmt.Printf("Received Heartbeat %v\n", receivedHbPacket)	
						
					}
				
				}
			}
			return nil
		}() */
	}
	return nil, fsm.State
}