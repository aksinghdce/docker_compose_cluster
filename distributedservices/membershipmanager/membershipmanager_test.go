package membershipmanager

import (
	"fmt"
	"os"
	"testing"
	"sort"
)

/*
Tests whether a machine goes from state 0 to state 1 or 2
A machine knows its hostname and based on that it can decide
whether to become an introducer machine.

If the machine becomes an introducer machine, then this machine
can accept "add" requests and heartbeats. Only one machine can in
state 1 because there is only one introducer for this assignment
*/
func TestState0(t *testing.T) {
	mmm := GetInstance()
	if mmm.MyState.CurrentState == 1 || mmm.MyState.CurrentState == 2 {
		t.Logf("successfully changed state")
	} else {
		t.Fatalf("State not changed")
	}
}

func TestState1(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Errorf("Error Hostname Resolution: %s", err.Error())
	}
	if hostname != "leader.assignment2" {
		t.Logf("This Machine is not in State 1, hostname:%s", hostname)
		return
	}
	mmm := GetInstance()

	
	if len(mmm.MyState.ClusterMap) < 1 {
		t.Fatal("Size of map is 0\n")
	}
	
}

func TestState2(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Errorf("Error Hostname Resolution: %s", err.Error())
	}

	fmt.Printf("My hostname:%s\n", hostname)
}

func TestState3(t *testing.T) {
	mmm := GetInstance()
	
	sortedIp := mmm.SortCurrentGroupInfo()
	if sort.StringsAreSorted(sortedIp) {
		t.Log("Ips are sorted")
	}else {
		t.Fatal("Not Sorted")
	}
}