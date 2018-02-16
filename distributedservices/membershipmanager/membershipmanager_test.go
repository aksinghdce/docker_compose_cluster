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
	if mmm.myState.currentState == 1 {
		t.Logf("successfully changed state")
	} else {
		t.Fatalf("State not changed")
	}

	internaleventforstate0 := InternalEvent{}

	/*Processing the State0 default event must take the state to 1 or 2*/
	mmm.ProcessInternalEvent(internaleventforstate0)
	/*Test whether the groupInfo grew in size
	 */

}
