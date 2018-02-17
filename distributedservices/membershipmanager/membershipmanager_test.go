package membershipmanager

import (
	"testing"
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
	state := State{
		currentState: 0,
		leaderIp:     "124.0.0.1",
		leaderPort:   10001,
		managedNodes: []string{},
		amITheLeader: false,
		clusterMap:   nil,
	}

	mmm := NewMembershipManager(state)
	if mmm.myState.currentState == 1 || mmm.myState.currentState == 2 {
		t.Logf("successfully changed state")
	} else {
		t.Fatalf("State not changed")
	}
}

func TestState1(t *testing.T) {
	state := State{
		currentState: 1,
		leaderIp:     "124.0.0.1",
		leaderPort:   10001,
		managedNodes: []string{},
		amITheLeader: false,
		clusterMap:   nil,
	}

	mmm := NewMembershipManager(state)

	internaleventforstate0 := InternalEvent{}

	/*Processing the State0 default event must take the state to 1 or 2*/
	/*The following is an infinite loop*/
	mmm.ProcessInternalEvent(internaleventforstate0)
	/*Test whether the groupInfo grew in size
	 */
}
