------------------------- MODULE membership_service -------------------------
(*Specification for maintaining a consistent view of membership list*)
EXTENDS Integers
VARIABLE node, chan
CONSTANT IP, Data, ADD_Request, ACK_Response

TypeInvariant == /\ node \in [type:{1, 2}, rdy:{0, 1}, ack:{0, 1}, ip:IP] \*node has IP, it is either a Leader or a Peer
                 /\ IPs = IP
                 /\ chan \in [type:{1, 2}, rdy:{0, 1}, ack:{0, 1}, val:{<<>>, <<IPs, IPd>>, <<IP, MembershipList>>}] (*
                 chan is either of sending type or receiving type
                 chan carries either empty value (when it's not in use), [IPs and IPs] 
                 (when sending ADD request, source IP: IPs
                 , destination IP: IPd), or [IPs, IPd and Membership
                 list] (when sending regular heartbeats)*)

LeaderInit == \/ /\ node.type = 1
                  /\ node.rdy = node.ack \* Leader begins in listening mode: Listening for ADD req
                  /\ node.membership = <<>> \* Leader begins with empty membership list

PeerInit == \/ /\ node.type = 2
               /\ node.rdy = node.ack \* Peer begins in listening mode: Listening for ACK res
               /\ node.membership = <<>> \* Peer begins with empty membership list

PeerSendAddReq(d) == /\ PeerInit
                     /\ d = ADD_Request
                     /\ chan' = [chan EXCEPT !.val = Append(@, IP), !.rdy = 1 - @]             
LeaderReceiveAddEvent == \/ /\ LeaderInit
                            /\ node' = [node EXCEPT !.ack = 1 - @]
LeaderSendAckEvent(d) == /\ LeaderReceiveAddEvent
                         /\ \* send ack to the machine that sent ADD request
PeerReceiveAckEvent == /\ PeerSendAddReq
                       /\ \* peer do something after receiving ACK response
PeerSendRegularHb == /\ PeerReceiveAckEvent
                     /\ \*what to send in the regular Hb?
PeerReceiveRegularHb == /\ PeerReceiveAckEvent
                     /\ \*what to do with received Hb?
LeaderNext == /\ LeaderReceiveAddEvent
              /\ node' = [node EXCEPT !.membership = Append(@, chan.val.IPs)]
LeaderSendRegularHb == /\ (\exists d \in Data : LeaderSendAckEvent(d))
                       /\ \* What to send in regular HB as a Leader
LeaderReceiveRegularHb == /\ (\exists d \in Data : LeaderSendAckEvent(d))
                       /\ \* What to do with received regular HB as a Leader
(*How to write a temporal logic for membership service such that

1. There are two types of machines in the system - Leader and Peer
2. Leader begins with a blank membership list
3. Peer knows about the Leader's Ip address
4. Peer begins by sending ADD request to Leader, ADD request carries requester's IP address.
5. Leader listens for ADD request
6. Leader receives ADD request from Peer i
7. Leader updates its membership list with thie Peer i's IP Address and timestamp
8. Leader sends ACK with current membership list of Ip Addresses and timestamps
9. Leader listens for heartbeat messages from all Peers
10. All Peers send heartbeat messages to Leader
11. All Peers send heartbeat messages to two Ip addresses higher than itself, mod calculations is considered
12. All Peers receive heartbeat messages from two peers less than itself, mod calculations is considered*)

=============================================================================
\* Modification History
\* Last modified Wed Apr 11 13:33:58 PDT 2018 by aksin
\* Created Tue Mar 27 15:24:27 PDT 2018 by aksin
