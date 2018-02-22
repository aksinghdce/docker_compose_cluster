package multicastheartbeater

import (
	"app/utilities"
	"context"
	"encoding/json"
	"net"
	"fmt"
)

/*
Specification:
Output: This function will return a channel to the caller. The caller can write
objects of type utilities.Heartbeat on this channel.false

Input: The function takes the address and the port it will send the data TO

TO-DO: Write unit test for it.
*/
func SendHeartBeatMessages(ctx context.Context, toAddress, toPort string) chan utilities.HeartBeat {
	heartbeatChannelIn := make(chan utilities.HeartBeat)
	go func() {
		toAddress += ":"
		toAddress += toPort
		toAddr, err := net.ResolveUDPAddr("udp", toAddress)
		if err != nil {
			utilities.Log(ctx, err.Error())
		}
		Conn, err := net.DialUDP("udp", nil, toAddr)
		if err != nil {
			//utilities.Log(ctx, err.Error())
			fmt.Printf("Error DialUDP:%s\n", err.Error())
			return
		}
		defer Conn.Close()
		for {
			select {
			case hb := <-heartbeatChannelIn:
				hb.FromTo.ToIp = toAddress
				jsonData, err := json.Marshal(hb)
				_, err = Conn.Write(jsonData)
				if err != nil {
					//utilities.Log(ctx, err.Error())
					fmt.Printf("Conn.Write:%s\n", err.Error())
					return
				}
			}
		}
	}()
	return heartbeatChannelIn
}
