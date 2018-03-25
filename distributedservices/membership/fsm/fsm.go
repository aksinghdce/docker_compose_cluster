package fsm

import (
	"app/membership/communication"
	"fmt"
	"app/membership/utilities"
	"context"
	"math/rand"
	"errors"
	"app/log"
)

/*This module is responsible for managing heartbeats
*/

type Fsm struct {
	State int
}

func Init(initialState int) *Fsm {
	instance := &Fsm{
		State : initialState,
	}
	return instance
}


func (fsm *Fsm) ProcessFsm() (error, int){
	
	switch {
	case fsm.State == 1:
		//Listen for "ADD" requests from peers
		//Forward the request to Membership service
		//Send Ack back to the peer
		ctx := context.Background()
		listenChannel, speakChannel := communication.Comm(ctx, 10001, 10002)
		go func() error{
			for {
				select {
				case receivedHbPacket := <-listenChannel:
					fmt.Printf("Recceived %v\n", receivedHbPacket)
					if receivedHbPacket.Req == 1 {
						log.Log(ctx, fmt.Sprintf("Received ADD request from IP:%s\n", receivedHbPacket.FromIp.String()))
						ips := utilities.MyIpAddress()
						if len(ips) <= 0 {
							return errors.New("Error accessing local ip address")
						}
						speakChannel <- utilities.Packet{
							FromIp: ips[0],
							ToIp: receivedHbPacket.FromIp,
							Seq: rand.Int63(),
							Req: 2,
						}	
					}
				}
			}
			return nil
		}()

		/*This go routine sends and receives Heartbeat packets*/
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
						/*Send HB back*/
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
		}()
	case fsm.State == 2:
		/*State 2 is transient state to send ADD request to Leader and Wait for Ack*/
		ctx := context.Background()
		listenChannel, speakChannel := communication.Comm(ctx, 10002, 10001)
		for{
			select {
			case receivedHbPacket := <-listenChannel :
				if receivedHbPacket.Req == 2 {
					fmt.Printf("Received ACK in STATE 2%v\n", receivedHbPacket)
					/*Send ADD event to Membership.go*/
					/*Move to state 3*/
					return nil, 3
				}
			default:
				fmt.Printf("Send ADD request from State 2\n")
				ips := utilities.MyIpAddress()
				if len(ips) <= 0 {
					return errors.New("Error accessing local ip address"), 2
				}
				speakChannel <- utilities.Packet{
					FromIp: ips[0],
					ToIp: ips[0],
					Seq: rand.Int63(),
					Req: 1,
				}
			}
			
		}
		
	case fsm.State == 3:
		fmt.Printf("Moved to state 3\n")
		ctx := context.Background()
		listenChannel2, _ := communication.Comm(ctx, 50002, 50001)
		go func() error{
			for {
				select {
				case receivedHbPacket := <-listenChannel2:
					if receivedHbPacket.Req == 3 {
						fmt.Printf("Received Heartbeat %v\n", receivedHbPacket)	
						/*Inform Membership.go*/
					}
				/*Find out the ip address of machines to send HB to*/
				}
			}
			return nil
		}()
	}
	return nil, fsm.State
}