package main

import (
	"testing"
)

func TestLocalGrep(t *testing.T) {
	c := LocalGrep("grep", "tanuki", "machine1.log")

	localGrepResults := <-c
	if len(localGrepResults) > 0 {
		t.Log("LocalGrep returned non null results string")
	}
}
