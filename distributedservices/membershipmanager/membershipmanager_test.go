package membershipmanager

import "testing"

func TestState0(t *testing.T) {
	state := State{
		currentState: 0,
		leaderIp:     "124.0.0.1",
		leaderPort:   10001,
		managedNodes: []string{},
		amITheLeader: false,
	}

	internaleventforstate0 := InternalEvent{
		state:       0,
		stateObject: state,
	}

	erm := MembershipTreeManager{
		myLeader:  "124.0.0.1:10001",
		groupInfo: []string{},
	}

	status := erm.ProcessInternalEvent(internaleventforstate0)
	if status == 0 {
		t.Logf("Successful transition of state")
	}
}
