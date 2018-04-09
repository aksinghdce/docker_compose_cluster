package communication

import (
	"app/log"
	"app/membership/utilities"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

/*
 */
func CommSend(ctx context.Context, sendport int) (chan utilities.Packet, chan bool) {
	speak := make(chan utilities.Packet)
	stop_speaking := make(chan bool)
	//go routine that will speak out to the world at large, whatever it receives
	//on the second output channel
	go func() {
		rerun := true
		for {
			time.Sleep(10 * time.Millisecond)
			if rerun == true {
				rerun = speaker(ctx, speak, stop_speaking, sendport)
			}
		}
	}()

	return speak, stop_speaking
}

func CommReceive(ctx context.Context, receiveport int) (chan utilities.Packet, chan bool) {
	listen := make(chan utilities.Packet)
	stop := make(chan bool)
	//go routine that will listen for incoming datagrams and return channel as first
	//item in the output
	go func() {
		rerun := true
		for {
			time.Sleep(10 * time.Millisecond)
			if rerun == true {
				rerun = listener(ctx, listen, stop, receiveport)
			}
		}
	}()

	return listen, stop
}

func Close(c io.Closer) bool {
	err := c.Close()
	if err != nil {
		fmt.Printf("Error in Comm Close:%v\n", err)
	}
	return false
}

/*
Will listen on ipv4 ipaddress configured by docker
*/
func listener(ctx context.Context, listenChannel chan utilities.Packet, stop chan bool, port int) bool {
	ips := utilities.MyIpAddress()
	if len(ips) <= 0 {
		fmt.Printf("Error: Local Ip\n")
	}
	myaddr := &net.UDPAddr{
		IP:   ips[0],
		Port: port,
	}
	conn, err := net.ListenUDP("udp", myaddr)
	if err != nil {
		fmt.Printf("Error in Comm ListenUDP:%v\n", err)
		log.Log(ctx, err.Error())
		return true
	}
	conn.SetReadBuffer(1048576)

	defer func() {
		Close(conn)
	}()

	for {
		select {
		case _ = <-stop:
			Close(conn)
			return false
		default:
			buf := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("Error in Comm ReadFromUDP:%v\n", err)
				log.Log(ctx, err.Error())
				return true
			}
			buf = buf[:n]
			var Result utilities.Packet
			err = json.Unmarshal(buf, &Result)
			if err != nil {
				log.Log(ctx, err.Error())
				return true
			}
			listenChannel <- Result
		}
	}
}

func speaker(ctx context.Context, dialChannel chan utilities.Packet, stop_speaking chan bool, port int) bool {
	sendThisPacket := <-dialChannel
	//Might require an error check in the next instruction
	ips := utilities.MyIpAddress()
	if len(ips) <= 0 {
		fmt.Printf("Error: Local Ip\n")
	}
	fromAddress := &net.UDPAddr{IP: ips[0]}
	toAddress := &net.UDPAddr{IP: sendThisPacket.ToIp, Port: port}
	Conn, err := net.DialUDP("udp", fromAddress, toAddress)
	if err != nil {
		//fmt.Printf("Error DialUDP:%s\n", err.Error())
		return true
	}

	defer func() {
		Close(Conn)
	}()

	for {
		timeout := time.After(100 * time.Millisecond)
		select {
		case _ = <-stop_speaking:
			Close(Conn)
			return false
		case <-timeout:
			fmt.Printf("Nothing received from up the stack; retrying\n")
		default:
			jsonData, err := json.Marshal(sendThisPacket)
			if err != nil {
				fmt.Printf("Error:%v\n", err.Error())
			}
			_, err = Conn.Write(jsonData)
			if err != nil {
				fmt.Printf("Conn.Write:%s\n", err.Error())
				return true
			}
			// Be ready for next iteration
			sendThisPacket = <-dialChannel
		}
	}
}
