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
		i := 0
		data0 := <-dataAndControl.DataC
		ok := true

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
			if !ok {
				fmt.Printf("Channel is closed\n")
				return
			}
			if data0.Req == 2 {
				fmt.Printf("Sending back ACK%v\n", data0)
			}

			jsonData, err := json.Marshal(data0)
			if err != nil {
				fmt.Printf("Error:%v\n", err.Error())
				panic("Can't marshall data\n")
			}

			i += 1
			fmt.Fprintf(os.Stdout, "Data:%d sent\n", i)
			_, err = Conn.Write(jsonData)
			if err != nil {
				fmt.Printf("Conn.Write:%s\n", err.Error())
				//panic("Can't write on Connection\n")
				break send_loop
			}

			select {
			case data, ok := <-dataAndControl.DataC:
				if !ok {
					fmt.Printf("Channel is closed\n")
					return
				} else {
					data0 = data
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

	go func() {
		peer_sent_this := make(chan utilities.Packet)
		close_connection := make(chan string)

		go func() {
			for {
				select {
				case data := <-peer_sent_this:
					fmt.Printf("Data from peer:%v\n", data)
					dataAndControl.DataC <- data
				case <-close_connection:
					fmt.Fprintf(os.Stdout, "RECEIVED REQUEST TO CLOSE CONNECTION\n")
					close_connection <- "closing connection immediately"
					return
				}
			}
		}()

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
			case <-dataAndControl.ControlC:
				close_connection <- "Stop sending data above"
				fmt.Fprintf(os.Stdout, "stopped sending data:%v\n", <-close_connection)
				err := conn.Close()
				if err != nil {
					fmt.Printf("There was an error closing the connection\n")
				}
				close(dataAndControl.DataC)
				return
			case peer_sent_this <- Result:
				fmt.Printf("Data sent UP\n")
			}
		}
	}()

	return dataAndControl
}
