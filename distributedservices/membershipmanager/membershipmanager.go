package membershipmanager

import (
	"fmt"

	"app/multicastheartbeatserver"
)

/*
State transition scheme:

State 0:
		New Node
		data :
		Constants:
			1. Leader's UDP address and port
			Saved Context : [A Binary Tree]
		Events:
			1. Begin clean slate
			Actions:
				Get leader's hostname and Move to state 1 or state 2
State 1:
		Leader Node
		data : List of machines that are sending me heartbeats
		Internal protocol:
			Internal Events:
				1. Number of machines became > 3
				Action : Appoint the first machine in the list as the sub-tier leader
						1. Ask new machines to heartbeat to sub-tier leader
						2. Update local data structure to reflect which machine is the sub-tier
						leader and what is it's assignment list.
				2. Received request to provide membership list
				Action : Send membership list to requesting machine or local user
				3. One of my appointed sub-leader died
				Action : Update internal data structure and
					Send updates to other appointed subleaders
		Events:
			1. It receives a join request
				Action : Run Internal protocol : Internal event 1
			2. It receives a leave request
				Action : Send update to sub-leaders
			3. It receives a heartbeatloss
				Action : Internal protocol -> Internal event 3
State 2:
		Peer Node
		Internal Protocol:
			Internal Events:
				1. Leader asked me to be a sub-leader
				Action : Listen on the port number provided by the leader for heartbeats
						from subscription list
				2. Someone died in my subscription list
				Action : Update internal data structure and Report to my leader
		data :
		Events:
			1. default : receive and send heartbeats and report loss
			2. It receives a leave
			3. It receives a heartbeatloss
			4. Leader sent updated "new node added / "
*/

type State struct {
	currentState int8
	/*State 0 data*/
	leaderIp     string
	leaderPort   int
	managedNodes []string
	amITheLeader bool
}

type Event interface {
	getSource() string
	getStimulus() string
	getArtifact() string
	getEvent() string
}

type InternalEvent struct {
	state       int8
	stateObject State
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
	ProcessInternalEvent(intevent InternalEvent) string
	GetGroupInfo() []string
	AddNodeToGroup() (error, string)
	RemoveNodeFromGroup() (error, string)
}

type MembershipTreeManager struct {
	myLeader  string
	groupInfo []string
}

func (erm *MembershipTreeManager) ProcessInternalEvent(intevent InternalEvent) string {
	fmt.Println("Internal event:", intevent)
	if intevent.state == 0 {
		fmt.Println("My state is:", intevent.state)
		udps := multicastheartbeatserver.UdpServer{}

		ch := udps.ListenAndReport()
		output := <-ch
		fmt.Printf("Channel reads:%s", output)
		return output
	}
	return ""
}

func (erm *MembershipTreeManager) GetGroupInfo() []string {
	return erm.groupInfo
}

func (erm *MembershipTreeManager) AddNodeToGroup() (error, string) {
	return nil, "success"
}

func (erm *MembershipTreeManager) RemoveNodeFromGroup() (error, string) {
	return nil, "success"
}

/*
Identity is an interface that prvides the following methods
GetHostname
GetIpAddress
IsLeader
CurrentTimeStamp
*/
func GetMembershipManager(event Event) MembershipTreeManager {
	erm := MembershipTreeManager{
		myLeader:  "124.0.0.1:10001",
		groupInfo: []string{"amit", "kumar", "singh"},
	}
	return erm
}
