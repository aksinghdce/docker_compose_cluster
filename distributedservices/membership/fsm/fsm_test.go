package fsm

import (
    "testing"
    "net"
    "app/membership/utilities"
    "context"
    "math/rand"
	"fmt"
    "app/membership/communication"
)

var udpServerTests = []struct {
	snet, saddr string // server endpoint
	tnet, taddr string // target endpoint for client
}{
	{snet: "udp", saddr: "127.0.0.1", tnet: "udp", taddr: "127.0.0.1"},
}


func TestFsm(t *testing.T) {
    done := make(chan bool)
    
    fsm1 := Init(1)
    fsm2 := Init(2)
    fsm1.ProcessFsm()
    err, newState := fsm2.ProcessFsm()
    if err == nil {
        fsm2 = Init(newState)
        fsm2.ProcessFsm()
    }
    //Run communication test in parallel to do some stress testing
    //Because 
    go func() {
        ctx := context.Background()
	    listenChannel, speakChannel := communication.Comm(ctx, 50000, 50000)
	    for _, tt := range udpServerTests {
		packet := utilities.Packet{
			FromIp: net.ParseIP(tt.taddr),
			ToIp: net.ParseIP(tt.saddr),
			Seq: rand.Int63(),
		}
		var receivedPacket utilities.Packet
		    for i := 0; i<5; i++ {
			    //Sending or Receive test packet
			    select {
				//Receiving test packet
			    case receivedPacket = <-listenChannel:
                    fmt.Printf("Packet:%v received\n",receivedPacket)
                //Send by default
                default:
				    speakChannel <- packet
			    }
		    }
        }
     done <- true
    }()
    
    <-done
    
    t.Log("Well done!")
}