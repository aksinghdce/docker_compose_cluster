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
	"sort"
	"strings"

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

/*Constants
*/

const (
	SEND_HEARTBEAT_EVERY = 10 * time.Millisecond
	DELETE_OLDER_THAN = 1 * time.Second
)

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
	MyIp       string
	SendPort     int
	ListenPort int
	ManagedNodes   []string
	ClusterMap     map[string]time.Time
	RequestContext context.Context
}

type InternalEvent struct {
	RequestNumber count64
	Ctx           context.Context
}

/*
Specification:
The Membership Manager is a Singleton. Helps in unit testing
and initialization from the main function concurrently.
*/
type MManagerSingleton struct {
	LeaderMultiCastIp string
	LeaderUniCastIp string
	MyState   State
	GroupInfo map[string]bool /*This is the list every node will use to run the algorithm
	for membership service in FSM State 3*/
	LastHeartbeatReceived utilities.HeartBeat
}

var instance *MManagerSingleton
var once sync.Once

/*
FSM State 0 to FSM State 1 and FSM State 2 happens early
based on hostname
*/
func GetInstance() *MManagerSingleton {
	once.Do(func() {
		instance = &MManagerSingleton{
			LeaderMultiCastIp: "224.0.0.1",
			LeaderUniCastIp: "",
			MyState: State{
				CurrentState: 0,
				ManagedNodes: []string{},
				ClusterMap:   make(map[string]time.Time),
			},
			GroupInfo: make(map[string]bool),
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

func (erm *MManagerSingleton) DeleteOlderHeartbeats(dur time.Duration) {
	for ip, timestamp := range erm.MyState.ClusterMap {
		if timestamp.Sub(time.Now()) >= dur {
			delete(erm.MyState.ClusterMap, ip)
		}
	}
}

func (erm *MManagerSingleton) AddNodeToGroup(intev InternalEvent, senderIp string) {
	erm.MyState.ClusterMap[senderIp] = time.Now()
	erm.CheckStaleness()
}

func (erm *MManagerSingleton) CheckStaleness() {
	for ip, _ := range erm.MyState.ClusterMap {
		if len(ip) != 0 {
			stale := time.Now().Sub(erm.MyState.ClusterMap[ip]) > (100 * time.Millisecond)
			erm.GroupInfo[ip] = stale
		}
	}
}

func (erm *MManagerSingleton) SendAckToAddRequester(intev InternalEvent, ip, port string) {
	heartbeatChannelOut := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, ip, port)
	hbMessage := utilities.HeartBeat{
		Cluster:   erm.GroupInfo,
		ReqNumber: intev.RequestNumber.get(),
		ReqCode:   2, //1 is for ADD request, 2 is for Acknowledgement
		FromTo: utilities.MessageAddressVector{
			FromIp: "",
			ToIp: ip,
		},
	}
	heartbeatChannelOut <- hbMessage
}

func (erm *MManagerSingleton) SortCurrentGroupInfo() []string{
	removedEmpty := make([]string, 0)
	for ip, stale := range erm.GroupInfo {
		if len(ip) != 0 && !stale{
			removedEmpty = append(removedEmpty, ip)
		}
	}
	sort.Strings(removedEmpty)
	return removedEmpty
}

func (erm *MManagerSingleton) WhomToSendHb() (string, error) {
	sortedIps := erm.SortCurrentGroupInfo()
	myIndex := -1
	for i, ip := range sortedIps {
		if erm.MyState.MyIp == ip {
			myIndex = i
		}
	}
	if myIndex == -1 {
		err := fmt.Errorf("your own ip:%v is not in your groupinfo\n", erm.MyState.MyIp)
		return "", err
	}
	if myIndex < (len(sortedIps) - 1) {
		return sortedIps[myIndex + 1], nil
	} else {
		return sortedIps[0], nil
	}
}


func (erm *MManagerSingleton) ConsolidateInfo(remoteMap map[string]bool) {
	
	if len(remoteMap) > len(erm.GroupInfo) {
		for key, state := range remoteMap {
			erm.GroupInfo[key] = state
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
		// For all requests membership service of state 1 listens
		// on 224.0.0.1:10001
		// for regular heartbeats, it must listen on it's unicast ip address
		// A peer in State 2 will send "ADD" request on multicast address because
		// At run time State 2 nodes only know their own ip address. Practically
		// every node is discovering the listener. Once the listener Add's it begins receiving
		// heartbeats from at least one node.
		ch := multicastheartbeatserver.CatchMultiCastDatagramsAndBounce(intev.Ctx, "224.0.0.1", "10001")
		sortedIps := erm.SortCurrentGroupInfo()
		channelArr := make([]chan utilities.HeartBeat, 0)
		for _, ip := range sortedIps {
			channelArr = append(channelArr, multicastheartbeater.SendHeartBeatMessages(intev.Ctx, ip, "50012"))
		}

		for {
			timeout := time.After(SEND_HEARTBEAT_EVERY)

			select {
			case s := <-ch:
				/*
					Expect an ADD request. Invoke the aggregator's collector
					routine to updat the internal data structures.
				*/
				//Delete heartbeats older than 20 milliseconds
				erm.DeleteOlderHeartbeats(DELETE_OLDER_THAN)
				erm.AddNodeToGroup(intev, s.FromTo.FromIp)
				if s.ReqCode == 1 {
					erm.SendAckToAddRequester(intev, s.FromTo.FromIp, "50009")
				}
			case <-timeout:
				for i, chpeer := range channelArr {
					chpeer <- utilities.HeartBeat{
						Cluster:   erm.GroupInfo,
						ReqNumber: intev.RequestNumber.get(),
						ReqCode:   3, //1 is for ADD request
						FromTo: utilities.MessageAddressVector{
							FromIp: erm.MyState.MyIp,
							ToIp: sortedIps[i],
						},
					}
				}
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

		heartbeatChannelIn := multicastheartbeatserver.CatchUniCastDatagramsAndBounce(intev.Ctx, "50009")

		for {
			/*
				Prepare add requests to be sent to the Introducer
			*/
			hbMessage := utilities.HeartBeat{
				Cluster:   nil,
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
				ip_port := strings.Split(hbRcv.FromTo.ToIp, ":")
				erm.MyState.MyIp = ip_port[0]

				// Set leader's unicast ip
				ip_port_leader := strings.Split(hbRcv.FromTo.FromIp, ":")
				erm.LeaderUniCastIp = ip_port_leader[0]
				fmt.Printf("Updated myIp to:%v and Leader unicast to: %v\n", erm.MyState.MyIp, erm.LeaderUniCastIp)
				if hbRcv.ReqCode == 2 {
					utilities.Log(intev.Ctx, "STATE Transition 2->3\n")
					erm.MyState.CurrentState = 3
					erm.GroupInfo = hbRcv.Cluster
					erm.LastHeartbeatReceived = hbRcv
					// Ask the caller to rerun this function: To change state to 3
					return true
				}
			default:
				heartbeatChannelOut <- hbMessage
			}
		}
	case erm.MyState.CurrentState == 3:
		heartbeatChannelToListener := multicastheartbeater.SendHeartBeatMessages(intev.Ctx, "224.0.0.1", "10001")
		heartbeatChannelIn := multicastheartbeatserver.CatchUniCastDatagramsAndBounce(intev.Ctx, "50012")
		sendTo, err := erm.WhomToSendHb()
		if err != nil {
			utilities.Log(intev.Ctx, err.Error())
		}
		var heartbeatChannelOut chan utilities.HeartBeat
			
		if len(sendTo) > 0 {
			heartbeatChannelOut = multicastheartbeater.SendHeartBeatMessages(intev.Ctx, sendTo, "50012")
		}else {
			fmt.Print("blank sendTo")
		}

		for {
			//We will delete heartbeats older than 100 milliseconds
			timeout := time.After(SEND_HEARTBEAT_EVERY)
			select {
			case hbst := <-heartbeatChannelIn:
				//fmt.Printf("I:%s have:%v and received:%v\n", erm.MyState.MyIp, erm.GroupInfo, hbst.Cluster)
				if hbst.ReqCode == 3 {
					//erm.ConsolidateInfo(hbst.Cluster)
					erm.GroupInfo = hbst.Cluster
				}
				ip_port := strings.Split(hbst.FromTo.ToIp, ":")
				erm.MyState.MyIp = ip_port[0]
				erm.DeleteOlderHeartbeats(DELETE_OLDER_THAN)
				erm.AddNodeToGroup(intev, hbst.FromTo.FromIp)
			case <-timeout:
				if len(sendTo) > 0 {
					heartbeatChannelOut <- utilities.HeartBeat{
						Cluster:   erm.GroupInfo,
						ReqNumber: intev.RequestNumber.get(),
						ReqCode:   3, //1 is for ADD request
						FromTo: utilities.MessageAddressVector{
							FromIp: erm.MyState.MyIp,
							ToIp: sendTo,
						},
					}
				}

				heartbeatChannelToListener <- utilities.HeartBeat{
					Cluster:   erm.GroupInfo,
					ReqNumber: intev.RequestNumber.get(),
					ReqCode:   3, //1 is for ADD request
					FromTo: utilities.MessageAddressVector{
						FromIp: erm.MyState.MyIp,
						ToIp: erm.LeaderUniCastIp,
					},
				}
			}

		}
	}
	return false
}
