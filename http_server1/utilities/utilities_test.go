package utilities

import (
	"testing"
)

func TestLocalGrep(t *testing.T) {
	tt := []struct{
		commandstringsslice []string
		output string
	}{
		{[]string{"grep", "-c", "8080", "/go/src/app/Dockerfile"}, "1"},
	}

	for _, tc := range tt {
		var output string
		output = LocalGrep(tc.commandstringsslice)
		outputexpected := tc.output
		outputexpected += "\n"
		if output != outputexpected {
			t.Fatalf(`Sprintf("%%s", empty("7")) = %q want %q`, output, "7")
		}
	}
	
}
