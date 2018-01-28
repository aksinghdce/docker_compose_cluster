package utilities

import (
	"testing"
)

func TestLocalGrep(t *testing.T) {
	var output string
	output = LocalGrep("grep", "-c", "tanuki", "/go/src/app/machine1.log")
	if output != "7\n" {
		t.Fatalf(`Sprintf("%%s", empty("7")) = %q want %q`, output, "7")
	}
}
