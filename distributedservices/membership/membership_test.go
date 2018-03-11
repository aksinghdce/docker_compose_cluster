package membership

import (
	"testing"
	"app/membership/utilities"
)

func TestMembership(t *testing.T) {
	m := Membership{
		chanIn : make(chan utilities.Packet),
		chanOut : make(chan utilities.Packet),
	}
	cin, cout := m.KeepMembershipUpdated()	
	for i:=0; i<10; i++ {
		t.Logf("Received in test:%v", <-cin)
		cout <- utilities.Packet{}
	}
	t.Log("Ran successfully")
}