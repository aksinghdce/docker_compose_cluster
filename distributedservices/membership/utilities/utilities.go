package utilities

import (
	"sync/atomic"
)

/*
These utilities are meant for the use by membership service.
All the utilities must be tested for accuracy and speed
*/

/*
An internal counter for requests to the network
To maintain a vector clock in association with hostname

For every event that the node sends or receives as a membership
service the Sequence will be incremented atomicly.

This counter must be managed by the membership service, 
We will follow the philosophy of single point of maintainance.
http://www.ifsq.org/single-point-of-maintenance.html

*/

type count64 int64

type vectorClock struct{
	Sequence int64
	Hostname string
}

func (c *count64) increment() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *count64) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}