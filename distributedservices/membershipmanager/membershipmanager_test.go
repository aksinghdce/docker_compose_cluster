package membershipmanager

import (
	"fmt"
	"os"
	"testing"
	"time"
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
	}
	mmm := GetInstance()
	internaleventforstate1 := InternalEvent{}
	go func() {
		mmm.ProcessInternalEvent(internaleventforstate1)
	}()
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		if len(mmm.MyState.ClusterMap) == 0 {
			t.Fatalf("Map size 0, even after running")
		}
	}
	// else {
	// 	state := State{
	// 		CurrentState: 1,
	// 		LeaderIp:     "124.0.0.1",
	// 		LeaderPort:   10001,
	// 		ManagedNodes: []string{},
	// 		AmITheLeader: false,
	// 		ClusterMap:   nil,
	// 	}

	// 	mmm := NewMembershipManager(state)

	// 	internaleventforstate0 := InternalEvent{}

	// 	/*Processing the State0 default event must take the state to 1 or 2*/
	// 	/*The following is an infinite loop*/
	// 	mmm.ProcessInternalEvent(internaleventforstate0)
	// 	/*Test whether the groupInfo grew in size
	// 	 */
	// }

}

func TestState2(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Errorf("Error Hostname Resolution: %s", err.Error())
	}

	fmt.Printf("My hostname:%s\n", hostname)

	// if hostname != "leader.assignment2" {

	// 	state := State{
	// 		CurrentState:   2,
	// 		LeaderIp:       "124.0.0.1",
	// 		LeaderPort:     10001,
	// 		ManagedNodes:   []string{},
	// 		AmITheLeader:   false,
	// 		ClusterMap:     nil,
	// 		RequestContext: nil,
	// 	}

	// 	mmm := NewMembershipManager(state)

	// 	internaleventforstate0 := InternalEvent{}

	// 	mmm.ProcessInternalEvent(internaleventforstate0)

	// }

}
