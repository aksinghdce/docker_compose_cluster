package fsm

import (
	"app/log"
	"app/membership/communication"
	"app/membership/utilities"
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func SendAddReqToLeader() chan bool {
	//command to stop sending ADD req
	ctx := context.Background()
	// a control to stop sending ADD request to Leader
	done := make(chan bool)
	speakChannel, stop_speaking := communication.CommSend(ctx, 50000)
	go func() {
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error getting ip\n")
			return
		}
	TheForLoopSendAdd:
		for {
			select {
			case stop := <-done:
				if stop == true {
					fmt.Printf("Stopping to send ADD\n")
					stop_speaking <- true
					break TheForLoopSendAdd
				}
			default:
				//fmt.Printf("Sending ADD req to LEADER\n")
				speakChannel <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   net.ParseIP("172.16.238.2"),
					Seq:    rand.Int63(),
					Req:    1,
				}
			}

		}
	}()
	return done
}

func SendHbReqToIp(ip net.IP) chan bool {
	//command to stop sending ADD req
	ctx := context.Background()
	// a control to stop sending ADD request to Leader
	done := make(chan bool)
	speakChannel, stop_speaking := communication.CommSend(ctx, 50001)
	go func() {
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error getting ip\n")
			return
		}
	TheForLoopSendHB:
		for {
			//Send HB every 100 milliseconds
			//<-time.After(100 * time.Millisecond)
			select {
			case stop := <-done:
				if stop == true {
					fmt.Printf("Stopping to send HB\n")
					stop_speaking <- true
					break TheForLoopSendHB
				}
			default:
				speakChannel <- utilities.Packet{
					FromIp: ips[0],
					ToIp:   ip,
					Seq:    rand.Int63(),
					Req:    3,
				}
			}

		}
	}()
	return done
}

func SendAcknowledgement() chan utilities.Packet {
	//send what we receive on this channel
	ctx := context.Background()
	sendFromThisChan := make(chan utilities.Packet)
	speakChannel, stop_speaking := communication.CommSend(ctx, 50001)
	go func() {
	LoopSendAck:
		for {
			select {
			case toBeSent := <-sendFromThisChan:
				speakChannel <- toBeSent
				time.Sleep(time.Second)
				stop_speaking <- true
				break LoopSendAck
			}
		}
		return
	}()
	return sendFromThisChan
}

func ReceiveAddRequest() chan utilities.Packet {
	addRequestChannel := make(chan utilities.Packet)
	ctx := context.Background()
	listenChannel, _ := communication.CommReceive(ctx, 50000)
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
					continue TheForLoopState1
				}
			}
		}
	}()
	return addRequestChannel
}

func ReceiveHbRequest() chan utilities.Packet {
	hbRequestChannel := make(chan utilities.Packet)
	ctx := context.Background()
	listenChannel, _ := communication.CommReceive(ctx, 50001)
	go func() {
	TheForLoopState1:
		for {
			select {
			case receivedHbPacket := <-listenChannel:
				if receivedHbPacket.Req == 3 {
					log.Log(ctx, fmt.Sprintf("Received Hb request from IP:%s\n", receivedHbPacket.FromIp.String()))
					ips := utilities.MyIpAddress()
					if len(ips) <= 0 {
						fmt.Printf("Error getting ip\n")
					}
					hbRequestChannel <- receivedHbPacket
					continue TheForLoopState1
				}
			}
		}
	}()
	return hbRequestChannel
}

func ReceiveAddAcknowledgement() (chan utilities.Packet, chan bool) {
	addAckChannel := make(chan utilities.Packet)
	stop := make(chan bool)
	ctx := context.Background()
	listenChannel, over := communication.CommReceive(ctx, 50001)
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
					continue TheForLoopState1
				}
			case _ = <-stop:
				over <- true
				break TheForLoopState1
			}
		}
	}()
	return addAckChannel, stop
}
