package communication

import (
	"app/membership/utilities"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

var udpServerTests = []struct {
	snet, saddr string // server endpoint
	tnet, taddr string // target endpoint for client
}{
	/* {snet: "udp", saddr: "0.0.0.0", tnet: "udp", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: ":0", tnet: "udp", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "0.0.0.0:0", tnet: "udp", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "[::]:0", tnet: "udp", taddr: "::1"},

	   {snet: "udp", saddr: ":0", tnet: "udp", taddr: "::1"},
	   {snet: "udp", saddr: "0.0.0.0:0", tnet: "udp", taddr: "::1"},
	   {snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp", taddr: "::1"},
	   {snet: "udp", saddr: "[::]:0", tnet: "udp", taddr: "127.0.0.1"},

	   {snet: "udp", saddr: ":0", tnet: "udp4", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "0.0.0.0:0", tnet: "udp4", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp4", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "[::]:0", tnet: "udp6", taddr: "::1"},

	   {snet: "udp", saddr: ":0", tnet: "udp6", taddr: "::1"},
	   {snet: "udp", saddr: "0.0.0.0:0", tnet: "udp6", taddr: "::1"},
	   {snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp6", taddr: "::1"},
	   {snet: "udp", saddr: "[::]:0", tnet: "udp4", taddr: "127.0.0.1"},

	   {snet: "udp", saddr: "127.0.0.1:0", tnet: "udp", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "[::ffff:127.0.0.1]:0", tnet: "udp", taddr: "127.0.0.1"},
	   {snet: "udp", saddr: "[::1]:0", tnet: "udp", taddr: "::1"},

	   {snet: "udp4", saddr: ":0", tnet: "udp4", taddr: "127.0.0.1"},
	   {snet: "udp4", saddr: "0.0.0.0:0", tnet: "udp4", taddr: "127.0.0.1"},
	   {snet: "udp4", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp4", taddr: "127.0.0.1"},

	   {snet: "udp4", saddr: "127.0.0.1:0", tnet: "udp4", taddr: "127.0.0.1"},

	   {snet: "udp6", saddr: ":0", tnet: "udp6", taddr: "::1"},
	   {snet: "udp6", saddr: "[::]:0", tnet: "udp6", taddr: "::1"},
	   {snet: "udp6", saddr: "[::1]:0", tnet: "udp6", taddr: "::1"}, */
	{snet: "udp6", saddr: "172.16.238.2", tnet: "udp6", taddr: "::1"},
	/* {snet: "udp6", saddr: "172.16.238.3", tnet: "udp6", taddr: "::1"},
	{snet: "udp6", saddr: "172.16.238.4", tnet: "udp6", taddr: "::1"},
	{snet: "udp6", saddr: "172.16.238.5", tnet: "udp6", taddr: "::1"}, */
}

func TestUDPServer(t *testing.T) {
	Ips := utilities.MyIpAddress()
	if len(Ips) <= 0 {
		t.Fatal("Couldn't get IP\n")
	}
	fmt.Printf("Ip:%v\n", Ips[0])

	listenChannel := GetComm2()("receive", 50000)
	speakChannel := GetComm2()("send", 50000)
	for caseNumber, tt := range udpServerTests {
		fmt.Printf("Running Test Number:%d\n", caseNumber)
		packet := utilities.Packet{
			FromIp: Ips[0],
			ToIp:   net.ParseIP(tt.saddr),
			Seq:    rand.Int63(),
		}

		for {
			timeout := time.After(2 * time.Second)
			select {
			//Sending/Receiving test packet
			case received_something:= <-listenChannel:
				
					fmt.Printf("-%v-", received_something)
				
			case <-timeout:
				packet.Seq = rand.Int63()
				speakChannel <- packet
				
			}
		}
	}
}
