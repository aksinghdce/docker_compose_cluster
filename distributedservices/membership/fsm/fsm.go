package fsm

import (
	"app/membership/communication"
	"fmt"
	"app/membership/utilities"
	"context"
	"net"
	"math/rand"
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
	
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	fmt.Println("Error getting hostname")
	// }
	// if hostname == "leader.assignment2" {
	// 	instance.State = 1
	// } else {
	// 	instance.State = 2
	// }	
	return instance
}


func (fsm *Fsm) ProcessFsm() bool {
	switch {
	case fsm.State == 1:
		//Listen for "ADD" requests from peers
		//Forward the request to Membership service
		//Send Ack back to the peer
		ctx := context.Background()
		listenChannel, _ := communication.Comm(ctx, 50000, 50000)
		go func() {
			for i:=0;i<10; i++ {
				fmt.Printf("Received in State 1:%v\n", <-listenChannel)
			}
		}()
		
		ProcessEvent()
	case fsm.State == 2:
		ctx := context.Background()
		_, speakChannel := communication.Comm(ctx, 50000, 50000)
		go func() {
			for i:=0; i<10; i++ {
				fmt.Printf("Sending Packet from State 2\n")
				speakChannel <- utilities.Packet{
					FromIp: net.ParseIP("[::ffff:0.0.0.0]:0"),
					ToIp: net.ParseIP("[::]:0"),
					Seq: rand.Int63(),
					Req: 1,
				}
			}
		}()
		
		ProcessEvent()
	}
	return true
}

/*
This function will send messages to Membership Service
It will ask Membership service for current state so that
it can send that to peers as heartbeat message.

It will send peer's heartbeat messages to Membership service
so that the state can be updated.
*/
func ProcessEvent() {
	//Get channels from Membership service
	//Do plumbing to sendHeartbeat and receiveHeartbeat
	go sendHeartbeat()
	go receiveHeartbeat()
}

func sendHeartbeat() {
	for {

	}
}

func receiveHeartbeat() {
	for {

	}
}












/**********************************************************************************************************/
/*
Specification:

Input: InternalEvent : This carries a context.Context and a sequence number
Output: returns whether the process want to run again to transition state

Processing:
This is the function that runs a finite state machine with 3 states
as described at the beginning of the file.

Receives all the udp datagrams received on
multicast ip address to receive add request

Keeps the add requests in a hashtable
the hash function hashes the ip address.

The hashtable is updating constantly with the last
time of packet arrival

This hashtable is used to construct a sorted
list with ip addresses
*/
// func (fsm *Fsm) Run(intev MainEvent) bool {
// 	switch {
// 	case fsm.State == 1:
// 		ch := multicastheartbeatserver.CatchMultiCastDatagramsAndBounce(intev.Ctx, "224.0.0.1", "10001")
// 		sortedIps := erm.SortCurrentGroupInfo()
// 		channelArr := make([]chan utilities.HeartBeat, 0)
// 		for _, ip := range sortedIps {
// 			channelArr = append(channelArr, multicastheartbeater.SendHeartBeatMessages(intev.Ctx, ip, "50012"))
// 		}

// 		for {
// 			timeout := time.After(SEND_HEARTBEAT_EVERY)

// 			select {
// 			case s := <-ch:
// 				erm.AddNodeToGroup(intev, s.FromTo.FromIp)
// 				if s.ReqCode == 1 {
// 					erm.SendAckToAddRequester(intev, s.FromTo.FromIp, "50009")
// 				}
// 			case <-timeout:
// 				for i, chpeer := range channelArr {
// 					chpeer <- utilities.HeartBeat{
// 						Cluster:   erm.GroupInfo,
// 						ReqNumber: intev.RequestNumber.get(),
// 						ReqCode:   3, //1 is for ADD request
// 						FromTo: utilities.MessageAddressVector{
// 							FromIp: erm.MyState.MyIp,
// 							ToIp: sortedIps[i],
// 						},
// 					}
// 				}
// 			}
// 		}
// 	case erm.MyState.CurrentState == 2:
// 		heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, "224.0.0.1", "10001")

// 		heartbeatChannelIn := multicastheartbeatserver.CatchUniCastDatagramsAndBounce(intev.Ctx, "50009")

// 		for {
// 			hbMessage := utilities.HeartBeat{
// 				Cluster:   nil,
// 				ReqNumber: intev.RequestNumber.get(),
// 				ReqCode:   1, //1 is for ADD request
// 				FromTo: utilities.MessageAddressVector{
// 					FromIp: "",
// 					ToIp: "224.0.0.1",
// 				},
// 			}

// 			intev.RequestNumber.increment()

// 			select {
// 			case hbRcv := <-heartbeatChannelIn:
// 				ip_port := strings.Split(hbRcv.FromTo.ToIp, ":")
// 				erm.MyState.MyIp = ip_port[0]
// 				ip_port_leader := strings.Split(hbRcv.FromTo.FromIp, ":")
// 				erm.LeaderUniCastIp = ip_port_leader[0]
// 				fmt.Printf("Updated myIp to:%v and Leader unicast to: %v\n", erm.MyState.MyIp, erm.LeaderUniCastIp)
// 				if hbRcv.ReqCode == 2 {
// 					utilities.Log(intev.Ctx, "STATE Transition 2->3\n")
// 					erm.MyState.CurrentState = 3
// 					erm.GroupInfo = hbRcv.Cluster
// 					erm.LastHeartbeatReceived = hbRcv
// 					// Ask the caller to rerun this function: To change state to 3
// 					return true
// 				}
// 			default:
// 				heartbeatChannelOut <- hbMessage
// 			}
// 		}
// 	case erm.MyState.CurrentState == 3:
// 		heartbeatChannelToListener := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, "224.0.0.1", "10001")
// 		heartbeatChannelIn := multicastheartbeatserver.CatchUniCastDatagramsAndBounce(intev.Ctx, "50012")
// 		sendTo, err := erm.WhomToSendHb()
// 		if err != nil {
// 			utilities.Log(intev.Ctx, err.Error())
// 		}
// 		var heartbeatChannelOut chan utilities.HeartBeat
			
// 		if len(sendTo) > 0 {
// 			heartbeatChannelOut = multicastheartbeater.SendHeartBeatMessages(intev.Ctx, sendTo, "50012")
// 		}else {
// 			fmt.Print("blank sendTo")
// 		}

// 		for {
// 			timeout := time.After(SEND_HEARTBEAT_EVERY)
// 			select {
// 			case hbst := <-heartbeatChannelIn:
// 				if hbst.ReqCode == 3 {
// 					erm.ConsolidateInfo(hbst.Cluster)
// 				}
// 				ip_port := strings.Split(hbst.FromTo.ToIp, ":")
// 				erm.MyState.MyIp = ip_port[0]
// 				erm.DeleteOlderHeartbeats(DELETE_OLDER_THAN)
// 				erm.AddNodeToGroup(intev, hbst.FromTo.FromIp)
// 			case <-timeout:
// 				if len(sendTo) > 0 {
// 					heartbeatChannelOut <- utilities.HeartBeat{
// 						Cluster:   erm.GroupInfo,
// 						ReqNumber: intev.RequestNumber.get(),
// 						ReqCode:   3, //1 is for ADD request
// 						FromTo: utilities.MessageAddressVector{
// 							FromIp: erm.MyState.MyIp,
// 							ToIp: sendTo,
// 						},
// 					}
// 				}

// 				heartbeatChannelToListener <- utilities.HeartBeat{
// 					Cluster:   erm.GroupInfo,
// 					ReqNumber: intev.RequestNumber.get(),
// 					ReqCode:   3, //1 is for ADD request
// 					FromTo: utilities.MessageAddressVector{
// 						FromIp: erm.MyState.MyIp,
// 						ToIp: erm.LeaderUniCastIp,
// 					},
// 				}
// 			}

// 		}
// 	}
// 	return false
// }