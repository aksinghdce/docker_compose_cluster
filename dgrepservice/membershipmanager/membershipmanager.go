package membershipmanager

type Event interface {
	getSource() string
	getStimulus() string
	getArtifact() string
	getEvent() string
}

type AddNodeEvent struct {
	hostname  string
	ipAddress string
	timeStamp string
}

func (ane *AddNodeEvent) getSource() string {
	return ane.hostname
}

func (ane *AddNodeEvent) getStimulus() string {
	return ane.ipAddress
}

func (ane *AddNodeEvent) getArtifact() string {
	return ane.timeStamp
}

func (AddNodeEvent *AddNodeEvent) getEvent() string {
	return "event 1"
}

/*
membershipmanager package manages a statemachine
The statemachine keeps the distributed cluster state
*/
type MembershipManager interface {
	GetGroupInfo() []string
	AddNodeToGroup() (error, string)
	RemoveNodeFromGroup() (error, string)
}

type ExtendedRingManager struct {
	groupInfo []string
}

func (erm *ExtendedRingManager) GetGroupInfo() []string {
	return erm.groupInfo
}

func (erm *ExtendedRingManager) AddNodeToGroup() (error, string) {
	return nil, "success"
}

func (erm *ExtendedRingManager) RemoveNodeFromGroup() (error, string) {
	return nil, "success"
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
