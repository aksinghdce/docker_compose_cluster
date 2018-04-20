package fsm

import (
	"app/membership"
	"app/membership/communication"
	"app/membership/utilities"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

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
	_, ChanSm := fsm.Mserv.KeepMembershipUpdated()
	switch {
	case fsm.State == 1:
		add_raised_panic := make(chan string)

		//ADD Process Loop : Listen for "ADD" requests and send "ACK"
		//GO ROUTINE TO RECEIVE ADD REQUESTS AND SEND ACK
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered in ProcessFsm STATE 1:%v !!", r)
					add_raised_panic <- "ADD request handler dying\n"
				}
			}()
			restart_comm_receive := true
			ips := utilities.MyIpAddress()
			if len(ips) <= 0 {
				fmt.Printf("Error getting ip\n")
			}
			for {
				var addreq_data_control communication.Channels2
				var data0 utilities.Packet
				if restart_comm_receive {
					addreq_data_control = communication.GetComm2()("receive", 50000)
					data0 = <-addreq_data_control.DataC
				}
				//fmt.Fprintf(os.Stdout, "ADD request received from:%v\n", data0.FromIp)
				//addreq_data_control.ControlC <- "stop receiving ADD and exit"
				//fmt.Printf("ADD-request-receiving-channel closed:%v\n", <-addreq_data_control.ControlC)
				restart_comm_receive = true
				ackres_data_control := communication.GetComm2()("send", 50001) //Send ACK 3 UDP packets for ensuring receipt
				for i := 0; i < 3; i++ {
					ackres_data_control.DataC <- utilities.Packet{
						FromIp: ips[0],
						ToIp:   data0.FromIp,
						Seq:    rand.Int63(),
						Req:    2,
					}
				}
				//ackres_data_control.ControlC <- "Stop sending ACK"
				//fmt.Printf("So done stopping:%v\n", <-ackres_data_control.ControlC)
				fmt.Printf("ACK sent 3 times to:%v\n", data0.FromIp)
				ChanSm <- data0
				fmt.Printf("ADD request forwarded to Membership service:%v\n", data0.FromIp)

				select {
				case addR := <-addreq_data_control.DataC: //ADD request received
					data0 = addR
				case <-time.After(6 * time.Second):
					restart_comm_receive = false
				}
				//Go for next ADD request
			}
		}()

		/* 		go func() {
		   			defer func() {
		   				if r := recover(); r != nil {
		   					fmt.Printf("Recovered in ProcessFsm STATE 1:%v !!", r)
		   					hb_raised_panic <- "HB request handler dying\n"
		   				}
		   			}()
		   		HB_Loop:
		   			for {
		   				hbreq_data_control := communication.GetComm2()("receive", 50002)
		   				//time.Sleep(1 * time.Second)
		   				//addreq is a channel on which ADD requests from peer are received
		   				ips := utilities.MyIpAddress()
		   				if len(ips) <= 0 {
		   					fmt.Printf("Error getting ip\n")
		   				}
		   				select {
		   				case addR := <-hbreq_data_control.DataC: //HB request received
		   					hbreq_data_control.ControlC <- "stop receiving ADD and exit"
		   					fmt.Printf("HB-request-receiving-channel closed:%v\n", <-hbreq_data_control.ControlC)
		   					hbres_data_control := communication.GetComm2()("send", 50003) //Send ACK 3 UDP packets for ensuring receipt
		   					for i := 0; i < 3; i++ {
		   						hbres_data_control.DataC <- utilities.Packet{
		   							FromIp: ips[0],
		   							ToIp:   addR.FromIp,
		   							Seq:    rand.Int63(),
		   							Req:    2,
		   						}
		   					}
		   					hbres_data_control.ControlC <- "Stop sending ACK"
		   					fmt.Printf("So done stopping:%v\n", <-hbres_data_control.ControlC)
		   					fmt.Printf("ACK sent 3 times to:%v\n", addR.FromIp)
		   					ChanSm <- addR
		   					//fmt.Printf("ADD request forwarded to Membership service:%v\n", addR.FromIp)
		   					continue HB_Loop
		   				case <-time.After(10 * time.Millisecond):
		   					continue HB_Loop
		   				}
		   				//Go for next ADD request
		   			}
		   		}() */

		fmt.Fprintf(os.Stdout, "Add died:%v\n", <-add_raised_panic)

	case fsm.State == 2:
		/*State 2 is a state for non-leader node
		send ADD request to Leader and Wait for Ack*/
		//Keep sending "ADD" request to leader
		added_as_member := make(chan string)
		go func() {
			for {
				stop_sending_add_req := SendAddReqToLeader()
				channel := communication.GetComm2()("receive", 50001)
				select {
				case ack := <-channel.DataC:
					if ack.Req == 2 {
						fmt.Printf("ack received:%v\n", ack)
						stop_sending_add_req <- "Done sending ADD"
						fmt.Printf("Done?:%v\n", <-stop_sending_add_req)
						added_as_member <- "added"
						fmt.Printf("request for command to stop:%v\n", <-added_as_member)
						return
					}
				}
			}
		}()
		fmt.Printf("Adder said:%v\n", <-added_as_member)
		added_as_member <- "quit now"

		/*Regular heartbeat begin here*/
		fmt.Printf("Regular heartbeats begin\n")
		for {
			ips := utilities.MyIpAddress()
			if len(ips) <= 0 {
				return errors.New("Can't get IP address\n")
			}

			go func() {
				for {
					//Create channel to send HB for 10 Milliseconds
					channelS := communication.GetComm2()("send", 50002)
					select {
					case channelS.DataC <- utilities.Packet{
						FromIp: ips[0],
						ToIp:   net.ParseIP("172.16.238.4"),
						Seq:    rand.Int63(),
						Req:    3,
					}:
						fmt.Printf("HB Sent, now sleep for 100 milliseconds\n")
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()

			go func() {
				for {
					channel := communication.GetComm2()("receive", 50003)
					select {
					case hbR := <-channel.DataC:
						fmt.Printf("Received %v\n", hbR)
						time.Sleep(10 * time.Millisecond)
					}
				}
			}()
		}

	}
	return errors.New("Shouldn't have returned from here\n")
}
