package fsm

import (
	"app/membership/communication"
	"app/membership/utilities"
	"fmt"
	"math/rand"
	"net"
)

func SendAddReqToLeader() chan string {
	stop_sending_add_req := make(chan string)
	go func() {
		data_and_control := communication.GetComm2()("send", 50000)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in SendAddReqToLeader", r)
				return
			}
		}()

		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error getting ip\n")
			return
		}

		i := 0
		for {
			i += 1
			select {
			case data_and_control.DataC <- utilities.Packet{
				FromIp: ips[0],
				ToIp:   net.ParseIP("172.16.238.2"),
				Seq:    rand.Int63(),
				Req:    1,
			}:

			case <-stop_sending_add_req:
				data_and_control.ControlC <- "Stop sending ADD"
				fmt.Printf("Have you stopped:%v\n", <-data_and_control.ControlC)
				stop_sending_add_req <- "Done stopping!"
				return
			}
		}
	}()
	return stop_sending_add_req
}
