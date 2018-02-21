package membershipmanager

/*
membershipmanager package manages a statemachine
The statemachine keeps the distributed cluster state
*/

import (
	"app/utilities"
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
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

/*
An internal counter for requests to the network
To maintain a vector clock in association with hostname
*/
type count64 int64

func (c *count64) increment() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *count64) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

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

type InternalEvent struct {
	RequestNumber count64
	Ctx           context.Context
}

type MembershipManager interface {
	/*ProcessInternalEvent manages all the states*/
	ProcessInternalEvent(intevent InternalEvent) (bool, *MManagerSingleton)
	GetGroupInfo() []string // NOT IMPLEMENTED YET
	AddNodeToGroup(chan utilities.HeartBeat) error
	RemoveNodeFromGroup() (error, string) // NOT IMPLEMENTED YET
}

/*
Specification:
The Membership Manager is a Singleton. Helps in unit testing
and initialization from the main function concurrently.
*/
type MManagerSingleton struct {
	MyState   State
	GroupInfo []string /*This is the list every node will use to run the algorithm
	for membership service in FSM State 3*/
	LastHeartbeatReceived utilities.HeartBeat
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
			GroupInfo: nil,
			LastHeartbeatReceived: utilities.HeartBeat{
				Cluster: nil,
				ReqNumber: 0,
				ReqCode: 0,
				FromTo: utilities.MessageAddressVector{
					FromIp: "",
					ToIp: "",
				},
			},
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

func (erm *MManagerSingleton) AddNodeToGroup(intev InternalEvent, hb utilities.HeartBeat) error {
	erm.MyState.ClusterMap[hb.FromTo.FromIp] = hb
	//Create a list of ip addresses added in the ClusterMap
	nodeList := make([]string, 5, 30)
	for ip, _ := range erm.MyState.ClusterMap {
		if ip != "" {
			nodeList = append(nodeList, ip)
		}
	}
	//Assign the latest list of nodes to the GroupInfo
	erm.GroupInfo = nodeList
	heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, hb.FromTo.FromIp, "50009")
	hbMessage := utilities.HeartBeat{
		Cluster:   erm.GroupInfo,
		ReqNumber: intev.RequestNumber.get(),
		ReqCode:   2, //1 is for ADD request, 2 is for Acknowledgement
		FromTo: utilities.MessageAddressVector{
			FromIp: "",
			ToIp: hb.FromTo.FromIp,
		},
	}
	heartbeatChannelOut <- hbMessage
	return nil
}

func (erm *MManagerSingleton) SendState3HeartBeats(intev InternalEvent) {
	for _, ip := range erm.GroupInfo {
		if ip != "" {
			heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, ip, "50012")
			go func() {

				hbMessage := utilities.HeartBeat{
					Cluster:   erm.GroupInfo,
					ReqNumber: intev.RequestNumber.get(),
					ReqCode:   3, //1 is for ADD request
				}
				heartbeatChannelOut <- hbMessage

			}()
		}
	}

}

/*
Specification:

Input: InternalEvent : This carries a context.Context and a sequence number
Output: returns whether the process want to run again to transition state

Processing:
This is the function that runs a finite state machine with 3 states
as described at the beginning of the file.

Receives all the udp datagrams received on
multicast ip address to receive add request

Keeps the add requests in a hashtable
the hash function hashes the ip address.

The hashtable is updating constantly with the last
time of packet arrival

This hashtable is used to construct a sorted
list with ip addresses
*/
func (erm *MManagerSingleton) ProcessInternalEvent(intev InternalEvent) bool {
	switch {
	case erm.MyState.CurrentState == 1:
		// For "Add to the group" requests membership service of state 1 listens
		// on 224.0.0.1:10001
		// for regular heartbeats, it must listen on it's unicast ip address
		// A peer in State 2 will send "ADD" request on multicast address because
		// At run time State 2 nodes only know their own ip address. Practically
		// every node is discovering the listener. Once the listener Add's it begins receiving
		// heartbeats from at least one node.
		ch := multicastheartbeatserver.CatchMultiCastDatagramsAndBounce(intev.Ctx, "224.0.0.1", "10001")
		/*Listen to Add request only for 1 second and react to it by sending the received heartbeat
		to collector go routine.
		*/

		for {
			timeout := time.After(2 * time.Second)
			select {
			case s := <-ch:
				/*
					Expect an ADD request. Invoke the aggregator's collector
					routine to updat the internal data structures.
				*/
				erm.AddNodeToGroup(intev, s)
			case <-timeout:
				/*Run State 3 go routines by populating a channel*/
				erm.SendState3HeartBeats(intev)

			default:
				// Do other activities like sending membership
				// heartbeats to successors in the circle
				continue
			}
		}
	case erm.MyState.CurrentState == 2:
		/*In this state the machine send ADD request to multicast udp port for 5 seconds
		every 100 milliseconds and then returns.
		So, we need the caller to call this function multiple times
		till the machine goes into State 3
		*/
		/*SendHeartBeatMessages returns a channel in which you can write your heartbeat messages
		 */
		// heartbeatChannelOut is a channel of utilities.HeartBeat. It returns heartbeats received on
		// Multicast udp port 10001. We are sending on port 10002
		heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, "224.0.0.1", "10001")

		// ch is a channel of utilities.HeartBeatUpperStack to listen to heartbeats on unicast
		// udp port 10002
		heartbeatChannelIn := multicastheartbeatserver.CatchUniCastDatagramsAndBounce(intev.Ctx, "50009")

		for {
			timeout := time.After(5 * time.Second)
			/*
				Prepare add requests to be sent to the Introducer
			*/
			hbMessage := utilities.HeartBeat{
				Cluster:   []string{},
				ReqNumber: intev.RequestNumber.get(),
				ReqCode:   1, //1 is for ADD request
				FromTo: utilities.MessageAddressVector{
					FromIp: "",
					ToIp: "224.0.0.1",
				},
			}

			intev.RequestNumber.increment()

			select {
			case hbRcv := <-heartbeatChannelIn:
				/*If the heartbeat message contains ReqCode 2, then we must stop sending
				ADD requests and instead send Keep requests to our successor in the GroupInfo
				*/
				if hbRcv.ReqCode == 2 {
					utilities.Log(intev.Ctx, "STATE Transition 2->3\n")
					erm.MyState.CurrentState = 3
					erm.GroupInfo = hbRcv.Cluster
					erm.LastHeartbeatReceived = hbRcv
					// Ask the caller to rerun this function: To change state to 3
					return true
				}
			case <-timeout:
				// Add will be attempted for 5 seconds, every 100 milliseconds
				fmt.Printf("INFO:time to do something else\n")
			default:
				// Send ADD request every 100 millisecond
				//time.Sleep(100 * time.Millisecond)
				heartbeatChannelOut <- hbMessage
			}
		}
	case erm.MyState.CurrentState == 3:
		fmt.Print("Running in state 3 now\n")
		heartbeatChannelIn3 := multicastheartbeatserver.CatchUniCastDatagramsAndBounce(intev.Ctx, "50012")

		for {
			timeout := time.After(5 * time.Second)
			select {
			case hbst3 := <-heartbeatChannelIn3:
				if hbst3.ReqCode == 3 {
					erm.GroupInfo = hbst3.Cluster
				}
			case <-timeout:
				fmt.Printf("Cluster Info:%v\n", erm)
				time.Sleep(1 * time.Millisecond)
			}

		}
	}
	return false
}
