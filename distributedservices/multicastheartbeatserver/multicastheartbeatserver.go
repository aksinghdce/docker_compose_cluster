package multicastheartbeatserver

import (
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

/*
Specification
Use this function only in FSM State 1
Output: I return the channel on which you can read what I read on multicast port
Input: The caller context, the multicast ip address, the udp port to listen to ADD request
*/
func CatchMultiCastDatagramsAndBounce(ctx context.Context, iListenOnIp, iListenOnPort string) chan utilities.HeartBeatUpperStack {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	c := make(chan utilities.HeartBeatUpperStack)
	port, errconv := strconv.Atoi(iListenOnPort)
	if errconv != nil {
		utilities.Log(ctx, errconv.Error())
	}

	conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(iListenOnIp),
		Port: port,
	})
	if err != nil {
		utilities.Log(ctx, err.Error())
	}
	go func() {
		for {
			defer conn.Close()

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
			var hbu utilities.HeartBeatUpperStack
			hbu.Hb = Result
			hbu.Ip = udpAddr.IP.String()
			fmt.Printf("Received Lower stack:%v\n", hbu)
			c <- hbu
		}
	}()
	return c
}

/*
Specification:
Output: I return the channel on which you can listen to the acknowledgement given to ADD request
Input: The udp port to listen on.
*/
func CatchUniCastDatagramsAndBounce(ctx context.Context, iListenOnPort string) chan utilities.HeartBeatUpperStack {
	/*I am the leader I will keep listening to events
	from my group. And keep writing the response on the channel chout*/
	c := make(chan utilities.HeartBeatUpperStack)
	//addrstr := string("224.0.0.1")
	port, errconv := strconv.Atoi(iListenOnPort)
	if errconv != nil {
		utilities.Log(ctx, errconv.Error())
	}
	myaddr := &net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", myaddr)
	if err != nil {
		utilities.Log(ctx, err.Error())
	}

	go func() {
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
			var hbu utilities.HeartBeatUpperStack
			hbu.Hb = Result
			hbu.Ip = udpAddr.IP.String()
			c <- hbu
		}
	}()
	return c
}
