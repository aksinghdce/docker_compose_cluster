package communication

import (
	"app/membership/utilities"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Channels2 struct {
	DataC    chan utilities.Packet
	ControlC chan string
}

/*
Returns a function
	 which returns a Channel2 object
Uses:
	1. Caller can parameterize the returned function for a "send" or "receive"
	communication. And the port for udp communication only.
*/
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
			if r := recover(); r != nil {
				fmt.Printf("Recovered in ProcessFsm STATE 1:%v !!", r)
				fmt.Printf("commSend2:CLOSING CONNECTION\n")
				close(dataAndControl.DataC)
				err := Conn.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error closing connection\n")
				}
				dataAndControl.ControlC <- "Yes Sir!"
			}
		}()
	send_loop:
		for {
			select {
			case data, ok := <-dataAndControl.DataC:
				if !ok {
					fmt.Printf("Channel is closed\n")
					return
				}
				if data.Req == 2 {
					fmt.Printf("Sending back ACK%v\n", data)
				}

				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Printf("Error:%v\n", err.Error())
					panic("Can't marshall data\n")
				}
				_, err = Conn.Write(jsonData)
				if err != nil {
					fmt.Printf("Conn.Write:%s\n", err.Error())
					//panic("Can't write on Connection\n")
					break send_loop
				}
			case <-dataAndControl.ControlC:
				panic("closing everything\n")
			}
		}
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
			//fmt.Printf("Error in Comm ListenUDP:%v\n", err)
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
				err := conn.Close()
				if err != nil {
					fmt.Printf("There was an error closing the connection\n")
				}
				close_connection <- "Done closing!"
				return
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
				dataAndControl.DataC <- data
			case <-dataAndControl.ControlC:
				fmt.Printf("RECEIVED REQUEST TO CLOSE CONNECTION\n")
				close_connection <- "close connection immediately"
				fmt.Printf("closed okay?:%v\n", <-close_connection)
				dataAndControl.ControlC <- "Okay done!"
				return
			}
		}
	}()

	return dataAndControl
}
