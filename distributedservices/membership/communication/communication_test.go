package communication

import (
	"app/membership/utilities"
	"context"
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
	{snet: "udp6", saddr: "172.16.238.3", tnet: "udp6", taddr: "::1"},
	{snet: "udp6", saddr: "172.16.238.4", tnet: "udp6", taddr: "::1"},
	{snet: "udp6", saddr: "172.16.238.5", tnet: "udp6", taddr: "::1"},
}

func TestUDPServer(t *testing.T) {
	Ips := utilities.MyIpAddress()
	if len(Ips) <= 0 {
		t.Fatal("Couldn't get IP\n")
	}
	fmt.Printf("Ip:%v\n", Ips[0])

	ctx := context.Background()
	listenChannel := CommReceive(ctx, 50000)
	speakChannel := CommSend(ctx, 50000)
	for caseNumber, tt := range udpServerTests {
		fmt.Printf("Running Test Number:%d\n", caseNumber)
		packet := utilities.Packet{
			FromIp: Ips[0],
			ToIp:   net.ParseIP(tt.saddr),
			Seq:    rand.Int63(),
		}
		var receivedPacket utilities.Packet
		for i := 0; i < 5; i++ {
			//Sending test packet
			speakChannel <- packet
			timeout := time.After(10 * time.Second)
			select {
			//Receiving test packet
			case receivedPacket = <-listenChannel:
				fmt.Printf("Packet:%v received\n", receivedPacket)
			case <-timeout:
				t.Fatal("Packet not received in 10 second")
			}
		}
	}
}
