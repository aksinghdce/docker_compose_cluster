package membershipmanager

type Event interface {
	getSource() string
	getStimulus() string
	getArtifact() string
	getEvent() string
}

type AddEvent struct {
	hostname  string
	ipAddress string
	timeStamp string
}

func (addEvent *AddEvent) getSource() string {
	return "source 1"
}

func (addEvent *AddEvent) getStimulus() string {
	return "stimulus 1"
}

func (addEvent *AddEvent) getArtifact() string {
	return "artifact 1"
}

func (addEvent *AddEvent) getEvent() string {
	return "event 1"
}

/*
membershipmanager package manages a statemachine
The statemachine keeps the distributed cluster state
*/
type MembershipManager interface {
	GetGroupInfo() []string
}

type ExtendedRingManager struct {
	groupInfo []string
}

func (erm *ExtendedRingManager) GetGroupInfo() []string {
	return erm.groupInfo
}

/*
Identity is an interface that prvides the following methods
GetHostname
GetIpAddress
IsLeader
CurrentTimeStamp
*/
func GetMembershipManager(event Event) ExtendedRingManager {
	erm := ExtendedRingManager{[]string{"amit", "kumar", "singh"}}
	return erm
}
