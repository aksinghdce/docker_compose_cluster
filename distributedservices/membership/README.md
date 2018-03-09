# A class that has a state and an fsm engine to change the state

# How to write the test case?
In the test case, the membership class will be initialized with a state and the fsm engine will be launched.
In the fsm engine, it will check the current state, the current environment and decide upon the first hard state.

The default state will be 0
Depending upon the environment the state will change to 1 or 2

In state 1, it will expect heartbeat messages with "ADD" request.
The request can be accepted by default. In state 1 the machine will accept
"ADD" requests and respond with "ACK".

In state 2, it will send heartbeats with "ADD" and move to state 3 if the "ACK" is received.

