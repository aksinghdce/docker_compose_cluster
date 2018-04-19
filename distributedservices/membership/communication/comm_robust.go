package communication

import (
	"app/membership/utilities"
	"context"
	"encoding/json"
	"fmt"
	"net"
)

type Channels2 struct {
	DataC    chan utilities.Packet
	ControlC chan string
}

func GetComm2() func(string, int) Channels2 {
	f := func(sendorreceive string, port int) Channels2 {
		var dataAndControl Channels2
		if sendorreceive == "send" {
			ctx := context.Background()
			dataAndControl = commSend2(ctx, port)
		} else if sendorreceive == "receive" {
			ctx := context.Background()
			dataAndControl = commReceive2(ctx, port)
		}
		return dataAndControl
	}
	return f
}

/*
 */
func commSend2(ctx context.Context, sendport int) Channels2 {

	dataAndControl := Channels2{
		DataC:    make(chan utilities.Packet),
		ControlC: make(chan string),
	}

	go func() {
		data0 := <-dataAndControl.DataC
		fmt.Printf("Received data to Send:%v\n", data0)
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error: Local Ip\n")
		}
		fromAddress := &net.UDPAddr{IP: ips[0]}
		toAddress := &net.UDPAddr{IP: data0.ToIp, Port: sendport}
		Conn, err := net.DialUDP("udp", fromAddress, toAddress)
		if err != nil {
			fmt.Printf("Error dialing up\n")
		}

		defer func() {
			fmt.Printf("commSend2:CLOSING CONNECTION\n")
			close(dataAndControl.DataC)
			Conn.Close()
		}()

	sender_loop:
		for {
			select {
			case data := <-dataAndControl.DataC:
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Printf("Error:%v\n", err.Error())
					continue sender_loop
				}
				_, err = Conn.Write(jsonData)
				if err != nil {
					fmt.Printf("Conn.Write:%s\n", err.Error())
					break sender_loop
				}
			case <-dataAndControl.ControlC:
				fmt.Printf("commSend2:CLOSING CONNECTION\n")
				close(dataAndControl.DataC)
				Conn.Close()
				dataAndControl.ControlC <- "Yes Sir!"
			}
		}
		return
	}()

	return dataAndControl
}

/*
 */
func commReceive2(ctx context.Context, sendport int) Channels2 {
	dataAndControl := Channels2{
		DataC:    make(chan utilities.Packet),
		ControlC: make(chan string),
	}

	//Go routine to read from UDP
	peer_sent_this := make(chan utilities.Packet)
	close_connection := make(chan string)
	go func() {
		ips := utilities.MyIpAddress()
		if len(ips) <= 0 {
			fmt.Printf("Error: Local Ip\n")
		}
		myaddr := &net.UDPAddr{
			IP:   ips[0],
			Port: sendport,
		}
		conn, err := net.ListenUDP("udp", myaddr)
		if err != nil {
			fmt.Printf("Error in Comm ListenUDP:%v\n", err)
			return
		}
		defer func() {
			fmt.Printf("commReceive2:CLOSING CONNECTION\n")
			close(dataAndControl.DataC)
			conn.Close()
		}()

	receiver_loop:
		for {
			buf := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("Error in Comm ReadFromUDP:%v\n", err)
				break receiver_loop
			}
			buf = buf[:n]
			var Result utilities.Packet
			err = json.Unmarshal(buf, &Result)
			if err != nil {
				fmt.Printf("Unmarshall Error\n")
				break receiver_loop
			}

			select {
			case <-close_connection:
				close(dataAndControl.DataC)
				conn.Close()
				close_connection <- "Done closing!"
			case peer_sent_this <- Result:
				fmt.Printf("Data sent UP\n")
			}
		}
	}()

	go func() {
		for {
			select {
			case data := <-peer_sent_this:
				fmt.Printf("Data from peer:%v\n", data)
			case <-dataAndControl.ControlC:
				fmt.Printf("RECEIVED REQUEST TO CLOSE CONNECTION\n")
				close_connection <- "close connection immediately"
				fmt.Printf("closed okay?:%v\n", <-close_connection)
				dataAndControl.ControlC <- "Okay done!"
			}
		}
	}()

	return dataAndControl
}
