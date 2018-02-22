package multicastheartbeatserver

import (
	"app/utilities"
	"context"
	"encoding/json"
	"net"
	"strconv"
)

/*
Specification
Use this function only in FSM State 1
Output: I return the channel on which you can read what I read on multicast port
Input: The caller context, the multicast ip address, the udp port to listen to ADD request
*/
func CatchMultiCastDatagramsAndBounce(ctx context.Context, iListenOnIp, iListenOnPort string) chan utilities.HeartBeat {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	c := make(chan utilities.HeartBeat)
	port, errconv := strconv.Atoi(iListenOnPort)
	if errconv != nil {
		utilities.Log(ctx, errconv.Error())
	}

	
	go func() {
		conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
			IP:   net.ParseIP(iListenOnIp),
			Port: port,
		})
		if err != nil {
			utilities.Log(ctx, err.Error())
		}
		defer conn.Close()
		for {
			buf := make([]byte, 100)
			n, udpAddr, err2 := conn.ReadFromUDP(buf)
			if err2 != nil {
				utilities.Log(ctx, err2.Error())
			}
			buf = buf[:n]
			var Result utilities.HeartBeat

			errUnmarshal := json.Unmarshal(buf, &Result)
			if errUnmarshal != nil {
				utilities.Log(ctx, errUnmarshal.Error())
			}
			Result.FromTo.FromIp = udpAddr.IP.String()
			c <- Result
		}
	}()
	return c
}

/*
Specification:
Output: I return the channel on which you can listen to the acknowledgement given to ADD request
Input: The udp port to listen on.
*/
func CatchUniCastDatagramsAndBounce(ctx context.Context, iListenOnPort string) chan utilities.HeartBeat {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	c := make(chan utilities.HeartBeat)
	//addrstr := string("224.0.0.1")
	port, errconv := strconv.Atoi(iListenOnPort)
	if errconv != nil {
		utilities.Log(ctx, errconv.Error())
	}
	myaddr := &net.UDPAddr{Port: port}
	

	go func() {
		conn, err := net.ListenUDP("udp", myaddr)
		if err != nil {
			utilities.Log(ctx, err.Error())
		}
		defer conn.Close()
		for {

			buf := make([]byte, 1024)
			n, udpAddr, err2 := conn.ReadFromUDP(buf)
			if err2 != nil {
				utilities.Log(ctx, err2.Error())
				continue
			}
			buf = buf[:n]
			var Result utilities.HeartBeat
			//fmt.Printf("Received:%v\n", string(buf))
			errUnmarshal := json.Unmarshal(buf, &Result)
			if errUnmarshal != nil {
				utilities.Log(ctx, errUnmarshal.Error())
				continue
			}
			Result.FromTo.FromIp = udpAddr.IP.String()
			c <- Result
		}
	}()
	return c
}
