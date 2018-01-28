package utilities

import (
	"testing"
)

func TestLocalGrep(t *testing.T) {
	var output string
	output = LocalGrep("grep", "-c", "tanuki", "machine1.log")
	//output = "7"
	if output != "7" {
		t.Errorf(`Sprintf("%%s", empty("7")) = %q want %q`, output, "7")
	}
}
