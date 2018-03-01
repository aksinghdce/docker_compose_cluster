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
	dial        bool   // test with Dial
}{
	{snet: "udp", saddr: "0.0.0.0", tnet: "udp", taddr: "127.0.0.1"},
}

func TestUDPServer(t *testing.T) {
	ctx := context.Background()
	
	for _, tt := range udpServerTests {
		packet := Packet{
			FromIp: net.ParseIP(tt.taddr),
			ToIp: net.ParseIP(tt.saddr),
			Seq: rand.Int63(),
		}
		var receivedPacket Packet
		//get communication channels and test the channels
		listenChannel, speakChannel := Comm(ctx, 50000, 50000)
		for i := 0; i<5; i++ {
			speakChannel <- packet
			timeout := time.After(10 * time.Second)
			select {
			case receivedPacket = <-listenChannel:
				fmt.Printf("Packet:%v received\n",receivedPacket)
			case <-timeout:
				t.Fatal("Packet not received in 10 second")
			}
		}
	}
}