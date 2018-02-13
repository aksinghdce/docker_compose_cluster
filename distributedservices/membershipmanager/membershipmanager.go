package membershipmanager

/*
State transition scheme:

State 0:
		New Node
		data : Constants
			Saved Context : [A Binary Tree]
		Events:
			1. Begin clean slate: Join group to leader
				Source : Default
			Actions:
				1. Read configuration:
				Get leader's hostname : Move to state 1 or state 2
State 1:
		Leader Node
		data : Constants
			Saved Context : A Binary Tree
		Events:
			1. It receives a join request
			2. It receives a leave
			3. It receives a heartbeatloss
State 2:
		Peer Node
		data : Constants
		Events:
			1. default : receive and send heartbeats and report loss
			2. It receives a leave
			3. It receives a heartbeatloss
			4. Leader sent updated "new node added / "
*/

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
