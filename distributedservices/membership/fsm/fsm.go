package fsm

/*This module is responsible for managing heartbeats
*/

/*
This function will send messages to Membership Service
It will ask Membership service for current state so that
it can send that to peers as heartbeat message.

It will send peer's heartbeat messages to Membership service
so that the state can be updated.
*/
func ProcessEvent() {
	//Get channels from Membership service
	//Do plumbing to sendHeartbeat and receiveHeartbeat
	go sendHeartbeat()
	go receiveHeartbeat()
}

func sendHeartbeat() {
	for {

	}
}

func receiveHeartbeat() {
	for {

	}
}