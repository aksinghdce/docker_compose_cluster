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

	/*Processing the State0 default event must take the state to 1 or 2*/
	erm.ProcessInternalEvent(internaleventforstate0)
	if erm.myState.currentState == 1 || erm.myState.currentState == 2 {
		t.Logf("Correctly transitioning from State 0 to 1 or 2")
	}
}
