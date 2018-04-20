package fsm

import (
	"testing"
	//"net"
	//"app/membership/utilities"
	//"context"
	//"math/rand"
	"fmt"
	//"app/membership/communication"
	//"app/membership/fsm"
	"os"
)

var udpServerTests = []struct {
	snet, saddr string // server endpoint
	tnet, taddr string // target endpoint for client
}{
	{snet: "udp", saddr: "127.0.0.1", tnet: "udp", taddr: "127.0.0.1"},
}

func TestFsm(t *testing.T) {
	//done := make(chan bool)
	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error:%s\n", err.Error())
	}
	fmt.Printf("HOSTNAME:%v\n", host)
	if host == "leader.assignment2" {
		fsm1 := Init(1)
		fsm1.ProcessFsm()
	} else {
		fsm2 := Init(2)
		err := fsm2.ProcessFsm()
		if err != nil {
			t.Fatalf("ProcessFsm exited with error:%v\n", err)
		}
	}

	t.Log("Well done!")
}
