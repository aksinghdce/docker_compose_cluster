package utilities

import (
	"testing"
)

func TestLocalGrep(t *testing.T) {
	var output string
	output = LocalGrep([]string{"grep", "-c", "8080", "/go/src/app/Dockerfile"})
	if output != "1\n" {
		t.Fatalf(`Sprintf("%%s", empty("7")) = %q want %q`, output, "7")
	}
}
