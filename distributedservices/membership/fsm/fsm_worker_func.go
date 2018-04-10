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

	// a control to stop sending ADD request to Leader
	done := make(chan bool)
	channel := communication.GetComm()("send", 50000)
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
					channel.ControlC <- true
					break TheForLoopSendAdd
				}
			default:
				//fmt.Printf("Sending ADD req to LEADER\n")
				channel.DataC <- utilities.Packet{
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
	// a control to stop sending ADD request to Leader
	done := make(chan bool)
	channel := communication.GetComm()("send", 50001)
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
					channel.ControlC <- true
					break TheForLoopSendHB
				}
			default:
				channel.DataC <- utilities.Packet{
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
	sendFromThisChan := make(chan utilities.Packet)
	channel := communication.GetComm()("send", 50001)
	go func() {
	LoopSendAck:
		for {
			select {
			case toBeSent := <-sendFromThisChan:
				channel.DataC <- toBeSent
				time.Sleep(time.Second)
				channel.ControlC <- true
				break LoopSendAck
			}
		}
		return
	}()
	return sendFromThisChan
}

func ReceiveAddRequest() chan utilities.Packet {
	addRequestChannel := make(chan utilities.Packet)
	channel := communication.GetComm()("receive", 50000)
	go func() {
	TheForLoopState1:
		for {
			select {
			case receivedHbPacket := <-channel.DataC:
				if receivedHbPacket.Req == 1 {
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
	channel := communication.GetComm()("receive", 50001)
	go func() {
	TheForLoopState1:
		for {
			select {
			case receivedHbPacket := <-channel.DataC:
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
	channel := communication.GetComm()("receive", 50001)
	go func() {
	TheForLoopState1:
		for {
			select {
			case receivedHbPacket := <-channel.DataC:
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
				channel.ControlC <- true
				break TheForLoopState1
			}
		}
	}()
	return addAckChannel, stop
}
