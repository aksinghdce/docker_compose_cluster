package membership

import "container/list"

/*
1. This structure will just have one state that it has to maintain
consistent. 

2. Heartbeat can be calculated from a consistent state
only.

3. This service will expose a method that the users of this package
can call, this method will exclusively be present in this file for
maintainability.
*/
type Membership struct{
	ring list.List
}

// func (m *Membership) init() {
// 	m.ring = new list.New()
// }

func (m *Membership) Insert(key string, value interface{}) bool {
	return false
}

func (m *Membership) Get(key string) interface{} {
	return false
}