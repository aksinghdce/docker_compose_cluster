package communication

import (
	"app/log"
	"app/membership/utilities"
	"context"
	"encoding/json"
	"net"
	"fmt"
	"time"
	"io"
)

/*
*/
func Comm(ctx context.Context, receiveport, sendport  int) (chan utilities.Packet, chan utilities.Packet) {
	listen := make(chan utilities.Packet)
	speak := make(chan utilities.Packet)
	//go routine that will listen for incoming datagrams and return channel as first
	//item in the output
	go func() {
		rerun := false
		for {
			time.Sleep(10 * time.Millisecond)
			if rerun == false {
				rerun = listener(ctx, listen, receiveport)
			}
		}
	}()
	//go routine that will speak out to the world at large, whatever it receives 
	//on the second output channel
	go func() {
		rerun := false
		for {
			time.Sleep(10 * time.Millisecond)
			if rerun == false {
				rerun = speaker(ctx, speak, sendport)
			}
		}
	}()
	return listen, speak
}


func Close(c io.Closer) bool{
	err := c.Close()
	if err != nil {
		fmt.Printf("Error in Comm Close:%v\n", err)
	}
	return false
}

func listener(ctx context.Context, listenChannel chan utilities.Packet, port int) bool{
	myaddr := &net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", myaddr)
	if err != nil {
		fmt.Printf("Error in Comm ListenUDP:%v\n", err)
		log.Log(ctx, err.Error())
		return false
	}
	conn.SetReadBuffer(1048576)

	defer func() {
		Close(conn)
	}()
	

	for {
		
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error in Comm ReadFromUDP:%v\n", err)
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
	return true		
}

func speaker(ctx context.Context, dialChannel chan utilities.Packet, port int) bool{
	sendThisPacket := <-dialChannel
	//Might require an error check in the next instruction
	toAddress := &net.UDPAddr{IP: sendThisPacket.ToIp, Port: port}
	Conn, err := net.DialUDP("udp", nil, toAddress)
	if err != nil {
		fmt.Printf("Error DialUDP:%s\n", err.Error())
		return false
	}

	defer func() {
		Close(Conn)
	}()
	
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
				return false
			}
			// Be ready for next iteration
			sendThisPacket = <-dialChannel
		}
	}
	return true
}