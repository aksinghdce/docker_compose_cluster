package membershipmanager

import (
	"app/utilities"
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
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
	CurrentState int8
	/*State 0 data*/
	LeaderIp       string
	LeaderPort     int
	ManagedNodes   []string
	AmITheLeader   bool
	ClusterMap     map[string]utilities.HeartBeat
	RequestContext context.Context
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
	AddNodeToGroup(chan utilities.HeartBeatUpperStack) error
	RemoveNodeFromGroup() (error, string)
}

type MManagerSingleton struct {
	MyState   State
	GroupInfo []string
}

var instance *MManagerSingleton
var once sync.Once

func GetInstance() *MManagerSingleton {
	once.Do(func() {
		instance = &MManagerSingleton{
			MyState: State{
				CurrentState: 0,
				LeaderIp:     "124.0.0.1",
				LeaderPort:   10001,
				ManagedNodes: []string{},
				AmITheLeader: false,
				ClusterMap:   make(map[string]utilities.HeartBeat),
			},
			GroupInfo: []string{},
		}

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("Error getting hostname")
		}
		if hostname == "leader.assignment2" {
			instance.MyState.CurrentState = 1
		} else {
			instance.MyState.CurrentState = 2
		}
	})
	return instance
}

/*Write heartbeat*/
func (erm *MManagerSingleton) aggregator(hb utilities.HeartBeatUpperStack) {
	_, ok := erm.MyState.ClusterMap[hb.Ip]
	if ok {
		// Populate the ClusterMap with value type as another map
		// In the sub-map keep the latest heartbeat from the sender
		// As soon as the map is added to, send the cluster map to
		// a go routine for reacting on it. This is how
		// receiving the heartbeats can be decoupled with reacting
		// on the heartbeat.
		// Make a tree diagram to illustrate a use case before coding
	}
}

func (erm *MManagerSingleton) AddNodeToGroup(hbu utilities.HeartBeatUpperStack) error {
	erm.MyState.ClusterMap[hbu.Ip] = hbu.Hb
	fmt.Printf("State:%v\n", erm.MyState)
	return nil
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
func (erm *MManagerSingleton) ProcessInternalEvent(intev InternalEvent) {
	switch {
	case erm.MyState.CurrentState == 1:
		// For "Add to the group" requests membership service of state 1 listens
		// on 224.0.0.1:10001
		// for regular heartbeats, it must listen on it's unicast ip address
		// A peer in State 2 will send "ADD" request on multicast address because
		// At run time State 2 nodes only know their own ip address. Practically
		// every node is discovering the listener. Once the listener Add's it begins receiving
		// heartbeats from at least one node.
		ch := multicastheartbeatserver.CatchMultiCastDatagramsAndBounce("224.0.0.1", "10001")
		/*Listen to Add request only for 1 second and react to it by sending the received heartbeat
		to collector go routine.
		*/
		timeout := time.After(10 * time.Second)
		for {
			select {
			case s := <-ch:
				/*
					Expect an ADD request. Invoke the aggregator's collector
					routine to updat the internal data structures.
				*/
				fmt.Printf("Received upper stack:%v\n", s)
				erm.AddNodeToGroup(s)
			case <-timeout:
				/*Run State 3 go routines by populating a channel*/
				fmt.Printf("Print the state %v\n", erm.MyState)
			default:
				// Do other activities like sending membership
				// heartbeats to successors in the circle

				time.Sleep(1 * time.Second)
			}
		}
	case erm.MyState.CurrentState == 2:
		/*In this state the machine send ADD request to multicast udp port for 5 seconds
		every 100 milliseconds and then returns.
		So, we need the caller to call this function multiple times
		till the machine goes into State 3
		*/
		r := rand.New(rand.NewSource(99))
		/*SendHeartBeatMessages returns a channel in which you can write your heartbeat messages
		 */
		// heartbeatChannelOut is a channel of utilities.HeartBeat. It returns heartbeats received on
		// Multicast udp port 10001. We are sending on port 10002
		heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages("224.0.0.1", "10001", "10002")

		// ch is a channel of utilities.HeartBeatUpperStack to listen to heartbeats on unicast
		// udp port 10003
		//ch := make(chan utilities.HeartBeatUpperStack)
		// For "Add to the group" requests membership service of state 1 listens
		// on 224.0.0.1:10001
		//go multicastheartbeatserver.CatchUniCastDatagramsAndBounce("10003", ch)
		timeout := time.After(5 * time.Second)
		for {
			/*
				Prepare add requests to be sent to the Introducer
			*/
			hbMessage := utilities.HeartBeat{
				Cluster:   []string{"amit", "kumar", "singh"},
				ReqNumber: r.Int63(),
				ReqCode:   5,
			}

			select {
			//case hbRcv := <-ch:
			//	fmt.Printf("Received Heartbeat:%v\n", hbRcv)
			case <-timeout:
				// Add will be attempted for 5 seconds, every 100 milliseconds
				fmt.Printf("INFO:time to do something else")
			default:
				// Send ADD request every 100 millisecond
				//time.Sleep(100 * time.Millisecond)
				heartbeatChannelOut <- hbMessage
			}
		}
	}
}

/*
func (erm *MManagerSingleton) SendAddRequest() {
	if erm.MyState.currentState == 2 {
		leaderIpAddress := erm.MyState.leaderIp + ":" + erm.MyState.leaderPort
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
