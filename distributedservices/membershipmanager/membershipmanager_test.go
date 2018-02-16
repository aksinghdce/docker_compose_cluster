package membershipmanager

import (
	"testing"
)

func TestState0(t *testing.T) {
	state := State{
		currentState: 0,
		leaderIp:     "124.0.0.1",
		leaderPort:   10001,
		managedNodes: []string{},
		amITheLeader: false,
		clusterMap:   nil,
	}

	mmm := NewMembershipManager(state, "224.0.1.1")
	internaleventforstate0 := InternalEvent{}

	// erm := MembershipTreeManager{
	// 	myState:   state,
	// 	myLeader:  "124.0.0.1:10001",
	// 	groupInfo: []string{},
	// }

	/*Processing the State0 default event must take the state to 1 or 2*/
	mmm.ProcessInternalEvent(internaleventforstate0)
	if mmm.myState.currentState == 1 {
		t.Logf("successfully changed state")
	} else {
		t.Fatalf("State not changed")
	}

}
