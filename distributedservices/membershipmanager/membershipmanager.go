package membershipmanager

import (
	"app/utilities"
	"fmt"
	"math/rand"
	"os"
	"time"

	"app/multicastheartbeater"
	"app/multicastheartbeatserver"
)

/*
State transition scheme:

State 0:
		New Node
		data :
		Environment: hostname, ip and port for multicast
		Constants:
			1. Leader's UDP address and port
			Saved Context : [A Binary Tree]
		Events:
			1. Begin clean slate
			Actions:
				Get leader's hostname and Move to state 1 or state 2
State 1:
		Leader Node
		data :
			1. List of machines that are sending me "add" request
				I add it to the group and send the result to my three successors
			2. List of machines that are sending me heartbeats
				I check whether they are in my group list, update their timestamp
				If they are not in my group, then ignore the heartbeat
			3. Cluster map
				List in step 1, my predecessors and my successors.
			4. My leader : nil for the 124.0.0.1:10001 leader (introducer)
		Internal protocol:
			Internal Events:
				0. New machine sent ping
					Given : Number of machines currently I am heartbeating with < 3
						Add machine as sub-leader and open a heartbeat channel with it
					Given : Number of machines I am heartbeating with == 4
					Action : Appoint the first machine in the list as the sub-tier leader
						The first machine in the list will deligate to the second machine
						in the list if the number of nodes it is managing is already 4
						1. Ask new machines to heartbeat to sub-tier leader
							Send sub-leader info to new machine
						2. Update local data structure to reflect which machine is the sub-tier
				2. Received request to provide membership list
					Action : Send membership list to requesting machine or local user
				3. One of my appointed sub-leader died
				Action : Update internal data structure and wait for orphans to ping me
		Events:
			1. It receives a join request
				Action : Run Internal protocol : Internal event 1
			2. It receives a leave request
				Action : Send update to sub-leaders
			3. It receives a heartbeatloss
				Action : Internal protocol -> Internal event 3
State 2:
		Peer Node asking to join
		Internal Protocol:
			Internal Events:
				1. Leader asked me to be a sub-leader
				Action : Listen on the port number provided by the leader for heartbeats
						from subscription list
				2. Someone died in my subscription list
				Action : Update internal data structure and Report to my leader
				3. My leader died
				Action : Pick the next available sub-leader
					How?
					Answer: In the cluster map, pick the next machine that is supposed to
					be alive. Heartbeat with it to add you and wait for the response.
					If the heartbeat response is not received, then ping the next in the series.
		data :
			1. List of machines that are sending me heartbeats
			2. Cluster map
		Events:
			1. default : receive and send heartbeats and report loss
			2. It receives a leave
			3. It receives a heartbeatloss
			4. Leader sent updated "new node added / "
State 3:
		A group member doing it's usual heartbeats
		Report missed heartbeats to subscribers

*/

type State struct {
	currentState int8
	/*State 0 data*/
	leaderIp     string
	leaderPort   int
	managedNodes []string
	amITheLeader bool
	clusterMap   map[string]string
}

/*
Functions to change the state

1. Add new node to clusterMap
2. Remove a node from clusterMap
3. Update last heartbeat timestamp for a node in cluster map
4. Update my leaderIp and leaderPort
5. Change my state
6. Add to my managed nodes
7. Remove from my managed nodes
*/

type Event interface {
	getSource() string
	getStimulus() string
	getArtifact() string
	getEvent() string
}

type InternalEvent struct{}

/*
membershipmanager package manages a statemachine
The statemachine keeps the distributed cluster state
*/
type MembershipManager interface {
	ProcessInternalEvent(intevent InternalEvent)
	GetGroupInfo() []string
	AddNodeToGroup() (error, string)
	RemoveNodeFromGroup() (error, string)
}

type MembershipTreeManager struct {
	// What is my current state?
	myState   State
	groupInfo []string
}

/*
Specification:

Input: InternalEvent is for future use only.
Output:
Processing:
The function runs in State 1 only
Receives all the udp datagrams received on
multicast ip address to receive add request

Keeps the add requests in a hashtable
the hash function hashes the ip address.

The hashtable is updating constantly with the last
time of packet arrival

This hashtable is used to construct a sorted
list with ip addresses
*/
func (erm *MembershipTreeManager) ProcessInternalEvent(intev InternalEvent) {
	switch {
	case erm.myState.currentState == 1:
		fmt.Println("internal state:", intev)
		ch := make(chan utilities.HeartBeatUpperStack)
		go multicastheartbeatserver.CatchDatagramsAndBounce(ch)
		timeout := time.After(10 * time.Second)
		for {
			select {
			case s := <-ch:
				//fmt.Println("Received:\n", s)
				fmt.Printf("Received upper stack:%v\n", s)
			case <-timeout:
				fmt.Printf("This function runs for 10 seconds only")
				return
			default:
				// Do other activities like sending membership
				// heartbeats to successors in the circle

				time.Sleep(1 * time.Second)
			}
		}
	case erm.myState.currentState == 2:
		r := rand.New(rand.NewSource(99))
		timeout := time.After(15 * time.Second)
		heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages("224.0.0.1", "10001", "10002")

		for {
			hbMessage := utilities.HeartBeat{
				Cluster:   []string{"amit", "kumar", "singh"},
				ReqNumber: r.Int63(),
				ReqCode:   5,
			}

			select {
			case <-timeout:
				fmt.Printf("INFO:Send heartbeats for 15 minutes")
				break
			default:
				time.Sleep(1 * time.Second)
				heartbeatChannelOut <- hbMessage
			}
		}
	}
}

func NewMembershipManager(state State) *MembershipTreeManager {
	erm := new(MembershipTreeManager)
	erm.myState = state
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname")
	}
	if hostname == "leader.assignment2" {
		erm.myState.currentState = 1
	} else {
		erm.myState.currentState = 2
	}
	erm.groupInfo = []string{}
	return erm
}

/*
func (erm *MembershipTreeManager) SendAddRequest() {
	if erm.myState.currentState == 2 {
		leaderIpAddress := erm.myState.leaderIp + ":" + erm.myState.leaderPort
		ServerAddr, err := net.ResolveUDPAddr("udp", leaderIpAddress)
		CheckError(err)

		LocalAddr, err := net.ResolveUDPAddr("udp", "")
		CheckError(err)
		fmt.Println("Localaddress:", LocalAddr.String())
		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		CheckError(err)

		defer Conn.Close()
		i := 0
		for {
			msg := strconv.Itoa(i)
			i++
			msg += "add"
			buf := []byte(msg)
			_, err := Conn.Write(buf)
			if err != nil {
				fmt.Println(msg, err)
			}
			time.Sleep(time.Second * 1)
		}
	}

}
*/
