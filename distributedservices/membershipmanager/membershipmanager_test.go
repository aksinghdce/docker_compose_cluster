package membershipmanager

import "testing"

func TestState0(t *testing.T) {
	state := State{
		currentState: 0,
		leaderIp:     "124.0.0.1",
		leaderPort:   10001,
		managedNodes: []string{},
		amITheLeader: false,
		clusterMap:   nil,
	}

	internaleventforstate0 := InternalEvent{}

	erm := MembershipTreeManager{
		myState:   state,
		myLeader:  "124.0.0.1:10001",
		groupInfo: []string{},
	}

	status := erm.ProcessInternalEvent(internaleventforstate0)
	if status != "" {
		t.Logf("Got something on channel:%s", status)
	}

	if status == "state2" {
		t.Logf("Moving to state 2")
	}
}
