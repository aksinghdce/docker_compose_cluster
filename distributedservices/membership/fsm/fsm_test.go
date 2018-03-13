package fsm

import "testing"

func TestFsm(t *testing.T) {
    fsm1 := Init(1)
    fsm2 := Init(2)
    fsm1.ProcessFsm()
    fsm2.ProcessFsm()
    t.Log("Well done!")
}