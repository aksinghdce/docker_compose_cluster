package communication

import (
	"app/log"
	"app/membership/utilities"
	"context"
	"encoding/json"
	"net"
	"fmt"
	"time"
)

/*
*/
func Comm(ctx context.Context, receiveport, sendport  int) (chan utilities.Packet, chan utilities.Packet) {
	listen := make(chan utilities.Packet)
	speak := make(chan utilities.Packet)
	
	//go routine that will listen for incoming datagrams and return channel as first
	//item in the output
	go listener(ctx, listen, receiveport)

	//go routine that will speak out to the world at large, whatever it receives 
	//on the second output channel
	go speaker(ctx, speak, sendport)
	return listen, speak
}

func listener(ctx context.Context, listenChannel chan utilities.Packet, port int) {
	myaddr := &net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", myaddr)
	if err != nil {
		log.Log(ctx, err.Error())
	}
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Log(ctx, err.Error())
			continue
		}
		buf = buf[:n]
		var Result utilities.Packet
		err = json.Unmarshal(buf, &Result)
		if err != nil {
			log.Log(ctx, err.Error())
			continue
		}
		listenChannel <- Result
	}		
}

func speaker(ctx context.Context, dialChannel chan utilities.Packet, port int) {
	sendThisPacket := <-dialChannel
	//Might require an error check in the next instruction
	toAddress := &net.UDPAddr{IP: sendThisPacket.ToIp, Port: port}
	Conn, err := net.DialUDP("udp", nil, toAddress)
	if err != nil {
		fmt.Printf("Error DialUDP:%s\n", err.Error())
		return
	}
	defer Conn.Close()
	for {
		timeout := time.After(100 * time.Millisecond)
		select {
		case <-timeout : 
			fmt.Printf("Nothing received from up the stack; retrying\n")
		default:
			jsonData, err := json.Marshal(sendThisPacket)
			if err != nil {
				fmt.Printf("Error:%v\n", err.Error())
			}
			_, err = Conn.Write(jsonData)
			if err != nil {
				fmt.Printf("Conn.Write:%s\n", err.Error())
				return
			}
			// Be ready for next iteration
			sendThisPacket = <-dialChannel
		}
	}
}