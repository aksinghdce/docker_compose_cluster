package communication

import (
	"testing"
	"context"
	"net"
	"math/rand"
	"fmt"
	"time"
)

var udpServerTests = []struct {
	snet, saddr string // server endpoint
	tnet, taddr string // target endpoint for client
}{
	{snet: "udp", saddr: "0.0.0.0", tnet: "udp", taddr: "127.0.0.1"},
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
	{snet: "udp6", saddr: "[::1]:0", tnet: "udp6", taddr: "::1"},
}

func TestUDPServer(t *testing.T) {
	ctx := context.Background()
	listenChannel, speakChannel := Comm(ctx, 50000, 50000)
	for caseNumber, tt := range udpServerTests {
		fmt.Printf("Running Test Number:%d\n", caseNumber)
		packet := Packet{
			FromIp: net.ParseIP(tt.taddr),
			ToIp: net.ParseIP(tt.saddr),
			Seq: rand.Int63(),
		}
		var receivedPacket Packet
		for i := 0; i<5; i++ {
			//Sending test packet
			speakChannel <- packet
			timeout := time.After(10 * time.Second)
			select {
				//Receiving test packet
			case receivedPacket = <-listenChannel:
				fmt.Printf("Packet:%v received\n",receivedPacket)
			case <-timeout:
				t.Fatal("Packet not received in 10 second")
			}
		}
	}
}